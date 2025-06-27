import { useQuery } from '@tanstack/react-query';
import { SessionCard } from './SessionCard';
import { useUser } from '../contexts/UserContext';
import { useNotification } from '../contexts/NotificationContext';
import { Pagination } from './Pagination';
import { useNavigate, useSearch } from '@tanstack/react-router';
import { TalkSessionListResponse } from '@/types/session';
import { Route } from '../routes/talksessions';

const ITEMS_PER_PAGE = 10;

interface SessionListProps {
  initialPage?: number;
}

export const SessionList = ({ initialPage = 1 }: SessionListProps) => {
  const { currentUser, isLoading: isUserLoading } = useUser();
  const navigate = useNavigate({ from: Route.fullPath });
  const search = useSearch({ from: Route.fullPath });

  const { data, isLoading: isSessionsLoading, refetch } = useQuery({
    queryKey: ['sessions', search.page],
    queryFn: async () => {
      const offset = ((search.page || 1) - 1) * ITEMS_PER_PAGE;
      const response = await fetch(`/v1/manage/talksessions/list?limit=${ITEMS_PER_PAGE}&offset=${offset}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
      });

      if (!response.ok) {
        throw new Error('セッションの取得に失敗しました');
      }

      const result = await response.json();
      return result as TalkSessionListResponse;
    },
  });

  const handlePageChange = (page: number) => {
    navigate({
      search: (prev) => ({ ...prev, page }),
    });
  };

  if (isUserLoading || isSessionsLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-gray-900"></div>
      </div>
    );
  }

  if (!currentUser) {
    return null;
  }

  const totalPages = Math.ceil((data?.totalCount || 0) / ITEMS_PER_PAGE);

  return (
    <div className="space-y-4">
      {totalPages > 1 && (
        <div className="mb-4">
          <Pagination
            currentPage={search.page || 1}
            totalPages={totalPages}
            onPageChange={handlePageChange}
          />
        </div>
      )}

      {data?.talkSessionStats && data.talkSessionStats.length > 0 ? (
        data.talkSessionStats.map((session: any) => (
          <SessionCard
            key={session.talkSessionID}
            session={session}
            onRefetch={() => {
              refetch();
            }}
          />
        ))
      ) : (
        <div className="text-center py-8 text-gray-500">
          セッションが見つかりませんでした
        </div>
      )}

      {totalPages > 1 && (
        <div className="mt-8">
          <Pagination
            currentPage={search.page || 1}
            totalPages={totalPages}
            onPageChange={handlePageChange}
          />
        </div>
      )}
    </div>
  );
};
