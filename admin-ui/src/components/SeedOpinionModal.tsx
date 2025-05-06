import { useState } from 'react';
import { useNotification } from '../contexts/NotificationContext';

interface SeedOpinionModalProps {
  sessionId: string;
  onClose: () => void;
}

export const SeedOpinionModal = ({ sessionId, onClose }: SeedOpinionModalProps) => {
  const [opinionContent, setOpinionContent] = useState('');
  const [referenceURL, setReferenceURL] = useState('');
  const { showNotification } = useNotification();

  const submitSeedOpinion = async () => {
    try {
      if (!opinionContent || opinionContent.length < 5 || opinionContent.length > 140) {
        showNotification('意見は5〜140文字で入力してください', 'error');
        return;
      }

      const formData = new FormData();
      formData.append('talkSessionID', sessionId);
      formData.append('opinionContent', opinionContent);
      formData.append('isSeed', 'true');

      if (referenceURL) formData.append('referenceURL', referenceURL);

      const response = await fetch('/opinions', {
        method: 'POST',
        body: formData,
      });

      if (response.ok) {
        showNotification('シード意見を投稿しました', 'success');
        onClose();
        setOpinionContent('');
        setReferenceURL('');
      } else {
        const errorData = await response.json();
        showNotification(`投稿に失敗しました: ${errorData.message || '不明なエラー'}`, 'error');
      }
    } catch (error) {
      console.error('Error:', error);
      showNotification('エラーが発生しました', 'error');
    }
  };

  return (
    <div className="fixed inset-0 overflow-y-auto z-50 flex items-center justify-center">
      <div className="fixed inset-0 bg-black opacity-50 transition-opacity" onClick={onClose}></div>

      <div className="relative bg-white rounded-lg max-w-lg w-full mx-4 shadow-xl">
        <div className="flex justify-between items-center px-6 py-4 border-b">
          <h3 className="text-lg font-medium text-gray-900">シード意見の投稿</h3>
          <button type="button" onClick={onClose} className="text-gray-400 hover:text-gray-600">
            <i className="fas fa-times"></i>
          </button>
        </div>

        <div className="p-6">
          <div className="mb-4">
            <label htmlFor="opinionContent" className="block text-sm font-medium text-gray-700 mb-1">意見内容 <span className="text-red-600">*</span></label>
            <textarea
              id="opinionContent"
              value={opinionContent}
              onChange={(e) => setOpinionContent(e.target.value)}
              rows={4}
              required
              className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            ></textarea>
            <p className="mt-1 text-sm text-gray-500">5〜140文字で入力してください</p>
          </div>

          <div className="mb-4">
            <label htmlFor="referenceURL" className="block text-sm font-medium text-gray-700 mb-1">参考URL (任意)</label>
            <input
              type="url"
              id="referenceURL"
              value={referenceURL}
              onChange={(e) => setReferenceURL(e.target.value)}
              className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>

          <div className="mt-6 flex justify-end space-x-3">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              キャンセル
            </button>
            <button
              type="button"
              onClick={submitSeedOpinion}
              className="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              投稿する
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};
