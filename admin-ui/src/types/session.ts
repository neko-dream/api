export interface Owner {
  DisplayName: string;
  DisplayID: string;
}

export interface Session {
  TalkSessionID: string;
  Theme: string;
  Hidden: boolean;
  Owner: Owner;
  CreatedAt: string;
  ScheduledEndTime: string;
  OpinionCount: number;
  VoteCount: number;
  VoteUserCount: number;
}

export interface TalkSessionListResponse {
  TotalCount: number;
  TalkSessionStats: Session[];
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
