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
  const [copied, setCopied] = useState(false);

  if (!currentUser) {
    return null;
  }

  return (
    <Card title={
      <div className="flex items-center gap-3">
        <div className="w-2 h-2 rounded-full bg-gradient-to-r from-blue-500 to-indigo-500"></div>
        <span className="text-xl font-semibold text-gray-800">{session.theme}</span>
      </div>
    } headerRight={
      <span
        className={`inline-flex items-center px-4 py-1.5 rounded-full text-xs font-medium transition-all duration-200
          ${session.hidden
            ? 'bg-red-100 text-red-700 shadow-sm'
            : 'bg-emerald-100 text-emerald-700 shadow-sm'
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
      <div className="-m-6 p-6 bg-gradient-to-r from-gray-50 to-gray-50/50">
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-3">
            <div className="flex items-start gap-3">
              <div className="p-2 bg-blue-100 rounded-lg">
                <i className="fas fa-fingerprint text-blue-600"></i>
              </div>
              <div className="flex-1">
                <p className="text-xs text-gray-500">セッションID</p>
                <div className="flex items-center gap-2 mt-1">
                  <p className="font-mono text-xs bg-white px-2 py-1 rounded shadow-sm inline-block break-all">{session.talkSessionID}</p>
                  <button
                    onClick={() => {
                      navigator.clipboard.writeText(session.talkSessionID);
                      setCopied(true);
                      setTimeout(() => setCopied(false), 2000);
                    }}
                    className="p-1.5 bg-white rounded shadow-sm hover:bg-gray-50 transition-colors relative"
                    title="IDをコピー"
                  >
                    <i className={`${copied ? 'fas fa-check' : 'far fa-copy'} text-gray-600 text-xs`}></i>
                    {copied && (
                      <span className="absolute -top-8 left-1/2 transform -translate-x-1/2 text-xs bg-gray-800 text-white px-2 py-1 rounded whitespace-nowrap">
                        コピーしました
                      </span>
                    )}
                  </button>
                </div>
              </div>
            </div>
            <div className="flex items-start gap-3">
              <div className="p-2 bg-purple-100 rounded-lg">
                <i className="fas fa-user text-purple-600"></i>
              </div>
              <div>
                <p className="text-xs text-gray-500">作成者</p>
                <p className="font-semibold text-sm text-gray-800">{session.owner.displayName}</p>
              </div>
            </div>
          </div>
          <div className="space-y-3">
            <div className="flex items-start gap-3">
              <div className="p-2 bg-indigo-100 rounded-lg">
                <i className="far fa-calendar text-indigo-600"></i>
              </div>
              <div>
                <p className="text-xs text-gray-500">期間</p>
                <p className="text-sm text-gray-800">{new Date(session.createdAt).toLocaleDateString()} 〜</p>
                <p className="text-sm text-gray-800">{new Date(session.scheduledEndTime).toLocaleDateString()}</p>
              </div>
            </div>
          </div>
        </div>
        
        <div className="grid grid-cols-3 gap-4 mt-6">
          <div className="bg-white rounded-xl p-3 text-center shadow-sm">
            <i className="far fa-comment-dots text-lg text-blue-500 mb-1"></i>
            <p className="text-lg font-bold text-gray-800">{session.opinionCount}</p>
            <p className="text-xs text-gray-500">意見数</p>
          </div>
          <div className="bg-white rounded-xl p-3 text-center shadow-sm">
            <i className="fas fa-thumbs-up text-lg text-green-500 mb-1"></i>
            <p className="text-lg font-bold text-gray-800">{session.voteCount}</p>
            <p className="text-xs text-gray-500">投票数</p>
          </div>
          <div className="bg-white rounded-xl p-3 text-center shadow-sm">
            <i className="fas fa-users text-lg text-purple-500 mb-1"></i>
            <p className="text-lg font-bold text-gray-800">{session.voteUserCount}</p>
            <p className="text-xs text-gray-500">参加者数</p>
          </div>
        </div>
      </div>

      <SessionReport session={session} onRefetch={onRefetch} />
    </Card>
  );
};
