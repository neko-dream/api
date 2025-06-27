export interface TokenInfo {
  aud: string;
  exp: string;
  iat: string;
  iss: string;
  sub: string;
  jti: string;
  displayID?: string;
  displayName?: string;
  isRegistered: boolean;
  isEmailVerified: boolean;
  orgType?: number;
  requiredPasswordChange: boolean;
  organizationRole?: string;
  organizationCode?: string;
  organizationID?: string;
}

export interface Organization {
  ID: string;
  name: string;
  code: string;
  type: number;
  roleName: string;
  role: number;
}

export interface OrganizationsResponse {
  organizations: Organization[];
}

export const authService = {
  async getTokenInfo(): Promise<TokenInfo | null> {
    try {
      const response = await fetch('/auth/token/info', {
        method: 'GET',
        credentials: 'include',
      });

      if (!response.ok) {
        // If 403, redirect to login page (only if not already on login page)
        if (response.status === 403 && !window.location.pathname.includes('/login')) {
          window.location.href = '/admin/login';
        }
        return null;
      }

      const data = await response.json();
      return data;
    } catch (error) {
      console.error('Failed to fetch token info:', error);
      return null;
    }
  },

  async logout(): Promise<void> {
    try {
      await fetch('/auth/revoke', {
        method: 'POST',
        credentials: 'include',
      });
    } catch (error) {
      console.error('Failed to logout:', error);
    }
  },

  async getOrganizations(): Promise<Organization[]> {
    try {
      const response = await fetch('/organizations', {
        method: 'GET',
        credentials: 'include',
      });

      if (!response.ok) {
        return [];
      }

      const data: OrganizationsResponse = await response.json();
      return data.organizations || [];
    } catch (error) {
      console.error('Failed to fetch organizations:', error);
      return [];
    }
  },
};
