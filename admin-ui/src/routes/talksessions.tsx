import { createFileRoute } from '@tanstack/react-router';
import { SessionList } from '../components/SessionList';

export const Route = createFileRoute('/talksessions')({
  component: Talksessions,
  validateSearch: (search: Record<string, unknown>): { page: number } => ({
    page: Number(search.page) || 1,
  }),
});

function Talksessions() {
  const { page } = Route.useSearch();
  return <SessionList initialPage={page} />;
}
