import { useState } from 'react';
import { marked } from 'marked';
import DOMPurify from 'dompurify';
import { Session } from '../types/session';
import { useNotification } from '../contexts/NotificationContext';
import { useMutation } from '@tanstack/react-query';
import { Link } from '@tanstack/react-router';
import { useUser } from '../contexts/UserContext';
import { SeedOpinionModal } from './SeedOpinionModal';

interface SessionReportProps {
  session: Session;
  onRefetch: () => void;
}

interface ApiResponse {
  status?: string;
  code?: string;
  report?: string;
}

// API呼び出しの処理を分離
const api = {
  toggleHideReport: async (talkSessionID: string, hideStatus: boolean): Promise<ApiResponse> => {
    try {
      const formData = new URLSearchParams();
      formData.append('hidden', hideStatus.toString());
      const response = await fetch(`/v1/manage/talksessions/${talkSessionID}/analysis/report`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: formData.toString(),
      });
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      return response.json();
    } catch (error) {
      console.error('Error in toggleHideReport:', error);
      throw error;
    }
  },

  generateAnalysis: async (talkSessionID: string, type: string): Promise<ApiResponse> => {
    try {
      const formData = new URLSearchParams();
      formData.append('type', type);
      const response = await fetch(`/v1/manage/talksessions/${talkSessionID}/analysis/regenerate`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: formData.toString(),
      });
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      return response.json();
    } catch (error) {
      console.error('Error in generateAnalysis:', error);
      throw error;
    }
  },

  getReport: async (talkSessionID: string): Promise<ApiResponse> => {
    try {
      const response = await fetch(`/v1/manage/talksessions/${talkSessionID}/analysis/report`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      return response.json();
    } catch (error) {
      console.error('Error in getReport:', error);
      throw error;
    }
  },
};

