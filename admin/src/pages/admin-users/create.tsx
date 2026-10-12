import { useState } from "react";
import { useNavigate } from "react-router";
import { App, Button, Card, Form, Input, Typography } from "antd";

import { httpClient } from "../../providers";

interface CreateAdminUserValues {
  email: string;
  displayName?: string;
}

export const AdminUserCreate = () => {
  const navigate = useNavigate();
  const { message } = App.useApp();
  const [submitting, setSubmitting] = useState(false);

  const onFinish = async (values: CreateAdminUserValues) => {
    setSubmitting(true);
    try {
      await httpClient.post("/v1/admin/admin-users", values);
      message.success("管理者を作成しました");
      navigate("/admin-users");
    } catch (e) {
      const status = (e as { response?: { status?: number } })?.response?.status;
      if (status === 409) {
        message.error("既に登録済みのメールアドレスです");
      } else {
        message.error("作成に失敗しました");
      }
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div style={{ padding: 24, maxWidth: 480 }}>
      <Typography.Title level={4}>管理者作成</Typography.Title>
      <Card>
        <Form<CreateAdminUserValues> layout="vertical" onFinish={onFinish}>
          <Form.Item
            label="メールアドレス (Google アカウント)"
            name="email"
            rules={[{ required: true, type: "email" }]}
          >
            <Input placeholder="user@example.com" />
          </Form.Item>
          <Form.Item label="表示名" name="displayName">
            <Input />
          </Form.Item>
          <Button type="primary" htmlType="submit" loading={submitting}>
            作成
          </Button>
        </Form>
      </Card>
    </div>
  );
};
