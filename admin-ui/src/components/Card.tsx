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
    <div className={`bg-white rounded-lg shadow-md pt-4 pb-4 ${className} mb-4 transition-shadow duration-300`}>
      {(title || headerLeft || headerRight) && (
        <div className="flex justify-between items-center mb-4">
          <div className="flex items-center">
            {headerLeft}
            {title && !headerLeft && <h2 className="text-xl font-semibold mx-4">{title}</h2>}
          </div>
          {headerRight}
        </div>
      )}
      {children}
    </div>
  );
};
