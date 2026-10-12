import { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { Alert, Button, Space, Table, Typography } from "antd";

import { httpClient } from "../../providers";

interface AdminUserItem {
  id: string;
  email: string;
  displayName: string;
}

export const AdminUserList = () => {
  const [data, setData] = useState<AdminUserItem[]>([]);
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const [error, setError] = useState(false);
  const [retry, setRetry] = useState(0);

  useEffect(() => {
    let active = true;
    setLoading(true);
    setError(false);
    httpClient
      .get<{ adminUsers: AdminUserItem[] }>("/v1/admin/admin-users")
      .then((res) => {
        if (active) setData(res.data.adminUsers ?? []);
      })
      .catch(() => {
        if (active) {
          setData([]);
          setError(true);
        }
      })
      .finally(() => {
        if (active) setLoading(false);
      });
    return () => {
      active = false;
    };
  }, [retry]);

  return (
    <div style={{ padding: 24 }}>
      <Space style={{ marginBottom: 16, justifyContent: "space-between", width: "100%" }}>
        <Typography.Title level={4} style={{ margin: 0 }}>
          管理者
        </Typography.Title>
        <Button type="primary" onClick={() => navigate("/admin-users/create")}>
          新規作成
        </Button>
      </Space>
      {error && (
        <Alert
          type="error"
          showIcon
          message="管理者一覧を取得できませんでした"
          action={<Button onClick={() => setRetry((value) => value + 1)}>再試行</Button>}
        />
      )}
      <Table
        rowKey="id"
        loading={loading}
        dataSource={data}
        columns={[
          { title: "メールアドレス", dataIndex: "email" },
          { title: "表示名", dataIndex: "displayName" },
        ]}
      />
    </div>
  );
};
