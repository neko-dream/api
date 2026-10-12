import { useList } from "@refinedev/core";
import { List } from "@refinedev/antd";
import { Alert, Button, Input, Space, Table, Tag } from "antd";
import { Link, useSearchParams } from "react-router";
import type { ManagedUser } from "./types";
import { formatUserDate } from "./types";

export function UserList() {
  const [params, setParams] = useSearchParams();
  const candidate = Number(params.get("page") ?? 1);
  const page =
    Number.isSafeInteger(candidate) && candidate >= 1 && candidate <= 100000000 ? candidate : 1;
  const search = params.get("search") ?? "";
  const { query, result } = useList<ManagedUser>({
    resource: "users",
    pagination: { currentPage: page, pageSize: 20 },
    filters: [{ field: "search", operator: "contains", value: search }],
  });
  return (
    <List title="ユーザー" canCreate={false}>
      <Space direction="vertical" size="middle" style={{ width: "100%" }}>
        <Input.Search
          key={search}
          defaultValue={search}
          maxLength={200}
          allowClear
          placeholder="表示名・表示ID・ユーザーIDで検索"
          aria-label="ユーザーを検索"
          enterButton="検索"
          onSearch={(value) => setParams({ search: value.trim(), page: "1" })}
        />
        {query.isError && (
          <Alert
            type="error"
            showIcon
            message="ユーザー一覧を取得できませんでした"
            action={<Button onClick={() => query.refetch()}>再試行</Button>}
          />
        )}
        <Table<ManagedUser>
          rowKey="userID"
          loading={query.isFetching}
          dataSource={query.isError ? [] : result.data}
          scroll={{ x: 1000 }}
          locale={{
            emptyText: query.isError ? "取得に失敗しました" : "該当するユーザーはいません",
          }}
          pagination={{
            current: page,
            pageSize: 20,
            total: result.total,
            showSizeChanger: false,
            showTotal: (total) => `全${total}件`,
            onChange: (next) => setParams({ search, page: String(next) }),
          }}
          columns={[
            {
              title: "表示名",
              dataIndex: "displayName",
              render: (value, user) => (
                <Link to={`/users/show/${user.userID}`}>{value || "未設定"}</Link>
              ),
            },
            { title: "表示ID", dataIndex: "displayID", render: (value) => value || "—" },
            { title: "メールアドレス", dataIndex: "email", render: (value) => value || "未登録" },
            {
              title: "状態",
              render: (_, user) =>
                user.withdrawalDate ? <Tag>退会済み</Tag> : <Tag color="green">有効</Tag>,
            },
            { title: "登録日時", dataIndex: "createdAt", render: formatUserDate },
            { title: "ユーザーID", dataIndex: "userID" },
          ]}
        />
      </Space>
    </List>
  );
}
