import { useCallback, useEffect, useRef, useState } from "react";
import { useShow } from "@refinedev/core";
import { Show } from "@refinedev/antd";
import {
  Alert,
  App,
  Button,
  Card,
  Descriptions,
  Empty,
  List,
  Popconfirm,
  Result,
  Space,
  Switch,
  Table,
  Tabs,
  Tag,
  Typography,
} from "antd";
import axios from "axios";
import ReactMarkdown from "react-markdown";
import { Link, useParams } from "react-router";

import { env } from "../../config/env";
import { httpClient } from "../../providers";
import type { ManagedTalkSession } from "./edit-form";
import {
  type ManagedOpinion,
  type RegenerateType,
  type TalkSessionAnalysis,
  groupSummaries,
  opinionStatusTags,
  regenerateBody,
  reportVisibilityBody,
  safeExternalURL,
  shouldReloadAfterRegenerate,
  voteSummary,
} from "./detail-view";

const formatDate = (iso?: string) => (iso ? new Date(iso).toLocaleString("ja-JP") : "—");
const yesNo = (value: boolean, yes: string, no: string) =>
  value ? <Tag color="green">{yes}</Tag> : <Tag>{no}</Tag>;
const formHeaders = { "Content-Type": "application/x-www-form-urlencoded" };

export function TalkSessionShow() {
  const { id = "" } = useParams();
  const { query, result: session } = useShow<ManagedTalkSession>({ resource: "talksessions", id });
  const notFound = axios.isAxiosError(query.error) && query.error.response?.status === 404;

  return (
    <Show
      title={session?.theme ?? "セッション詳細"}
      isLoading={query.isLoading}
      canDelete={false}
      headerButtons={() => (
        <Link to={`/talksessions/edit/${encodeURIComponent(id)}`}>
          <Button type="primary">編集する</Button>
        </Link>
      )}
    >
      {query.isError ? (
        notFound ? (
          <Result
            status="404"
            title="セッションが見つかりません"
            extra={<Link to="/talksessions">セッション一覧へ戻る</Link>}
          />
        ) : (
          <Alert
            type="error"
            showIcon
            message="セッションを取得できませんでした"
            action={<Button onClick={() => query.refetch()}>再試行</Button>}
          />
        )
      ) : (
        session && (
          <Tabs
            items={[
              { key: "overview", label: "概要", children: <Overview session={session} /> },
              // セッションが変わったらタブを作り直し、ページ番号などを引き継がない
              { key: "opinions", label: "意見", children: <OpinionsTab key={id} talkSessionID={id} /> },
              {
                key: "analysis",
                label: "AI分析",
                children: (
                  <AnalysisTab key={id} session={session} onSessionChange={() => query.refetch()} />
                ),
              },
            ]}
          />
        )
      )}
    </Show>
  );
}

function Overview({ session }: { session: ManagedTalkSession }) {
  return (
    <Descriptions bordered column={1} size="small">
      <Descriptions.Item label="セッションID">
        <Typography.Text copyable>{session.talkSessionID}</Typography.Text>
      </Descriptions.Item>
      <Descriptions.Item label="テーマ">{session.theme}</Descriptions.Item>
      <Descriptions.Item label="オーナー">
        {session.owner.displayName || "未設定"}
        {session.owner.userID && (
          <>
            {"（"}
            <Link to={`/users/show/${encodeURIComponent(session.owner.userID)}`}>ユーザー詳細</Link>
            {"）"}
          </>
        )}
      </Descriptions.Item>
      <Descriptions.Item label="募集の終了日時">{formatDate(session.scheduledEndTime)}</Descriptions.Item>
      <Descriptions.Item label="場所">
        {[session.prefecture, session.city].filter(Boolean).join(" ") || "未設定"}
      </Descriptions.Item>
      <Descriptions.Item label="参加者数">{session.participantCount}人</Descriptions.Item>
      <Descriptions.Item label="ホームの一覧">{yesNo(!session.hideTop, "表示する", "表示しない")}</Descriptions.Item>
      <Descriptions.Item label="レポート">{yesNo(!session.hideReport, "公開する", "公開しない")}</Descriptions.Item>
      <Descriptions.Item label="意見の分析">{yesNo(!session.disableAnalysis, "行う", "行わない")}</Descriptions.Item>
      <Descriptions.Item label="参加制限">
        {session.restrictions.length ? (
          <Space wrap>
            {session.restrictions.map((key) => (
              <Tag key={key}>{key}</Tag>
            ))}
          </Space>
        ) : (
          "なし"
        )}
      </Descriptions.Item>
      <Descriptions.Item label="作成日時">{formatDate(session.createdAt)}</Descriptions.Item>
      <Descriptions.Item label="更新日時">{formatDate(session.updatedAt)}</Descriptions.Item>
    </Descriptions>
  );
}

const OPINION_PAGE_SIZE = 20;

