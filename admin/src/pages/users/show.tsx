import { useCustomMutation, useInvalidate, useShow } from "@refinedev/core";
import { Show } from "@refinedev/antd";
import {
  Alert,
  Button,
  Descriptions,
  Popconfirm,
  Result,
  Space,
  Tag,
  Typography,
} from "antd";
import axios from "axios";
import { Link } from "react-router";
import type { ManagedUser } from "./types";
import {
  canGrantSessionCreation,
  formatUserDate,
  providerLabel,
  sessionCreationPermissionLabel,
} from "./types";
import { env } from "../../config/env";
import { ResendUserInvitation } from "./resend-invitation";

export function UserShow() {
  const { query, result } = useShow<ManagedUser>({ resource: "users" });
  const user = result;
  const { mutate, mutation } = useCustomMutation();
  const invalidate = useInvalidate();
  const setCreationPermission = () => {
    if (!user) return;
    const allowed = !user.permissions.canCreateTalkSession;
    mutate(
      {
        url: `${env.apiUrl}/v1/manage/users/${encodeURIComponent(user.userID)}/permissions`,
        method: "put",
        values: { canCreateTalkSession: allowed },
        successNotification: {
          type: "success",
          message: allowed
            ? "セッション作成を許可しました"
            : "セッション作成の許可を解除しました",
        },
        errorNotification: {
          type: "error",
          message: "作成権限を変更できませんでした",
          description:
            "ゲスト・退会済みユーザーには許可できません。接続状態を確認して再試行してください。",
        },
      },
      {
        onSuccess: () => {
          query.refetch();
          invalidate({ resource: "users", invalidates: ["list"] });
        },
      },
    );
  };
  const notFound =
    axios.isAxiosError(query.error) && query.error.response?.status === 404;
  return (
    <Show
      title="ユーザー詳細"
      isLoading={query.isLoading}
      canEdit={false}
      canDelete={false}
    >
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
          <Descriptions bordered column={1}>
            <Descriptions.Item label="表示名">
              {user.displayName || "未設定"}
            </Descriptions.Item>
            <Descriptions.Item label="表示ID">
              {user.displayID || "未設定"}
            </Descriptions.Item>
            <Descriptions.Item label="ユーザーID">
              <Typography.Text copyable>{user.userID}</Typography.Text>
            </Descriptions.Item>
            <Descriptions.Item label="メールアドレス">
              {user.email || "未登録"}
            </Descriptions.Item>
            <Descriptions.Item label="メールの確認">
              {user.email ? (user.emailVerified ? "確認済み" : "未確認") : "—"}
            </Descriptions.Item>
            {user.passwordSetupRequired && (
              <Descriptions.Item label="パスワード設定">
                本人による設定待ち
              </Descriptions.Item>
            )}
            {user.passwordSetupRequired && (
              <Descriptions.Item label="パスワード設定メール">
                <ResendUserInvitation
                  user={user}
                  onSuccess={() => {
                    query.refetch();
                    invalidate({ resource: "users", invalidates: ["list"] });
                  }}
                />
              </Descriptions.Item>
            )}
            <Descriptions.Item label="ログイン方法">
              <Space wrap>
                {user.providers.length
                  ? user.providers.map((provider) => (
                      <Tag key={provider}>{providerLabel(provider)}</Tag>
                    ))
                  : "未登録"}
              </Space>
            </Descriptions.Item>
            <Descriptions.Item label="セッション作成の許可">
              <Space wrap>
                <Tag
                  color={
                    user.permissions.canCreateTalkSession &&
                    canGrantSessionCreation(user)
                      ? "green"
                      : undefined
                  }
                >
                  {sessionCreationPermissionLabel(user)}
                </Tag>
                <Popconfirm
                  title={
                    user.permissions.canCreateTalkSession
                      ? "セッション作成の許可を解除しますか？"
                      : "セッション作成を許可しますか？"
                  }
                  description="組織への所属や、既存セッションの編集権限は変更しません。"
                  onConfirm={setCreationPermission}
                  okText="変更する"
                  cancelText="キャンセル"
                >
                  <Button
                    loading={mutation.isPending}
                    danger={user.permissions.canCreateTalkSession}
                    disabled={
                      query.isFetching ||
                      (!user.permissions.canCreateTalkSession &&
                        !canGrantSessionCreation(user))
                    }
                  >
                    {user.permissions.canCreateTalkSession
                      ? "許可を解除"
                      : "作成を許可"}
                  </Button>
                </Popconfirm>
              </Space>
            </Descriptions.Item>
            <Descriptions.Item label="登録日時">
              {formatUserDate(user.createdAt)}
            </Descriptions.Item>
            <Descriptions.Item label="更新日時">
              {formatUserDate(user.updatedAt)}
            </Descriptions.Item>
            <Descriptions.Item label="状態">
              {user.withdrawalDate ? (
                <Tag>退会済み</Tag>
              ) : (
                <Tag color="green">有効</Tag>
              )}
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
