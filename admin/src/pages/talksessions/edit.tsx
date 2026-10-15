import { useEffect, useState } from "react";
import { useCustomMutation, useInvalidate, useShow } from "@refinedev/core";
import { Edit } from "@refinedev/antd";
import {
  Alert,
  App,
  Button,
  Checkbox,
  Col,
  Descriptions,
  Form,
  Input,
  Result,
  Row,
  Space,
  Switch,
  Typography,
} from "antd";
import axios from "axios";
import { Link, useNavigate, useParams } from "react-router";

import { env } from "../../config/env";
import { httpClient } from "../../providers";
import { buildPreviewDocument } from "./description-preview";
import {
  type ManagedTalkSession,
  type RestrictionOption,
  type TalkSessionEditValues,
  needsRestrictionWarning,
  toEditRequest,
  toEditValues,
  validateEditValues,
} from "./edit-form";

// 表示設定は「隠すかどうか」で持っているので、画面では肯定形（表示する・行う）に反転して見せる
const invert = {
  getValueProps: (value: boolean) => ({ checked: !value }),
  normalize: (checked: boolean) => !checked,
  valuePropName: "checked",
};

export function TalkSessionEdit() {
  const { id = "" } = useParams();
  const navigate = useNavigate();
  const { modal } = App.useApp();
  const [form] = Form.useForm<TalkSessionEditValues>();
  const description = Form.useWatch("description", form) ?? "";
  const { query, result: session } = useShow<ManagedTalkSession>({
    resource: "talksessions",
    id,
  });
  const { mutate, mutation } = useCustomMutation();
  const invalidate = useInvalidate();
  const [options, setOptions] = useState<RestrictionOption[] | null>(null);
  const [optionsError, setOptionsError] = useState(false);

  useEffect(() => {
    let active = true;
    httpClient
      .get<RestrictionOption[]>("/talksessions/restrictions")
      .then((res) => active && setOptions(res.data ?? []))
      .catch(() => active && setOptionsError(true));
    return () => {
      active = false;
    };
  }, []);

  useEffect(() => {
    if (session) form.setFieldsValue(toEditValues(session));
  }, [session, form]);

  const save = (values: TalkSessionEditValues) => {
    if (!session) return;
    mutate(
      {
        url: `${env.apiUrl}/v1/manage/talksessions/${encodeURIComponent(session.talkSessionID)}`,
        method: "put",
        values: toEditRequest(values, session),
        successNotification: { type: "success", message: "セッションを更新しました" },
        errorNotification: (error) => ({
          type: "error",
          message: "セッションを更新できませんでした",
          description: error?.message,
        }),
      },
      {
        onSuccess: () => {
          invalidate({ resource: "talksessions", invalidates: ["list", "detail"], id });
          navigate("/talksessions");
        },
      },
    );
  };

  const onFinish = (values: TalkSessionEditValues) => {
    if (!session || !options) return;
    const errors = validateEditValues(values, session, options);
    const fields = Object.entries(errors).map(([name, message]) => ({
      name: name as keyof TalkSessionEditValues,
      errors: [message],
    }));
    if (fields.length > 0) {
      form.setFields(fields);
      return;
    }
    if (needsRestrictionWarning(session, values)) {
      modal.confirm({
        title: "参加制限を追加しますか？",
        content: `このセッションにはすでに${session.participantCount}人の参加者がいます。条件を満たさない参加者は、これ以降意見の投稿や投票ができなくなります。`,
        okText: "追加して保存",
        okButtonProps: { danger: true },
        cancelText: "キャンセル",
        onOk: () => save(values),
      });
      return;
    }
    save(values);
  };

  const notFound = axios.isAxiosError(query.error) && query.error.response?.status === 404;

  return (
    <Edit
      title="セッションの編集"
      isLoading={query.isLoading}
      canDelete={false}
      headerButtons={() => null}
      footerButtons={() => (
        <Button
          type="primary"
          loading={mutation.isPending}
          disabled={!session || !options}
          onClick={() => form.submit()}
        >
          保存する
        </Button>
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
          <Space direction="vertical" size="large" style={{ width: "100%" }}>
            <Descriptions bordered column={1} size="small">
              <Descriptions.Item label="セッションID">
                <Typography.Text copyable>{session.talkSessionID}</Typography.Text>
              </Descriptions.Item>
              <Descriptions.Item label="オーナー">
                {session.owner.displayName || "未設定"}
                {session.owner.userID && (
                  <>
                    {"（"}
                    <Link to={`/users/show/${encodeURIComponent(session.owner.userID)}`}>
                      ユーザー詳細
                    </Link>
                    {"）"}
                  </>
                )}
              </Descriptions.Item>
              <Descriptions.Item label="参加者数">{session.participantCount}人</Descriptions.Item>
            </Descriptions>

            <Form form={form} layout="vertical" onFinish={onFinish} disabled={mutation.isPending}>
              <Form.Item label="テーマ" name="theme" extra="100文字以内">
                <Input />
              </Form.Item>
              <Row gutter={16}>
                <Col xs={24} lg={12}>
                  <Form.Item
                    label="ストーリー（本文）"
                    name="description"
                    extra="HTMLのまま編集します。40000文字以内"
                  >
                    <Input.TextArea rows={14} style={{ fontFamily: "monospace" }} />
                  </Form.Item>
                </Col>
                <Col xs={24} lg={12}>
                  <Form.Item
                    label="プレビュー"
                    extra="web と同じ条件で表示します。使えないタグ（見出しなど）は文字のまま表示されます"
                  >
                    {/* 利用者が書いたHTMLなので、スクリプトを動かさない sandbox 付き iframe で表示する */}
                    <iframe
                      title="本文のプレビュー"
                      sandbox="allow-popups allow-popups-to-escape-sandbox"
                      srcDoc={buildPreviewDocument(description)}
                      style={{
                        width: "100%",
                        height: 318,
                        border: "1px solid #d9d9d9",
                        borderRadius: 6,
                        background: "#fff",
                      }}
                    />
                  </Form.Item>
                </Col>
              </Row>
              <Form.Item
                label="募集の終了日時"
                name="scheduledEndTime"
                extra="変更する場合は、現在より後の日時にしてください"
              >
                <Input type="datetime-local" style={{ maxWidth: 260 }} />
              </Form.Item>
              <Space wrap size="large">
                <Form.Item label="都道府県" name="prefecture">
                  <Input placeholder="未設定" />
                </Form.Item>
                <Form.Item label="市区町村" name="city">
                  <Input placeholder="未設定" />
                </Form.Item>
              </Space>

              <Typography.Title level={5}>表示</Typography.Title>
              <Form.Item label="ホームの一覧に表示する" name="hideTop" {...invert}>
                <Switch checkedChildren="表示する" unCheckedChildren="表示しない" />
              </Form.Item>
              <Form.Item label="レポートを公開する" name="hideReport" {...invert}>
                <Switch checkedChildren="公開する" unCheckedChildren="公開しない" />
              </Form.Item>
              <Form.Item
                label="意見の分析を行う"
                name="disableAnalysis"
                extra="無効から有効に戻すと、すぐに分析を始めます"
                {...invert}
              >
                <Switch checkedChildren="行う" unCheckedChildren="行わない" />
              </Form.Item>

              <Typography.Title level={5}>参加制限</Typography.Title>
              {optionsError ? (
                <Alert type="error" showIcon message="参加制限の一覧を取得できませんでした。画面を開き直してください。" />
              ) : (
                <Form.Item
                  name="restrictions"
                  extra={
                    session.participantCount > 0
                      ? `参加者が${session.participantCount}人います。制限を追加すると、条件を満たさない参加者は投稿や投票ができなくなります。`
                      : "チェックした情報を登録したユーザーだけが参加できます"
                  }
                >
                  <Checkbox.Group
                    options={(options ?? []).map((o) => ({ label: o.description, value: o.key }))}
                  />
                </Form.Item>
              )}
            </Form>
          </Space>
        )
      )}
    </Edit>
  );
}
