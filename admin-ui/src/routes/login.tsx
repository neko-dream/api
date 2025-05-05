import React from 'react';
import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/login')({
  component: Login,
});

function Login() {
  const handleGoogleLogin = () => {
    const redirectUrl = `${window.location.origin}/admin`;
    window.location.href = `/auth/google/login?redirect_url=${encodeURIComponent(redirectUrl)}`;
  };

  const handleLineLogin = () => {
    const redirectUrl = `${window.location.origin}/admin`;
    window.location.href = `/auth/line/login?redirect_url=${encodeURIComponent(redirectUrl)}`;
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
              <svg className="w-5 h-5 mr-2" viewBox="0 0 24 24">
                <path
                  fill="currentColor"
                  d="M12.545,10.239v3.821h5.445c-0.712,2.315-2.647,3.972-5.445,3.972c-3.332,0-6.033-2.701-6.033-6.032s2.701-6.032,6.033-6.032c1.498,0,2.866,0.549,3.921,1.453l2.814-2.814C17.503,2.988,15.139,2,12.545,2C7.021,2,2.543,6.477,2.543,12s4.478,10,10.002,10c8.396,0,10.249-7.85,9.426-11.748L12.545,10.239z"
                />
              </svg>
              Googleでログイン
            </button>

            <button
              onClick={handleLineLogin}
              className="w-full flex items-center justify-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
            >
              <svg className="w-5 h-5 mr-2" viewBox="0 0 24 24">
                <path
                  fill="currentColor"
                  d="M19.365 9.863c.349 0 .63.282.63.631 0 .349-.281.63-.63.63H17.61v1.125h1.755c.349 0 .63.282.63.63 0 .348-.281.63-.63.63h-2.386c-.348 0-.63-.282-.63-.63V8.108c0-.348.282-.63.63-.63h2.386c.349 0 .63.282.63.63 0 .349-.281.631-.63.631H17.61v1.124h1.755zm-4.855 3.017c0 .349-.282.63-.63.63h-1.125v-1.26h1.124c.349 0 .631.282.631.63zm-1.124-2.522h1.124c.349 0 .631.282.631.631 0 .348-.282.63-.631.63h-1.124V10.358zm-2.502 2.522c0 .349-.282.63-.63.63H9.013v-1.26h1.241c.349 0 .631.282.631.63zm-1.241-2.522h1.241c.349 0 .631.282.631.631 0 .348-.282.63-.631.63H9.013V10.358zm-2.503 2.522c0 .349-.281.63-.63.63H5.629v-1.26h.754c.349 0 .63.282.63.63zm-.63-2.522h.754c.349 0 .63.282.63.631 0 .348-.281.63-.63.63H5.629V10.358zm-2.116 1.26h1.755c.349 0 .63.282.63.63 0 .349-.281.631-.63.631H2.5c-.348 0-.63-.282-.63-.631V8.108c0-.348.282-.63.63-.63h2.386c.349 0 .63.282.63.63 0 .349-.281.631-.63.631H3.13v1.124h1.755c.349 0 .63.283.63.631 0 .349-.281.63-.63.63H3.13v1.125zm5.655-1.26H6.38v1.26h1.755c.349 0 .63-.281.63-.63 0-.348-.281-.63-.63-.63zm1.241 0h-1.241v1.26h1.241c.349 0 .631-.281.631-.63 0-.348-.282-.63-.631-.63zm2.502 0h-1.124v1.26h1.124c.349 0 .631-.281.631-.63 0-.348-.282-.63-.631-.63z"
                />
              </svg>
              LINEでログイン
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
