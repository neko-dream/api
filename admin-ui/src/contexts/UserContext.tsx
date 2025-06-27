import { createContext, useContext, useEffect, useState } from 'react';
import { useNavigate, useRouterState } from '@tanstack/react-router';
import { authService } from '@/services/auth';

interface UserContextType {
  currentUser?: {
    id: string;
    name?: string;
    email?: string;
    displayID?: string;
    organizationCode?: string;
    organizationID?: string;
    organizationRole?: string;
  };
  isLoading: boolean;
  logout: () => Promise<void>;
}

const UserContext = createContext<UserContextType | null>(null);

export const UserProvider = ({ children }: { children: React.ReactNode }) => {
  const [currentUser, setCurrentUser] = useState<UserContextType['currentUser']>();
  const [isLoading, setIsLoading] = useState(true);
  const navigate = useNavigate();
  const router = useRouterState();
  const isLoginPage = router.location.pathname === '/login';

  useEffect(() => {
    const checkAuth = async () => {
      // Skip auth check on login page
      if (isLoginPage) {
        setIsLoading(false);
        return;
      }
      setIsLoading(true);
      try {
        const tokenInfo = await authService.getTokenInfo();
        
        if (!tokenInfo) {
          // Not authenticated
          navigate({ to: '/login' });
          return;
        }

        // Check if user has organization parameters
        if (!tokenInfo.organizationCode || !tokenInfo.organizationID) {
          // User is authenticated but not associated with an organization
          navigate({ to: '/login' });
          return;
        }

        // Set user data
        setCurrentUser({
          id: tokenInfo.sub,
          name: tokenInfo.displayName,
          displayID: tokenInfo.displayID,
          organizationCode: tokenInfo.organizationCode,
          organizationID: tokenInfo.organizationID,
          organizationRole: tokenInfo.organizationRole,
        });
      } catch (error) {
        console.error('Auth check failed:', error);
        navigate({ to: '/login' });
      } finally {
        setIsLoading(false);
      }
    };

    checkAuth();
  }, [navigate, isLoginPage]);

  const logout = async () => {
    await authService.logout();
    setCurrentUser(undefined);
    navigate({ to: '/login' });
  };

  const value: UserContextType = {
    currentUser,
    isLoading,
    logout,
  };

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-gray-500">Loading...</div>
      </div>
    );
  }

  return <UserContext.Provider value={value}>{children}</UserContext.Provider>;
};

export const useUser = () => {
  const context = useContext(UserContext);
  if (!context) {
    throw new Error('useUser must be used within a UserProvider');
  }
  return context;
};