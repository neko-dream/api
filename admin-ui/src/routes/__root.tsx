import { createRootRoute, Outlet } from '@tanstack/react-router'
import { Layout } from '../components/Layout'
import { UserProvider } from '@/contexts/UserContext'
import { NotificationProvider } from '@/contexts/NotificationContext'

export const Route = createRootRoute({
  component: () => (
    <UserProvider >
      <NotificationProvider>
        <Layout>
          <div className="container mx-auto px-4 py-6">
            <Outlet />
          </div>
        </Layout>
      </NotificationProvider>
    </UserProvider>
  ),
})
