import React from 'react';

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
  const pages = Array.from({ length: totalPages }, (_, i) => i + 1);
  const maxVisiblePages = 5;
  let visiblePages = pages;

  if (totalPages > maxVisiblePages) {
    const start = Math.max(1, Math.min(currentPage - 2, totalPages - maxVisiblePages + 1));
    const end = Math.min(start + maxVisiblePages - 1, totalPages);
    visiblePages = pages.slice(start - 1, end);
  }

  return (
    <div className={`flex items-center justify-center space-x-2 ${className}`}>
      <button
        onClick={() => onPageChange(currentPage - 1)}
        disabled={currentPage === 1}
        className="px-4 py-2 rounded-lg text-sm font-medium text-gray-600 bg-white border border-gray-100 hover:bg-gray-50 hover:border-gray-200 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
      >
        <svg
          className="w-4 h-4 inline-block mr-1"
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
        前へ
      </button>

      {visiblePages[0] > 1 && (
        <>
          <button
            onClick={() => onPageChange(1)}
            className="px-4 py-2 rounded-lg text-sm font-medium text-gray-600 bg-white border border-gray-100 hover:bg-gray-50 hover:border-gray-200 transition-all duration-200"
          >
            1
          </button>
          {visiblePages[0] > 2 && (
            <span className="px-2 text-gray-300">...</span>
          )}
        </>
      )}

      {visiblePages.map((page) => (
        <button
          key={page}
          onClick={() => onPageChange(page)}
          className={`px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 ${currentPage === page
            ? 'bg-blue-500 text-white border border-blue-500 shadow-sm'
            : 'text-gray-600 bg-white border border-gray-100 hover:bg-gray-50 hover:border-gray-200'
            }`}
        >
          {page}
        </button>
      ))}

      {visiblePages[visiblePages.length - 1] < totalPages && (
        <>
          {visiblePages[visiblePages.length - 1] < totalPages - 1 && (
            <span className="px-2 text-gray-300">...</span>
          )}
          <button
            onClick={() => onPageChange(totalPages)}
            className="px-4 py-2 rounded-lg text-sm font-medium text-gray-600 bg-white border border-gray-100 hover:bg-gray-50 hover:border-gray-200 transition-all duration-200"
          >
            {totalPages}
          </button>
        </>
      )}

      <button
        onClick={() => onPageChange(currentPage + 1)}
        disabled={currentPage === totalPages}
        className="px-4 py-2 rounded-lg text-sm font-medium text-gray-600 bg-white border border-gray-100 hover:bg-gray-50 hover:border-gray-200 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
      >
        次へ
        <svg
          className="w-4 h-4 inline-block ml-1"
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
