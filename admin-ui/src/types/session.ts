export interface Owner {
  userID: string;
  displayName: string;
  displayID: string;
  iconURL: string;
  lastLoginAt: string;
  createdAt: string;
  updatedAt: string;
}

export interface Session {
  talkSessionID: string;
  theme: string;
  description: string;
  owner: Owner;
  scheduledEndTime: string;
  thumbnailURL: string;
  hidden: boolean;
  updatedAt: string;
  createdAt: string;
  opinionCount: number;
  opinionUserCount: number;
  voteCount: number;
  voteUserCount: number;
}

export interface TalkSessionListResponse {
  totalCount: number;
  talkSessionStats: Session[];
}

export interface ReportResponse {
  code?: number;
  report?: string;
}

export interface CurrentUser {
  UserID: string;
  displayName: string;
  displayID: string;
  iconURL: string;
}

export interface UserStats {
  Date: string;
  UserCount: number;
  UniqueActionUserCount: number;
}

export interface UserStatsResponse {
  UserStats: UserStats[];
}
