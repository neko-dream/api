import { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { Button, Space, Table, Typography } from "antd";

import { httpClient } from "../../providers";

interface AdminUserItem {
  id: string;
  email: string;
  displayName: string;
  role: number;
}

export const AdminUserList = () => {
  const [data, setData] = useState<AdminUserItem[]>([]);
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    let active = true;
    setLoading(true);
    httpClient
      .get<{ adminUsers: AdminUserItem[] }>("/v1/admin/admin-users")
      .then((res) => {
        if (active) setData(res.data.adminUsers ?? []);
      })
      .finally(() => {
        if (active) setLoading(false);
      });
    return () => {
      active = false;
    };
  }, []);

  return (
    <div style={{ padding: 24 }}>
      <Space style={{ marginBottom: 16, justifyContent: "space-between", width: "100%" }}>
        <Typography.Title level={4} style={{ margin: 0 }}>
          管理者アカウント
        </Typography.Title>
        <Button type="primary" onClick={() => navigate("/admin-users/create")}>
          新規作成
        </Button>
      </Space>
      <Table
        rowKey="id"
        loading={loading}
        dataSource={data}
        columns={[
          { title: "メールアドレス", dataIndex: "email" },
          { title: "表示名", dataIndex: "displayName" },
          { title: "ロール", dataIndex: "role" },
        ]}
      />
    </div>
  );
};