function OpinionsTab({ talkSessionID }: { talkSessionID: string }) {
  const [page, setPage] = useState(1);
  const [data, setData] = useState<{ totalCount: number; opinions: ManagedOpinion[] }>({
    totalCount: 0,
    opinions: [],
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(false);
  const [retry, setRetry] = useState(0);

  useEffect(() => {
    let active = true;
    setLoading(true);
    setError(false);
    httpClient
      .get<{ totalCount: number; opinions: ManagedOpinion[] }>(
        `/v1/manage/talksessions/${encodeURIComponent(talkSessionID)}/opinions`,
        { params: { limit: OPINION_PAGE_SIZE, offset: (page - 1) * OPINION_PAGE_SIZE } },
      )
      .then((res) => active && setData({ totalCount: res.data.totalCount, opinions: res.data.opinions ?? [] }))
      .catch(() => {
        if (!active) return;
        // ページ番号だけ進んで前のページの意見が残らないよう、表を空にする
        setData((prev) => ({ totalCount: prev.totalCount, opinions: [] }));
        setError(true);
      })
      .finally(() => active && setLoading(false));
    return () => {
      active = false;
    };
  }, [talkSessionID, page, retry]);

  return (
    <Space direction="vertical" style={{ width: "100%" }}>
      {error && (
        <Alert
          type="error"
          showIcon
          message="意見一覧を取得できませんでした"
          action={<Button onClick={() => setRetry((v) => v + 1)}>再試行</Button>}
        />
      )}
      <Typography.Text type="secondary">
        通報の対応で削除された意見や返信も含めて、新しい順に表示します。投票は意見ごとの集計だけで、誰がどう投票したかは表示しません。
      </Typography.Text>
      <Table
        rowKey="opinionID"
        loading={loading}
        dataSource={data.opinions}
        pagination={{
          current: page,
          pageSize: OPINION_PAGE_SIZE,
          total: data.totalCount,
          onChange: setPage,
          showSizeChanger: false,
          showTotal: (total) => `${total}件`,
        }}
        columns={[
          {
            title: "投稿日時",
            dataIndex: "createdAt",
            width: 170,
            render: (v: string) => formatDate(v),
          },
          {
            title: "投稿者",
            key: "author",
            width: 160,
            render: (_: unknown, row: ManagedOpinion) =>
              row.author.userID ? (
                <Link to={`/users/show/${encodeURIComponent(row.author.userID)}`}>
                  {row.author.displayName || row.author.displayID || "未設定"}
                </Link>
              ) : (
                "—"
              ),
          },
          {
            title: "内容",
            key: "content",
            render: (_: unknown, row: ManagedOpinion) => (
              <Space direction="vertical" size={4}>
                {opinionStatusTags(row).length > 0 && (
                  <Space wrap size={4}>
                    {opinionStatusTags(row).map((tag) => (
                      <Tag key={tag.label} color={tag.color}>
                        {tag.label}
                      </Tag>
                    ))}
                  </Space>
                )}
                <Typography.Text delete={row.deleted} style={{ whiteSpace: "pre-wrap" }}>
                  {row.content}
                </Typography.Text>
                {safeExternalURL(row.pictureURL) && (
                  <a href={safeExternalURL(row.pictureURL) ?? undefined} target="_blank" rel="noopener noreferrer">
                    画像を開く
                  </a>
                )}
                {safeExternalURL(row.referenceURL) && (
                  <a href={safeExternalURL(row.referenceURL) ?? undefined} target="_blank" rel="noopener noreferrer">
                    参考URL
                  </a>
                )}
              </Space>
            ),
          },
          {
            title: "投票（賛成 / 反対 / 保留）",
            key: "votes",
            width: 200,
            render: (_: unknown, row: ManagedOpinion) => {
              const summary = voteSummary(row);
              return (
                <Space direction="vertical" size={0}>
                  <span>
                    {row.agreeCount} / {row.disagreeCount} / {row.passCount}
                  </span>
                  <Typography.Text type="secondary">
                    計{summary.total}票
                    {summary.agreeRate !== null && `・賛成率 ${summary.agreeRate}%`}
                  </Typography.Text>
                </Space>
              );
            },
          },
        ]}
      />
    </Space>
  );
}

function AnalysisTab({
  session,
  onSessionChange,
}: {
  session: ManagedTalkSession;
  onSessionChange: () => void;
}) {
  const { message } = App.useApp();
  const id = session.talkSessionID;
  const [report, setReport] = useState<{ loading: boolean; text?: string; missing?: boolean; error?: boolean }>({
    loading: true,
  });
  const [analysis, setAnalysis] = useState<{ loading: boolean; data?: TalkSessionAnalysis; error?: boolean }>({
    loading: true,
  });
  const [toggling, setToggling] = useState(false);
  const [regenerating, setRegenerating] = useState<RegenerateType | null>(null);
  // 再生成や再試行で読み込み直したとき、先に始めた遅い応答で新しい表示を上書きしないよう、最後の読み込みだけを反映する
  const latestLoad = useRef(0);

  const load = useCallback(() => {
    const loadID = ++latestLoad.current;
    const isLatest = () => loadID === latestLoad.current;
    setReport({ loading: true });
    setAnalysis({ loading: true });
    httpClient
      .get<{ report?: string }>(`/v1/manage/talksessions/${encodeURIComponent(id)}/analysis/report`)
      .then((res) => isLatest() && setReport({ loading: false, text: res.data.report }))
      .catch(
        (error) =>
          isLatest() &&
          setReport(
            axios.isAxiosError(error) && error.response?.status === 404
              ? { loading: false, missing: true }
              : { loading: false, error: true },
          ),
      );
    // 公開APIなので、401 で管理者のリフレッシュを試みる管理画面用のクライアントは使わない
    axios
      .get<TalkSessionAnalysis>(`${env.apiUrl}/talksessions/${encodeURIComponent(id)}/analysis`)
      .then((res) => isLatest() && setAnalysis({ loading: false, data: res.data }))
      .catch(() => isLatest() && setAnalysis({ loading: false, error: true }));
  }, [id]);

  useEffect(() => {
    load();
    return () => {
      // タブを離れたあとに届いた応答は捨てる
      latestLoad.current++;
    };
  }, [load]);

  const toggleReport = (visible: boolean) => {
    setToggling(true);
    httpClient
      .post(`/v1/manage/talksessions/${encodeURIComponent(id)}/analysis/report`, reportVisibilityBody(!visible), {
        headers: formHeaders,
      })
      .then(() => {
        message.success(visible ? "レポートを公開しました" : "レポートを非公開にしました");
        onSessionChange();
      })
      .catch(() => message.error("レポートの公開設定を変更できませんでした"))
      .finally(() => setToggling(false));
  };

  const regenerate = (type: RegenerateType) => {
    setRegenerating(type);
    httpClient
      .post(`/v1/manage/talksessions/${encodeURIComponent(id)}/analysis/regenerate`, regenerateBody(type), {
        headers: formHeaders,
      })
      .then(() => {
        if (shouldReloadAfterRegenerate(type)) {
          message.success("レポートを再生成しました");
          load();
        } else {
          message.success("再生成を始めました。反映まで時間がかかるので、しばらくしてから開き直してください");
        }
      })
      .catch(() => message.error("再生成を始められませんでした"))
      .finally(() => setRegenerating(null));
  };

  const regenerateButton = (type: RegenerateType, label: string) => (
    <Popconfirm
      title={`${label}を再生成しますか？`}
      description="AIの利用料金がかかります。今の内容は上書きされます。"
      okText="再生成する"
      cancelText="キャンセル"
      onConfirm={() => regenerate(type)}
    >
      <Button loading={regenerating === type} disabled={regenerating !== null && regenerating !== type}>
        {label}を再生成
      </Button>
    </Popconfirm>
  );

  return (
    <Space direction="vertical" size="large" style={{ width: "100%" }}>
      <Space wrap>
        <span>レポートを公開する</span>
        <Switch
          checked={!session.hideReport}
          loading={toggling}
          checkedChildren="公開"
          unCheckedChildren="非公開"
          onChange={toggleReport}
        />
        {regenerateButton("report", "レポート")}
        {regenerateButton("group", "グループ分け")}
        {regenerateButton("image", "画像")}
      </Space>

      <Card title="レポート" size="small" loading={report.loading}>
        {report.error ? (
          <Alert type="error" showIcon message="レポートを取得できませんでした" action={<Button onClick={load}>再試行</Button>} />
        ) : report.missing || !report.text ? (
          <Empty description="まだレポートがありません" />
        ) : (
          // web のレポート表示と同じく react-markdown で整形する（HTML は描画しない）
          <Typography>
            <ReactMarkdown
              components={{
                a: ({ node: _node, ...props }) => <a {...props} target="_blank" rel="noopener noreferrer" />,
              }}
            >
              {report.text}
            </ReactMarkdown>
          </Typography>
        )}
      </Card>

      <Card title="グループ分け" size="small" loading={analysis.loading}>
        {analysis.error ? (
          <Alert type="error" showIcon message="分析結果を取得できませんでした" action={<Button onClick={load}>再試行</Button>} />
        ) : !analysis.data || groupSummaries(analysis.data).length === 0 ? (
          <Empty description="まだグループ分けの結果がありません" />
        ) : (
          <Space direction="vertical" style={{ width: "100%" }}>
            {groupSummaries(analysis.data).map((group) => (
              <Card
                key={group.groupID}
                type="inner"
                size="small"
                title={`${group.groupName}（${group.memberCount}人）`}
              >
                {group.opinions.length === 0 ? (
                  <Typography.Text type="secondary">代表的な意見はありません</Typography.Text>
                ) : (
                  <List
                    size="small"
                    dataSource={group.opinions}
                    renderItem={(item) => (
                      <List.Item
                        extra={
                          <Typography.Text type="secondary">
                            賛成 {item.agreeCount} / 反対 {item.disagreeCount} / 保留 {item.passCount}
                          </Typography.Text>
                        }
                      >
                        <List.Item.Meta
                          title={<span style={{ whiteSpace: "pre-wrap" }}>{item.opinion.content}</span>}
                          description={item.user.displayName || item.user.displayID}
                        />
                      </List.Item>
                    )}
                  />
                )}
              </Card>
            ))}
          </Space>
        )}
      </Card>
    </Space>
  );
}
