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
    <Card className="w-full max-w-md" title="ログイン">
      <div className="px-6">
        <p className="text-gray-600 mb-6">Googleアカウントでログインしてください</p>

        <div className="space-y-4">
          <div>
            <label htmlFor="organization-code" className="block text-sm font-medium text-gray-700 mb-1">
              組織（オプション）
            </label>
            {organizations.length > 0 ? (
              <select
                id="organization-code"
                value={organizationCode}
                onChange={(e) => setOrganizationCode(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
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
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                disabled={isLoadingOrgs}
              />
            )}
          </div>

          <button
            onClick={handleGoogleLogin}
            disabled={isLoadingOrgs}
            className="w-full py-2 px-4 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-md transition-colors duration-200 disabled:bg-gray-400 disabled:cursor-not-allowed"
          >
            {isLoadingOrgs ? '読み込み中...' : 'Googleでログイン'}
          </button>
        </div>
      </div>
    </Card>
  );
};
