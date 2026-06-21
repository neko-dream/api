import { useLogin } from "@refinedev/core";
import { Button, Card, Form, Input, Layout, Typography } from "antd";

interface LoginFormValues {
  email: string;
  password: string;
}

export const LoginPage = () => {
  const { mutate: login } = useLogin<LoginFormValues>();

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
        <Typography.Title level={3} style={{ textAlign: "center" }}>
          Kotohiro 管理画面
        </Typography.Title>
        <Form<LoginFormValues> layout="vertical" onFinish={(values) => login(values)}>
          <Form.Item
            label="メールアドレス"
            name="email"
            rules={[{ required: true, type: "email" }]}
          >
            <Input size="large" autoComplete="username" />
          </Form.Item>
          <Form.Item label="パスワード" name="password" rules={[{ required: true }]}>
            <Input.Password size="large" autoComplete="current-password" />
          </Form.Item>
          <Button type="primary" size="large" htmlType="submit" block>
            ログイン
          </Button>
        </Form>
      </Card>
    </Layout>
  );
};