export const SessionReport = ({ session, onRefetch }: SessionReportProps) => {
  const [expanded, setExpanded] = useState(false);
  const [reportContent, setReportContent] = useState<string | null>(null);
  const [isGeneratingAnalysis, setIsGeneratingAnalysis] = useState(false);
  const [isGeneratingReport, setIsGeneratingReport] = useState(false);
  const [showSeedModal, setShowSeedModal] = useState(false);
  const { showNotification } = useNotification();
  const { currentUser } = useUser();

  const toggleHideReport = useMutation<void, Error, boolean>({
    mutationFn: async (hideStatus) => {
      const result = await api.toggleHideReport(session.TalkSessionID, hideStatus);
      if (result.status === "success" || result.status === "ok") {
        onRefetch();
        showNotification(`レポートを${hideStatus ? '非表示' : '表示'}にしました`, 'success');
      } else {
        throw new Error('更新に失敗しました');
      }
    },
    onError: (error) => {
      console.error('Error in toggleHideReport:', error);
      showNotification(error.message || 'エラーが発生しました', 'error');
    },
  });

  const generateAnalysis = useMutation<void, Error, string>({
    mutationFn: async (type) => {
      const result = await api.generateAnalysis(session.TalkSessionID, type);
      if (result.status === "success" || result.status === "ok") {
        showNotification('分析を再実行しました', 'success');
      } else {
        throw new Error('分析の再実行に失敗しました');
      }
    },
    onMutate: (type) => {
      if (type === 'group') {
        setIsGeneratingAnalysis(true);
      } else if (type === 'report') {
        setIsGeneratingReport(true);
      }
    },
    onSettled: (_, __, type) => {
      if (type === 'group') {
        setIsGeneratingAnalysis(false);
      } else if (type === 'report') {
        setIsGeneratingReport(false);
      }
    },
    onError: (error) => {
      console.error('Error in generateAnalysis:', error);
      showNotification(error.message || 'エラーが発生しました', 'error');
    },
  });

  const toggleReport = async () => {
    if (!expanded) {
      setExpanded(true);
      if (!reportContent) {
        try {
          const result = await api.getReport(session.TalkSessionID);
          if (!result.code && result.report) {
            setReportContent(result.report);
          } else {
            setExpanded(false);
            showNotification('レポートの取得に失敗しました', 'error');
          }
        } catch (error) {
          console.error('Error fetching report:', error);
          setExpanded(false);
          showNotification('レポートの取得に失敗しました', 'error');
        }
      }
    } else {
      setExpanded(false);
    }
  };

  const getRenderedMarkdown = (markdownContent: string) => {
    try {
      const parsedMarkdown = marked.parse(markdownContent) as string;
      const sanitizedHtml = DOMPurify.sanitize(parsedMarkdown);
      return { __html: sanitizedHtml };
    } catch (error) {
      console.error('Error parsing markdown:', error);
      return { __html: 'マークダウンの解析に失敗しました' };
    }
  };

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center px-4 mt-4">
        <div className="flex items-center space-x-2">
          <button
            onClick={() => toggleHideReport.mutate(!session.Hidden)}
            disabled={toggleHideReport.isPending}
            className={`inline-flex items-center px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 ease-in-out
              ${session.Hidden
                ? 'bg-gradient-to-r from-red-50 to-rose-50 text-red-600 hover:text-red-700 border border-red-100 hover:border-red-200'
                : 'bg-gradient-to-r from-green-50 to-emerald-50 text-green-600 hover:text-green-700 border border-green-100 hover:border-green-200'
              }
              shadow-sm hover:shadow-md
              hover:scale-[1.02] active:scale-[0.98]
              disabled:opacity-50 disabled:cursor-not-allowed`}
          >
            {toggleHideReport.isPending ? (
              <>
                <svg className="animate-spin -ml-1 mr-2 h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                処理中...
              </>
            ) : (
              <>
                <svg
                  className="w-4 h-4 mr-2"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  {session.Hidden ? (
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
                レポート{session.Hidden ? '表示' : '非表示'}
              </>
            )}
          </button>

          {currentUser && currentUser.displayID === session.Owner.DisplayID && (
            <button
              onClick={() => setShowSeedModal(true)}
              className="inline-flex items-center px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 ease-in-out
                bg-gradient-to-r from-yellow-50 to-amber-50
                text-yellow-600 hover:text-yellow-700
                border border-yellow-100 hover:border-yellow-200
                shadow-sm hover:shadow-md
                hover:scale-[1.02] active:scale-[0.98]"
            >
              <svg
                className="w-4 h-4 mr-2"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"
                />
              </svg>
              シード意見投稿
            </button>
          )}
        </div>

        <div className="flex items-center space-x-2">
          <button
            onClick={() => generateAnalysis.mutate('group')}
            disabled={isGeneratingAnalysis}
            className="inline-flex items-center px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 ease-in-out
              bg-gradient-to-r from-blue-50 to-sky-50
              text-blue-600 hover:text-blue-700
              border border-blue-100 hover:border-blue-200
              shadow-sm hover:shadow-md
              hover:scale-[1.02] active:scale-[0.98]
              disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {isGeneratingAnalysis ? (
              <>
                <svg className="animate-spin -ml-1 mr-2 h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                分析中...
              </>
            ) : (
              <>
                <svg
                  className="w-4 h-4 mr-2"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M13 10V3L4 14h7v7l9-11h-7z"
                  />
                </svg>
                分析
              </>
            )}
          </button>

          <button
            onClick={() => generateAnalysis.mutate('report')}
            disabled={isGeneratingReport}
            className="inline-flex items-center px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 ease-in-out
              bg-gradient-to-r from-purple-50 to-violet-50
              text-purple-600 hover:text-purple-700
              border border-purple-100 hover:border-purple-200
              shadow-sm hover:shadow-md
              hover:scale-[1.02] active:scale-[0.98]
              disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {isGeneratingReport ? (
              <>
                <svg className="animate-spin -ml-1 mr-2 h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                生成中...
              </>
            ) : (
              <>
                <svg
                  className="w-4 h-4 mr-2"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                  />
                </svg>
                レポート
              </>
            )}
          </button>
        </div>
      </div>

      <div className="border-t border-gray-200">
        <button
          onClick={toggleReport}
          className={`w-full px-4 py-3 text-sm font-medium ${session.OpinionCount === 0 ? 'text-gray-400 cursor-not-allowed' : 'text-gray-700 hover:bg-gray-50'} transition-colors duration-200 flex items-center justify-center`}
          disabled={session.OpinionCount === 0}
          title={session.OpinionCount === 0 ? '意見が投稿されていないため、レポートをプレビューできません' : ''}
        >
          <i className="fas fa-file-alt mr-2"></i>
          <span>{expanded ? 'レポートを隠す' : 'レポートをプレビュー'}</span>
          <i className={`fas fa-chevron-down ml-2 transform transition-transform duration-300 ${expanded ? 'rotate-180' : ''}`}></i>
        </button>

        {expanded && (
          <div className="p-5 bg-gray-50 border-t border-gray-200 report-content">
            {!reportContent ? (
              <div className="flex justify-center py-4">
                <i className="fas fa-spinner fa-spin text-gray-400 text-2xl"></i>
              </div>
            ) : (
              <div className="markdown-content" dangerouslySetInnerHTML={getRenderedMarkdown(reportContent)} />
            )}
          </div>
        )}
      </div>

      {showSeedModal && (
        <SeedOpinionModal
          sessionId={session.TalkSessionID}
          onClose={() => setShowSeedModal(false)}
        />
      )}
    </div>
  );
};
