import { useCustomMutation } from "@refinedev/core";
import { App, Button, Popconfirm } from "antd";
import axios from "axios";
import { env } from "../../config/env";
import type { ManagedUser } from "./types";

export function ResendUserInvitation({
  user,
  onSuccess,
}: { user: ManagedUser; onSuccess: () => void }) {
  const { message } = App.useApp();
  const { mutate, mutation } = useCustomMutation<{
    userID: string;
    invitationSent: boolean;
  }>();
  const resend = () =>
    mutate(
      {
        url: `${env.apiUrl}/v1/manage/users/${encodeURIComponent(user.userID)}/invitation`,
        method: "post",
        values: {},
        successNotification: false,
        errorNotification: false,
      },
      {
        onSuccess: ({ data }) => {
          if (data.invitationSent)
            message.success("パスワード設定メールを再送しました");
          else
            message.warning(
              "メールを送信できませんでした。60秒以上待ってから再試行してください。",
              8,
            );
          onSuccess();
        },
        onError: (error) => {
          message.error(
            axios.isAxiosError(error) && error.response?.status === 429
              ? "再送は60秒以上待ってから行ってください"
              : "再送できませんでした。パスワードが設定済みでないか、接続状態を確認してください。",
          );
          onSuccess();
        },
      },
    );
  return (
    <Popconfirm
      title="パスワード設定メールを再送しますか？"
      description="以前の設定リンクは無効になります。"
      onConfirm={resend}
      okText="再送する"
      cancelText="キャンセル"
    >
      <Button loading={mutation.isPending} disabled={!!user.withdrawalDate}>
        設定メールを再送
      </Button>
    </Popconfirm>
  );
}
