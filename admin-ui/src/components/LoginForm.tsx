import { useState, useEffect } from 'react';
import { Card } from './Card';
import { authService, Organization } from '@/services/auth';

export const LoginForm = () => {
  const [organizationCode, setOrganizationCode] = useState('');
  const [organizations, setOrganizations] = useState<Organization[]>([]);
  const [isLoadingOrgs, setIsLoadingOrgs] = useState(false);
  const redirectUrl = `${window.location.origin}/admin/`;

  useEffect(() => {
    const fetchOrganizations = async () => {
      setIsLoadingOrgs(true);
      try {
        const orgs = await authService.getOrganizations();
        setOrganizations(orgs);
        // If only one organization, auto-select it
        if (orgs.length === 1) {
          setOrganizationCode(orgs[0].code);
        }
      } catch (error) {
        console.error('Failed to fetch organizations:', error);
      } finally {
        setIsLoadingOrgs(false);
      }
    };

    fetchOrganizations();
  }, []);

  const handleGoogleLogin = () => {
    const baseUrl = '/auth/google/login';
    const params = new URLSearchParams({
      redirect_url: redirectUrl,
    });
    
    // Only add organization_code if it has a value
    if (organizationCode) {
      params.append('organization_code', organizationCode);
    }

    const loginUrl = `${baseUrl}?${params.toString()}`;
    window.location.href = loginUrl;
  };

  return (
    <div className="w-full max-w-md">
      <div className="text-center mb-8">
        <h1 className="text-4xl font-bold bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent mb-2">
          ことひろ
        </h1>
        <p className="text-gray-600">管理画面へようこそ</p>
      </div>
      
      <Card className="backdrop-blur-sm bg-white/95 shadow-2xl border-0">
        <div className="space-y-6">
          <div>
            <label htmlFor="organization-code" className="block text-sm font-medium text-gray-700 mb-1">
              組織（オプション）
            </label>
            {organizations.length > 0 ? (
              <select
                id="organization-code"
                value={organizationCode}
                onChange={(e) => setOrganizationCode(e.target.value)}
                className="w-full px-4 py-3 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-gray-50 hover:bg-white transition-colors duration-200"
                disabled={isLoadingOrgs}
              >
                <option value="">組織を選択（オプション）</option>
                {organizations.map((org) => (
                  <option key={org.ID} value={org.code}>
                    {org.name} ({org.roleName})
                  </option>
                ))}
              </select>
            ) : (
              <input
                id="organization-code"
                type="text"
                placeholder="組織コード（オプション）"
                value={organizationCode}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setOrganizationCode(e.target.value)}
                className="w-full px-4 py-3 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-gray-50 hover:bg-white transition-colors duration-200"
                disabled={isLoadingOrgs}
              />
            )}
          </div>

          <button
            onClick={handleGoogleLogin}
            disabled={isLoadingOrgs}
            className="w-full py-4 px-6 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 text-white font-semibold rounded-xl transition-all duration-200 disabled:from-gray-400 disabled:to-gray-400 disabled:cursor-not-allowed shadow-lg hover:shadow-xl transform hover:-translate-y-0.5 flex items-center justify-center gap-3"
          >
            <svg className="w-5 h-5" viewBox="0 0 24 24" fill="currentColor">
              <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
              <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
              <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
              <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
            </svg>
            {isLoadingOrgs ? '読み込み中...' : 'Googleでログイン'}
          </button>
        </div>
      </Card>
    </div>
  );
};
