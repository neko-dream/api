import { useEffect, useState } from "react";
import { Alert, Button, Table, Tag, Typography } from "antd";
import { Link } from "react-router";

import { httpClient } from "../../providers";

interface TalkSessionStat {
  talkSessionID: string;
  theme: string;
  owner: { displayName?: string };
  scheduledEndTime: string;
  opinionCount: number;
  voteCount: number;
  hideTop: boolean;
}

const PAGE_SIZE = 20;

export const TalkSessionList = () => {
  const [data, setData] = useState<TalkSessionStat[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(1);
  const [error, setError] = useState(false);
  const [retry, setRetry] = useState(0);

  useEffect(() => {
    let active = true;
    setLoading(true);
    setError(false);
    const offset = (page - 1) * PAGE_SIZE;
    httpClient
      .get<{ talkSessionStats: TalkSessionStat[]; totalCount: number }>(
        `/v1/manage/talksessions/list?limit=${PAGE_SIZE}&offset=${offset}`,
      )
      .then((res) => {
        if (!active) return;
        setData(res.data.talkSessionStats ?? []);
        setTotal(res.data.totalCount ?? 0);
      })
      .catch(() => {
        if (active) {
          setData([]);
          setTotal(0);
          setError(true);
        }
      })
      .finally(() => {
        if (active) setLoading(false);
      });
    return () => {
      active = false;
    };
  }, [page, retry]);

  return (
    <div style={{ padding: 24 }}>
      <Typography.Title level={4}>トークセッション</Typography.Title>
      {error && (
        <Alert
          type="error"
          showIcon
          message="トークセッション一覧を取得できませんでした"
          action={<Button onClick={() => setRetry((value) => value + 1)}>再試行</Button>}
        />
      )}
      <Table
        rowKey="talkSessionID"
        loading={loading}
        dataSource={data}
        pagination={{
          current: page,
          pageSize: PAGE_SIZE,
          total,
          onChange: setPage,
          showSizeChanger: false,
        }}
        columns={[
          {
            title: "テーマ",
            dataIndex: "theme",
            render: (theme: string, row: TalkSessionStat) => (
              <Link to={`/talksessions/show/${encodeURIComponent(row.talkSessionID)}`}>{theme}</Link>
            ),
          },
          { title: "オーナー", dataIndex: ["owner", "displayName"] },
          { title: "意見数", dataIndex: "opinionCount", width: 90 },
          { title: "投票数", dataIndex: "voteCount", width: 90 },
          {
            title: "終了予定",
            dataIndex: "scheduledEndTime",
            render: (v: string) => (v ? new Date(v).toLocaleString("ja-JP") : ""),
          },
          {
            title: "ホーム一覧",
            dataIndex: "hideTop",
            width: 110,
            render: (hideTop: boolean) =>
              hideTop ? <Tag>表示しない</Tag> : <Tag color="green">表示する</Tag>,
          },
          {
            title: "",
            key: "actions",
            width: 80,
            render: (_: unknown, row: TalkSessionStat) => (
              <Link to={`/talksessions/edit/${encodeURIComponent(row.talkSessionID)}`}>編集</Link>
            ),
          },
        ]}
      />
    </div>
  );
};
