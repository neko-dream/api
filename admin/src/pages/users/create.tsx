import { Create, useForm } from "@refinedev/antd";
import type { HttpError } from "@refinedev/core";
import { Alert, App, Form, Input } from "antd";
import axios from "axios";
import { useNavigate } from "react-router";
import {
  trimUserField,
  validateUserDisplayName,
  validateUserEmail,
} from "./validation";

interface InviteUserValues {
  displayName: string;
  email: string;
}
interface UserInvitation {
  id: string;
  userID: string;
  invitationSent: boolean;
}

export function UserCreate() {
  const navigate = useNavigate();
  const { message } = App.useApp();
  const { formProps, saveButtonProps } = useForm<
    UserInvitation,
    HttpError,
    InviteUserValues
  >({
    resource: "users",
    action: "create",
    redirect: false,
    successNotification: false,
    errorNotification: false,
    onMutationSuccess: ({ data }) => {
      if (data.invitationSent) {
        message.success("ユーザーを作成し、パスワード設定メールを送信しました");
      } else {
        message.warning(
          "ユーザーは作成されましたが、メールを送信できませんでした。詳細画面から再送してください。",
          10,
        );
      }
      navigate(`/users/show/${encodeURIComponent(data.userID)}`);
    },
    onMutationError: (error) => {
      message.error(
        axios.isAxiosError(error) && error.response?.status === 409
          ? "このメールアドレスのユーザーは登録済みです。ユーザー詳細から設定メールを再送してください。"
          : "ユーザーを作成できませんでした。入力内容と接続状態を確認してください。",
      );
    },
  });

  return (
    <Create
      title="ユーザー作成"
      saveButtonProps={{
        ...saveButtonProps,
        children: "作成して設定メールを送信",
      }}
    >
      <div style={{ maxWidth: 480 }}>
        <Alert
          type="info"
          showIcon
          message="本人がメールからパスワードを設定します"
          description="管理画面の管理者や組織への所属は追加しません。セッション作成の許可は、作成後のユーザー詳細で設定してください。"
          style={{ marginBottom: 24 }}
        />
        <Form
          {...formProps}
          noValidate
          layout="vertical"
          onFinish={(values: InviteUserValues) =>
            formProps.onFinish?.({
              displayName: trimUserField(values.displayName),
              email: trimUserField(values.email),
            })
          }
        >
          <Form.Item
            label="表示名"
            name="displayName"
            rules={[
              {
                validator: (_, value) => {
                  const error = validateUserDisplayName(value);
                  return error
                    ? Promise.reject(new Error(error))
                    : Promise.resolve();
                },
              },
            ]}
          >
            <Input autoComplete="off" />
          </Form.Item>
          <Form.Item
            label="メールアドレス"
            name="email"
            rules={[
              {
                validator: (_, value) => {
                  const error = validateUserEmail(value);
                  return error
                    ? Promise.reject(new Error(error))
                    : Promise.resolve();
                },
              },
            ]}
          >
            <Input type="email" autoComplete="off" autoCapitalize="none" />
          </Form.Item>
        </Form>
      </div>
    </Create>
  );
}
