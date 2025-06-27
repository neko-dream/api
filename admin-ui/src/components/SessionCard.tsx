import { useState } from 'react';
import { Session } from '../types/session';
import { SessionReport } from './SessionReport';
import { useUser } from '../contexts/UserContext';
import { Card } from './Card';

interface SessionCardProps {
  session: Session;
  onRefetch: () => void;
}

export const SessionCard = ({ session, onRefetch }: SessionCardProps) => {
  const { currentUser } = useUser();

  if (!currentUser) {
    return null;
  }

  return (
    <Card title={session.theme} headerRight={
      <span
        className={`inline-flex items-center px-4 mx-4 py-1 rounded-full text-xs font-medium
          ${session.hidden
            ? 'bg-red-50/50 text-red-600 ring-1 ring-red-100'
            : 'bg-green-50/50 text-green-600 ring-1 ring-green-100'
          }`}
      >
        <svg
          className="w-3 h-3 mr-1.5"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        >
          {session.hidden ? (
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"
            />
          ) : (
            <>
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
              />
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
              />
            </>
          )}
        </svg>
        レポート: {session.hidden ? '非表示' : '表示'}
      </span>
    }>
      <div className="p-5 bg-gray-50">
        <div className="grid grid-cols-2 gap-3 text-sm text-gray-600">
          <div className="flex items-center">
            <i className="fas fa-fingerprint w-5 text-gray-400"></i>
            <span className="ml-2">セッションID: <span className="font-mono text-xs bg-gray-200 px-2 py-0.5 rounded">{session.talkSessionID}</span></span>
          </div>
          <div className="flex items-center">
            <i className="fas fa-user w-5 text-gray-400"></i>
            <span className="ml-2">作成者: <span className="font-semibold">{session.owner.displayName}</span></span>
          </div>
          <div className="flex items-center">
            <i className="far fa-calendar-plus w-5 text-gray-400"></i>
            <span className="ml-2">開始日時: {new Date(session.createdAt).toLocaleString()}</span>
          </div>
          <div className="flex items-center">
            <i className="far fa-calendar-check w-5 text-gray-400"></i>
            <span className="ml-2">終了日時: {new Date(session.scheduledEndTime).toLocaleString()}</span>
          </div>
          <div className="flex items-center">
            <i className="far fa-comment-dots w-5 text-gray-400"></i>
            <span className="ml-2">意見数: <span className="font-semibold">{session.opinionCount}</span></span>
          </div>
          <div>
            <i className="fas fa-thumbs-up w-5 text-gray-400"></i>
            <span className="ml-2">投票数: <span className="font-semibold">{session.voteCount}</span></span>
          </div>
          <div>
            <i className="fas fa-user-friends w-5 text-gray-400"></i>
            <span className="ml-2">投票ユーザー数: <span className="font-semibold">{session.voteUserCount}</span></span>
          </div>
        </div>
      </div>

      <SessionReport session={session} onRefetch={onRefetch} />
    </Card>
  );
};
