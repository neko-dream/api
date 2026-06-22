import { useEffect, useState } from "react";
import { Table, Tag, Typography } from "antd";

import { httpClient } from "../../providers";

interface TalkSessionStat {
  talkSessionID: string;
  theme: string;
  owner: { displayName?: string };
  scheduledEndTime: string;
  opinionCount: number;
  voteCount: number;
  hidden: boolean;
}

const PAGE_SIZE = 20;

export const TalkSessionList = () => {
  const [data, setData] = useState<TalkSessionStat[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(1);

  useEffect(() => {
    let active = true;
    setLoading(true);
    const offset = (page - 1) * PAGE_SIZE;
    httpClient
      .get<{ talkSessionStats: TalkSessionStat[]; totalCount: number }>(
        `/v1/manage/talksessions/list?limit=${PAGE_SIZE}&offset=${offset}`,
      )
      .then((res) => {
        if (!active) return;
        setData(res.data.talkSessionStats ?? []);
        setTotal(res.data.totalCount ?? 0);
      })
      .finally(() => {
        if (active) setLoading(false);
      });
    return () => {
      active = false;
    };
  }, [page]);

  return (
    <div style={{ padding: 24 }}>
      <Typography.Title level={4}>トークセッション</Typography.Title>
      <Table
        rowKey="talkSessionID"
        loading={loading}
        dataSource={data}
        pagination={{
          current: page,
          pageSize: PAGE_SIZE,
          total,
          onChange: setPage,
          showSizeChanger: false,
        }}
        columns={[
          { title: "テーマ", dataIndex: "theme" },
          { title: "オーナー", dataIndex: ["owner", "displayName"] },
          { title: "意見数", dataIndex: "opinionCount", width: 90 },
          { title: "投票数", dataIndex: "voteCount", width: 90 },
          {
            title: "終了予定",
            dataIndex: "scheduledEndTime",
            render: (v: string) => (v ? new Date(v).toLocaleString("ja-JP") : ""),
          },
          {
            title: "表示",
            dataIndex: "hidden",
            width: 90,
            render: (hidden: boolean) =>
              hidden ? <Tag color="red">非表示</Tag> : <Tag color="green">表示</Tag>,
          },
        ]}
      />
    </div>
  );
};
