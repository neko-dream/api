import { createFileRoute } from '@tanstack/react-router'
import { UserStatsGraph, UserStatsTotal } from '@/components/UserStats'
export const Route = createFileRoute('/')({
  component: Index,
})

function Index() {
  return (
    <div>
      <UserStatsGraph />
      <UserStatsTotal />
    </div>
  );
}
