import { useQuery } from "@tanstack/react-query";
import { useEffect, useRef, useState } from "react";
import * as d3 from "d3";
import { Card } from "./Card";
import { useNotification } from '@/contexts/NotificationContext';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer
} from 'recharts';

interface UserStats {
  date: string;
  userCount: number;
  uniqueActionUserCount: number;
  talkSessionCount: number;
}

type TimeRange = "daily" | "weekly";

export const UserStatsGraph = () => {
  const [timeRange, setTimeRange] = useState<TimeRange>("daily");

  const { data: userStats, isLoading: isUserStatsLoading } = useQuery<UserStats[], Error>({
    queryKey: ['userStats', timeRange],
    queryFn: async () => {
      const response = await fetch(`/v1/manage/users/stats/list?range=${timeRange}&offset=0&limit=10`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });
      if (!response.ok) {
        throw new Error('ユーザー統計の取得に失敗しました');
      }
      const data = await response.json() as UserStats[];
      return data.reverse();
    },
  });

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr);
    return `${date.getMonth() + 1}/${date.getDate()}`;
  };

  const CustomTooltip = ({ active, payload, label }: any) => {
    if (active && payload && payload.length) {
      return (
        <div className="bg-white p-4 rounded-lg shadow-lg border border-gray-200">
          <p className="font-medium text-gray-900">{formatDate(label)}</p>
          {payload.map((item: any, index: number) => (
            <p key={index} className="text-gray-600">
              {item.name}: {item.value}
            </p>
          ))}
        </div>
      );
    }
    return null;
  };

  if (isUserStatsLoading) {
    return (
      <Card
        title="ユーザー統計"
        headerRight={
          <div className="flex items-center space-x-2 px-4">
            <span className="text-sm text-gray-500">期間:</span>
            <select
              className="border rounded px-2 py-1 text-sm"
              value={timeRange}
              onChange={(e) => setTimeRange(e.target.value as TimeRange)}
            >
              <option value="daily">日次</option>
              <option value="weekly">週次</option>
            </select>
          </div>
        }
      >
        <div className="flex justify-center items-center h-64">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-gray-900"></div>
        </div>
      </Card>
    );
  }

  return (
    <Card
      title="ユーザー統計"
      headerRight={
        <div className="flex items-center space-x-2 px-4">
          <span className="text-sm text-gray-500">期間:</span>
          <select
            className="border rounded px-2 py-1 text-sm"
            value={timeRange}
            onChange={(e) => setTimeRange(e.target.value as TimeRange)}
          >
            <option value="daily">日次</option>
            <option value="weekly">週次</option>
          </select>
        </div>
      }
    >
      <div className="h-[500px] w-full">
        <ResponsiveContainer width="100%" height="100%">
          <LineChart
            data={userStats}
            margin={{
              top: 20,
              right: 30,
              left: 20,
              bottom: 60,
            }}
          >
            <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
            <XAxis
              dataKey="date"
              tickFormatter={formatDate}
              angle={-45}
              textAnchor="end"
              height={60}
              tick={{ fill: '#666', fontSize: 12 }}
            />
            <YAxis
              tick={{ fill: '#666', fontSize: 12 }}
            />
            <Tooltip content={<CustomTooltip />} />
            <Legend
              verticalAlign="top"
              height={36}
              formatter={(value) => (
                <span className="text-sm text-gray-600">{value}</span>
              )}
            />
            <Line
              type="monotone"
              dataKey="uniqueActionUserCount"
              name="アクティブユーザー数"
              stroke="#2196F3"
              strokeWidth={2}
              dot={{ fill: '#2196F3', r: 4 }}
              activeDot={{ r: 6 }}
            />
          </LineChart>
        </ResponsiveContainer>
      </div>
    </Card>
  );
};

export const UserStatsTotal = () => {
  const { showNotification } = useNotification();

  const { data, isLoading, error } = useQuery({
    queryKey: ['userStats'],
    queryFn: async () => {
      const response = await fetch('/v1/manage/users/stats/total', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        throw new Error('ユーザー統計の取得に失敗しました');
      }

      return await response.json() as UserStats;
    },
  });

  if (isLoading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {[...Array(3)].map((_, i) => (
          <div key={i} className="bg-white rounded-xl shadow-sm p-6 animate-pulse">
            <div className="h-4 bg-gray-200 rounded w-1/3 mb-4"></div>
            <div className="h-8 bg-gray-200 rounded w-1/2"></div>
          </div>
        ))}
      </div>
    );
  }

  if (error) {
    showNotification('ユーザー統計の取得に失敗しました', 'error');
    return null;
  }

  const stats = [
    {
      title: '総ユーザー数',
      value: data?.userCount,
      icon: (
        <svg className="w-6 h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
        </svg>
      ),
    },
    {
      title: 'アクティブユーザー数',
      value: data?.uniqueActionUserCount,
      icon: (
        <svg className="w-6 h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      ),
    },
    {
      title: 'セッション数',
      value: data?.talkSessionCount,
      icon: (
        <svg className="w-6 h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      ),
    }
  ];

  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
      {stats.map((stat, index) => (
        <div
          key={index}
          className="bg-white rounded-xl shadow-sm border border-gray-100 p-6 transition-all duration-300"
        >
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-sm font-medium text-gray-600">{stat.title}</h3>
            <div className="p-2 rounded-lg bg-gray-50">
              {stat.icon}
            </div>
          </div>
          <div className="flex items-baseline">
            <span className="text-3xl font-bold text-gray-900">{stat.value}</span>
          </div>
        </div>
      ))}
    </div>
  );
};
