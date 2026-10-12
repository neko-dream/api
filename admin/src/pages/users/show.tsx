import { useShow } from "@refinedev/core";
import { Show } from "@refinedev/antd";
import { Alert, Button, Descriptions, Result, Space, Tag, Typography } from "antd";
import axios from "axios";
import { Link } from "react-router";
import type { ManagedUser } from "./types";
import { formatUserDate, providerLabel } from "./types";

export function UserShow() {
  const { query, result } = useShow<ManagedUser>({ resource: "users" });
  const user = result;
  const notFound = axios.isAxiosError(query.error) && query.error.response?.status === 404;
  return (
    <Show title="ユーザー詳細" isLoading={query.isLoading} canEdit={false} canDelete={false}>
      {query.isError ? (
        notFound ? (
          <Result
            status="404"
            title="ユーザーが見つかりません"
            extra={<Link to="/users">ユーザー一覧へ戻る</Link>}
          />
        ) : (
          <Alert
            type="error"
            showIcon
            message="ユーザー詳細を取得できませんでした"
            action={<Button onClick={() => query.refetch()}>再試行</Button>}
          />
        )
      ) : (
        user && (
          <Descriptions bordered column={{ xs: 1, sm: 1, md: 2 }}>
            <Descriptions.Item label="表示名">{user.displayName || "未設定"}</Descriptions.Item>
            <Descriptions.Item label="表示ID">{user.displayID || "未設定"}</Descriptions.Item>
            <Descriptions.Item label="ユーザーID" span={2}>
              <Typography.Text copyable>{user.userID}</Typography.Text>
            </Descriptions.Item>
            <Descriptions.Item label="メールアドレス">{user.email || "未登録"}</Descriptions.Item>
            <Descriptions.Item label="メールの確認">
              {user.email ? (user.emailVerified ? "確認済み" : "未確認") : "—"}
            </Descriptions.Item>
            <Descriptions.Item label="ログイン方法" span={2}>
              <Space wrap>
                {user.providers.length
                  ? user.providers.map((provider) => (
                      <Tag key={provider}>{providerLabel(provider)}</Tag>
                    ))
                  : "未登録"}
              </Space>
            </Descriptions.Item>
            <Descriptions.Item label="登録日時">{formatUserDate(user.createdAt)}</Descriptions.Item>
            <Descriptions.Item label="更新日時">{formatUserDate(user.updatedAt)}</Descriptions.Item>
            <Descriptions.Item label="状態">
              {user.withdrawalDate ? <Tag>退会済み</Tag> : <Tag color="green">有効</Tag>}
            </Descriptions.Item>
            <Descriptions.Item label="退会日時">
              {formatUserDate(user.withdrawalDate)}
            </Descriptions.Item>
          </Descriptions>
        )
      )}
    </Show>
  );
}
