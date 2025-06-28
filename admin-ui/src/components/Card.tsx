import React from 'react';

interface CardProps {
  title?: string;
  headerLeft?: React.ReactNode;
  headerRight?: React.ReactNode;
  children: React.ReactNode;
  className?: string;
}

export const Card: React.FC<CardProps> = ({
  title,
  headerLeft,
  headerRight,
  children,
  className = ''
}) => {
  return (
    <div className={`bg-white rounded-xl shadow-sm border border-gray-100 hover:shadow-lg transition-all duration-300 ${className} mb-6`}>
      {(title || headerLeft || headerRight) && (
        <div className="flex justify-between items-center px-6 py-4 border-b border-gray-100">
          <div className="flex items-center">
            {headerLeft}
            {title && !headerLeft && <h2 className="text-xl font-semibold text-gray-800">{title}</h2>}
          </div>
          {headerRight}
        </div>
      )}
      <div className="p-6">
        {children}
      </div>
    </div>
  );
};
