import { useLogin } from "@refinedev/core";
import { GoogleLogin } from "@react-oauth/google";
import { Alert, Card, Layout, Space, Typography } from "antd";
import { env } from "../config/env";

export const LoginPage = () => {
  const { mutate: login } = useLogin();

  return (
    <Layout
      style={{
        minHeight: "100vh",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        background: "#f0f2f5",
      }}
    >
      <Card style={{ width: 360 }}>
        <Space direction="vertical" align="center" style={{ width: "100%" }} size="large">
          <Typography.Title level={3} style={{ textAlign: "center", margin: 0 }}>
            ことひろ 管理画面
          </Typography.Title>
          <Typography.Text type="secondary">Google アカウントでログイン</Typography.Text>
          {env.googleClientId ? (
            <GoogleLogin
              onSuccess={(credentialResponse) =>
                login({ credential: credentialResponse.credential })
              }
              onError={() => {
                // GoogleLogin 側のエラー。authProvider 側で資格情報なしとして扱われる。
              }}
            />
          ) : (
            <Alert
              type="error"
              message="Googleログインが未設定です"
              description="VITE_GOOGLE_CLIENT_IDを設定してください"
            />
          )}
        </Space>
      </Card>
    </Layout>
  );
};
