import React, { useEffect, useState } from 'react';

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  className?: string;
}

export const Pagination: React.FC<PaginationProps> = ({
  currentPage,
  totalPages,
  onPageChange,
  className = '',
}) => {
  // 最大表示数を常に4に固定
  const maxVisiblePages = 4;

  const getVisiblePages = () => {
    if (totalPages <= maxVisiblePages) {
      return Array.from({ length: totalPages }, (_, i) => i + 1);
    }
    // 1ページ目付近
    if (currentPage <= 2) {
      return [1, 2, 3, 4];
    }
    // 最終ページ付近
    if (currentPage >= totalPages - 1) {
      return [totalPages - 3, totalPages - 2, totalPages - 1, totalPages];
    }
    // 中間
    return [1, currentPage - 1, currentPage, totalPages];
  };

  const visiblePages = getVisiblePages();

  return (
    <div className={`flex items-center justify-center space-x-1 ${className}`}>
      <button
        onClick={() => onPageChange(currentPage - 1)}
        disabled={currentPage === 1}
        className="h-10 aspect-square px-0 py-0 rounded-xl border-2 text-base font-medium flex items-center justify-center leading-none bg-white text-gray-600 border-gray-200 hover:bg-gray-100 hover:border-gray-300 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
        aria-label="前のページ"
      >
        <svg
          className="w-5 h-5"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M15 19l-7-7 7-7"
          />
        </svg>
      </button>

      {visiblePages[0] > 1 && (
        <>
          <button
            onClick={() => onPageChange(1)}
            className="h-10 aspect-square px-0 py-0 rounded-xl border-2 text-base font-medium flex items-center justify-center leading-none bg-white text-gray-600 border-gray-100 hover:bg-gray-50 hover:border-gray-200 transition-all duration-200"
          >
            1
          </button>
          {visiblePages[0] > 2 && (
            <span className="px-1 text-gray-300">...</span>
          )}
        </>
      )}

      {visiblePages.map((page, idx) => {
        // ギャップ（...）の挿入
        if (idx > 0 && page - visiblePages[idx - 1] > 1) {
          return [
            <span key={`gap-${page}`} className="px-1 text-gray-300">...</span>,
            <button
              key={page}
              onClick={() => onPageChange(page)}
              className={`h-10 aspect-square px-0 py-0 rounded-xl border-2 text-base font-medium flex items-center justify-center leading-none transition-all duration-200 ${currentPage === page
                ? 'bg-blue-500 text-white border-blue-500 shadow-sm'
                : 'bg-white text-gray-600 border-gray-100 hover:bg-gray-50 hover:border-gray-200'}
                `}
            >
              {page}
            </button>
          ];
        }
        return (
          <button
            key={page}
            onClick={() => onPageChange(page)}
            className={`h-10 aspect-square px-0 py-0 rounded-xl border-2 text-base font-medium flex items-center justify-center leading-none transition-all duration-200 ${currentPage === page
              ? 'bg-blue-500 text-white border-blue-500 shadow-sm'
              : 'bg-white text-gray-600 border-gray-100 hover:bg-gray-50 hover:border-gray-200'}
              `}
          >
            {page}
          </button>
        );
      })}

      {visiblePages[visiblePages.length - 1] < totalPages && (
        <>
          {visiblePages[visiblePages.length - 1] < totalPages - 1 && (
            <span className="px-1 text-gray-300">...</span>
          )}
          <button
            onClick={() => onPageChange(totalPages)}
            className="h-10 aspect-square px-0 py-0 rounded-xl border-2 text-base font-medium flex items-center justify-center leading-none bg-white text-gray-600 border-gray-100 hover:bg-gray-50 hover:border-gray-200 transition-all duration-200"
          >
            {totalPages}
          </button>
        </>
      )}

      <button
        onClick={() => onPageChange(currentPage + 1)}
        disabled={currentPage === totalPages}
        className="h-10 aspect-square px-0 py-0 rounded-xl border-2 text-base font-medium flex items-center justify-center leading-none bg-white text-gray-600 border-gray-200 hover:bg-gray-100 hover:border-gray-300 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
        aria-label="次のページ"
      >
        <svg
          className="w-5 h-5"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M9 5l7 7-7 7"
          />
        </svg>
      </button>
    </div>
  );
};

