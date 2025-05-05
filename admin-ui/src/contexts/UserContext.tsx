import { createContext, useContext, useEffect, useState } from 'react';
import { useQuery } from '@tanstack/react-query';

interface User {
  displayName: string;
  displayID: string;
  Email: string;
  Role: string;
  IconURL?: string;
}

interface UserContextType {
  currentUser: User | null;
  isLoading: boolean;
  error: Error | null;
  refetch: () => Promise<void>;
}

const UserContext = createContext<UserContextType | null>(null);

export const UserProvider = ({ children }: { children: React.ReactNode }) => {
  const [currentUser, setCurrentUser] = useState<User | null>(null);

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['currentUser'],
    queryFn: async () => {
      const response = await fetch('/auth/token/info', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        throw new Error('ユーザー情報の取得に失敗しました');
      }

      const result = await response.json();
      if (result.code) {
        throw new Error('ユーザー情報の取得に失敗しました');
      }
      return result;
    },
    staleTime: 5 * 60 * 1000, // 5分間キャッシュ
  });

  useEffect(() => {
    if (data) {
      setCurrentUser(data);
    }
  }, [data]);

  const value = {
    currentUser,
    isLoading,
    error,
    refetch: async () => {
      await refetch();
    },
  };

  return <UserContext.Provider value={value}>{children}</UserContext.Provider>;
};

export const useUser = () => {
  const context = useContext(UserContext);
  if (!context) {
    throw new Error('useUser must be used within a UserProvider');
  }
  return context;
};
