import { Authenticated, Refine } from "@refinedev/core";
import { RefineKbar, RefineKbarProvider } from "@refinedev/kbar";
import {
  ErrorComponent,
  ThemedLayout,
  useNotificationProvider,
} from "@refinedev/antd";
import "@refinedev/antd/dist/reset.css";
import routerProvider, {
  CatchAllNavigate,
  NavigateToResource,
} from "@refinedev/react-router";
import { GoogleOAuthProvider } from "@react-oauth/google";
import { App as AntdApp, ConfigProvider, Typography } from "antd";
import { BrowserRouter, Outlet, Route, Routes } from "react-router";

import { env } from "./config/env";
import { LoginPage } from "./pages/login";
import { AdminUserCreate } from "./pages/admin-users/create";
import { AdminUserList } from "./pages/admin-users/list";
import { authProvider, dataProvider } from "./providers";

function App() {
  return (
    <GoogleOAuthProvider clientId={env.googleClientId}>
      <BrowserRouter>
      <RefineKbarProvider>
        <ConfigProvider>
          <AntdApp>
            <Refine
              dataProvider={dataProvider}
              authProvider={authProvider}
              routerProvider={routerProvider}
              notificationProvider={useNotificationProvider}
              resources={[
                {
                  name: "talksessions",
                  list: "/talksessions",
                  meta: { label: "トークセッション" },
                },
                {
                  name: "admin-users",
                  list: "/admin-users",
                  create: "/admin-users/create",
                  meta: { label: "管理者" },
                },
              ]}
              options={{
                syncWithLocation: true,
                warnWhenUnsavedChanges: true,
                title: { text: "Kotohiro 管理画面" },
              }}
            >
              <Routes>
                <Route
                  element={
                    <Authenticated
                      key="authenticated-routes"
                      fallback={<CatchAllNavigate to="/login" />}
                    >
                      <ThemedLayout>
                        <Outlet />
                      </ThemedLayout>
                    </Authenticated>
                  }
                >
                  <Route
                    index
                    element={<NavigateToResource resource="talksessions" />}
                  />
                  <Route
                    path="/talksessions"
                    element={
                      <Typography.Paragraph>
                        TODO: トークセッション一覧（ベースのみ。内容は今後実装）
                      </Typography.Paragraph>
                    }
                  />
                  <Route path="/admin-users" element={<AdminUserList />} />
                  <Route path="/admin-users/create" element={<AdminUserCreate />} />
                  <Route path="*" element={<ErrorComponent />} />
                </Route>

                <Route
                  element={
                    <Authenticated key="auth-pages" fallback={<Outlet />}>
                      <NavigateToResource resource="talksessions" />
                    </Authenticated>
                  }
                >
                  <Route path="/login" element={<LoginPage />} />
                </Route>
              </Routes>
              <RefineKbar />
            </Refine>
          </AntdApp>
        </ConfigProvider>
      </RefineKbarProvider>
      </BrowserRouter>
    </GoogleOAuthProvider>
  );
}

export default App;
