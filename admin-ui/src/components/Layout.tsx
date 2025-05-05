import React, { useState } from 'react';
import { Outlet } from 'react-router-dom';
import { Sidebar } from './Sidebar';
import { useUser } from '../contexts/UserContext';

interface LayoutProps {
  children?: React.ReactNode;
}

export const Layout: React.FC<LayoutProps> = ({ children }) => {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  const { currentUser } = useUser();

  if (!currentUser) {
    return null;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* モバイル用ヘッダー */}
      <div className="lg:hidden fixed top-0 left-0 right-0 z-50 bg-white shadow-sm">
        <div className="flex items-center justify-between px-4 py-3">
          <h1 className="text-xl font-bold text-gray-800">ことひろ</h1>
          <button
            onClick={() => setIsSidebarOpen(!isSidebarOpen)}
            className="p-2 rounded-md text-gray-500 hover:text-gray-600 focus:outline-none"
          >
            <svg
              className="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M4 6h16M4 12h16M4 18h16"
              />
            </svg>
          </button>
        </div>
      </div>

      <div className="flex">
        {/* デスクトップ用サイドバー */}
        <div className="hidden lg:block lg:w-72 lg:fixed lg:left-0 lg:top-0 lg:h-screen lg:bg-white lg:shadow-lg lg:z-40">
          <Sidebar isOpen={true} onClose={() => { }} />
        </div>

        {/* モバイル用サイドバー */}
        <div className={`lg:hidden fixed top-0 left-0 right-0 bg-white shadow-lg z-40 transform transition-transform duration-300 ease-in-out ${isSidebarOpen ? 'translate-y-0' : '-translate-y-full'}`}>
          <div className="max-h-[60vh] overflow-y-auto">
            <Sidebar isOpen={isSidebarOpen} onClose={() => setIsSidebarOpen(false)} />
          </div>
        </div>

        {/* メインコンテンツ */}
        <div className="flex-1 min-h-screen lg:ml-72">
          <div className="pt-16 lg:pt-0">
            <main className="p-6">
              {children || <Outlet />}
            </main>
          </div>
        </div>
      </div>
    </div>
  );
};
