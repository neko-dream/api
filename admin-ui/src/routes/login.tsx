import React from 'react';
import { createFileRoute, Navigate, useNavigate } from '@tanstack/react-router';
import { useUser } from '@/contexts/UserContext';

export const Route = createFileRoute('/login')({
  component: Login,
});

function Login() {
  const { currentUser } = useUser();

  if (currentUser) {
    return <Navigate to="/" />
  }

  const handleGoogleLogin = () => {
    const redirectUrl = `${window.location.origin}/admin`;
    window.location.href = `/auth/google/login?redirect_url=${encodeURIComponent(redirectUrl)}`;
  };
  return (
    <div className="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
      <div className="sm:mx-auto sm:w-full sm:max-w-md">
        <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
          ことひろ管理画面
        </h2>
        <p className="mt-2 text-center text-sm text-gray-600">
          ログインして管理画面にアクセスしてください
        </p>
      </div>

      <div className="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
        <div className="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10">
          <div className="space-y-6">
            <button
              onClick={handleGoogleLogin}
              className="w-full flex items-center justify-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              Googleでログイン
            </button>

          </div>
        </div>
      </div>
    </div>
  );
}
