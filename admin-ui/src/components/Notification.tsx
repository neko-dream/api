import { useEffect, useState } from 'react';

interface NotificationProps {
  message: string;
  type: 'success' | 'error';
  onClose: () => void;
}

export const Notification = ({ message, type, onClose }: NotificationProps) => {
  const [isVisible, setIsVisible] = useState(true);

  useEffect(() => {
    const timer = setTimeout(() => {
      setIsVisible(false);
      setTimeout(onClose, 300);
    }, 3000);

    return () => clearTimeout(timer);
  }, [onClose]);

  return (
    <div
      className={`fixed bottom-4 right-4 px-6 py-3 rounded-md text-white ${type === 'success' ? 'bg-green-600' : 'bg-red-600'
        } shadow-lg transition-opacity duration-300 ${isVisible ? 'opacity-100' : 'opacity-0'}`}
    >
      <div className="flex items-center">
        <i className={`fas ${type === 'success' ? 'fa-check-circle' : 'fa-exclamation-circle'} mr-2`}></i>
        <span>{message}</span>
      </div>
    </div>
  );
};
