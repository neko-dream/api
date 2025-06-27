import { createRootRoute, Outlet, useRouterState } from '@tanstack/react-router'
import { Layout } from '../components/Layout'
import { UserProvider } from '@/contexts/UserContext'
import { NotificationProvider } from '@/contexts/NotificationContext'

export const Route = createRootRoute({
  component: RootComponent,
});

function RootComponent() {
  const router = useRouterState();
  const isLoginPage = router.location.pathname === '/login';

  return (
    <UserProvider>
      <NotificationProvider>
        {isLoginPage ? (
          <Outlet />
        ) : (
          <Layout>
            <div className="container mx-auto px-4 py-6">
              <Outlet />
            </div>
          </Layout>
        )}
      </NotificationProvider>
    </UserProvider>
  );
}
