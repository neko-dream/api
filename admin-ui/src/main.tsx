import { StrictMode } from 'react'
import ReactDOM from 'react-dom/client'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { NotificationProvider } from './contexts/NotificationContext';
import { UserProvider } from './contexts/UserContext';
import { createRootRoute, createRoute, createRouter, RouterProvider } from '@tanstack/react-router';
import { routeTree } from './routeTree.gen'
import './index.css'

const queryClient = new QueryClient();

const router = createRouter({ routeTree })


declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}

const rootElement = document.getElementById('root')!
if (!rootElement.innerHTML) {
  const root = ReactDOM.createRoot(rootElement)
  root.render(
    <StrictMode>
      <QueryClientProvider client={queryClient}>
        <NotificationProvider>
          <UserProvider>
            <RouterProvider router={router} />
          </UserProvider>
        </NotificationProvider>
      </QueryClientProvider>
    </StrictMode>
  )
}
