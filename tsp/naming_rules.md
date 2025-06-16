# TypeSpec 命名規則ガイドライン

## 概要
このドキュメントは、Kotohiro APIプロジェクトにおけるTypeSpecの命名規則を定義します。
一貫性のあるAPI定義を維持するため、以下のルールに従ってください。

## モデル定義の命名規則

### モデル

- **PascalCase（大文字始まり）を使用**
- 単数形を使用
- ビジネスドメインの用語を使用
- 省略形は避ける

```tsp
// Good
model User {
  userId: string;
  displayName: string;
}

model Organization {
  organizationId: string;
  name: string;
}

// Bad
model user {}           // 小文字始まりは使用しない
model Users {}          // 複数形は使用しない
model Org {}           // 省略形は避ける
```

### プロパティ名

- **camelCase（小文字始まり）を使用**
- 説明的な名前を使用
- ID系は `{entity}ID` の形式で統一

```tsp
// Good
model TalkSession {
  talkSessionID: string;
  organizationID: string;
  createdAt: utcDateTime;
  isActive: boolean;
}

// Bad
model TalkSession {
  TalkSessionId: string;    // PascalCaseは使用しない
  org_id: string;          // snake_caseは使用しない
  created: utcDateTime;    // 曖昧な名前は避ける
}
```

### Enum定義

- **PascalCase（大文字始まり）を使用**
- Enum値もPascalCaseを使用
- 値の意味を明確に表現

```tsp
// Good
enum OrganizationUserRole {
  SuperAdmin: "SuperAdmin",
  Owner: "Owner",
  Admin: "Admin",
  Member: "Member",
}

enum VoteType {
  Agree: "AGREE",
  Disagree: "DISAGREE",
  Pass: "PASS",
}

// Bad
enum organizationUserRole {}  // 小文字始まりは使用しない
enum VOTE_TYPE {}            // 全て大文字は使用しない
```

## リクエスト/レスポンスの命名規則

### リクエストモデル

- 操作名 + `Request` または `Req`
- 操作の内容を明確に表現

```tsp
// Good
model CreateOrganizationRequest {
  name: string;
  code: string;
  orgType: OrganizationType;
}

model UpdateUserProfileRequest {
  displayName?: string;
  iconUrl?: string;
}

// Bad
model OrganizationRequest {}     // 操作が不明
model CreateOrgReq {}            // 省略形は避ける
```

### レスポンスモデル

- 操作名 + `Response` または単純な名前
- データの構造を明確に表現

```tsp
// Good
model OrganizationListResponse {
  organizations: Organization[];
  totalCount: int32;
}

model UserProfile {
  userId: string;
  displayName: string;
  iconUrl?: string;
}

// Bad
model OrgListRes {}              // 省略形は避ける
model Response {}                // 汎用的すぎる名前は避ける
```

## 名前空間（Namespace）の命名規則

- **PascalCase（大文字始まり）を使用**
- ドメインやリソースごとに整理
- 階層構造を適切に使用

```tsp
// Good
namespace Auth {
  model LoginRequest {}
  model TokenResponse {}
}

namespace Organization {
  model Organization {}
  model OrganizationUser {}
  model OrganizationAlias {}
}

// Bad
namespace auth {}                // 小文字始まりは使用しない
namespace OrganizationModels {}  // 冗長な名前は避ける
```

## 操作（Operation）の命名規則

### 操作ID

- **camelCase（小文字始まり）を使用**
- 動詞 + リソース名の形式
- RESTfulな命名を意識

```tsp
// Good
@route("/organizations")
@post
op createOrganization(...): Organization;

@route("/users/{userId}")
@get
op getUserById(...): User;

@route("/talk-sessions")
@get
op listTalkSessions(...): TalkSession[];

// Bad
@operationId("CreateOrganization")  // PascalCaseは使用しない
@operationId("get_user")            // snake_caseは使用しない
@operationId("organizations")       // 動詞がない
```

### HTTPメソッドとの対応

#### 基本的な対応

- GET: ドメインに応じた取得操作（`find`, `browse`, `search`など）
- POST: ビジネスアクション（`start`, `submit`, `register`, `invite`など）
- PUT/PATCH: 状態変更や更新（`edit`, `change`, `activate`, `complete`など）
- DELETE: 論理削除や無効化（`deactivate`, `cancel`, `revoke`など）

```tsp
// Good - ビジネスドメインを反映した命名
@post op startTalkSession(...): TalkSession;
@post op submitOpinion(...): Opinion;
@post op inviteOrganizationMember(...): void;
@patch op changeTalkSessionStatus(...): TalkSession;
@delete op deactivateOrganizationAlias(...): void;

// Avoid - 単純なCRUD操作
@post op createTalkSession(...): TalkSession;
@patch op updateTalkSession(...): TalkSession;
@delete op deleteTalkSession(...): void;
```

#### ドメイン別の例

**認証・ユーザー管理**
```tsp
@post op registerUser(...): User;
@post op authenticateUser(...): AuthToken;
@post op changePassword(...): void;
@delete op revokeSession(...): void;
```

## パスパラメータの命名規則

- **camelCase（小文字始まり）を使用**
- `{resource}ID` の形式で統一（IDは大文字）
- 説明的な名前を使用

```tsp
// Good
@route("/organizations/{organizationID}/users/{userID}")

@route("/talk-sessions/{talkSessionID}/opinions/{opinionID}")

// Bad
@route("/organizations/{id}")           // 何のIDか不明
@route("/organizations/{org_id}")       // snake_caseは使用しない
@route("/organizations/{orgId}")        // IDは大文字で統一
```

## クエリパラメータの命名規則

- **camelCase（小文字始まり）を使用**
- 一般的な規約に従う（page, limit, sort, filter等）

```tsp
// Good
model ListOrganizationsParams {
  @query page?: int32 = 1;
  @query limit?: int32 = 20;
  @query sortBy?: string;
  @query filterByType?: OrganizationType;
}

// Bad
model ListOrganizationsParams {
  @query Page?: int32;          // PascalCaseは使用しない
  @query per_page?: int32;      // snake_caseは使用しない
  @query s?: string;            // 省略形は避ける
}
```

## ファイル構成の規則

### ファイル名
- **kebab-case（ハイフン区切り）を使用**
- 内容を表す説明的な名前
- 適切にディレクトリで整理

```
tsp/
├── main.tsp
├── config/
│   └── service.tsp
├── models/
│   ├── common.tsp
│   ├── user.tsp
│   ├── organization.tsp
│   └── talk-session.tsp    // kebab-caseを使用
└── routes/
    ├── auth.tsp
    ├── user.tsp
    └── organization.tsp
```

### インポート構成
- 関連するモデルは適切にグループ化
- 循環参照を避ける
- 名前空間を活用して整理

```tsp
// Good - models/organization.tsp
import "./common.tsp";

namespace Organization {
  model Organization extends BaseModel {
    organizationID: string;
    name: string;
  }

  model OrganizationUser extends BaseModel {
    organizationUserID: string;
    role: OrganizationUserRole;
  }
}
```

## 一般的な規約

### 日時フィールド

- `createdAt`, `updatedAt` を使用
- `utcDateTime` 型を使用

```tsp
// Good
model BaseModel {
  createdAt: utcDateTime;
  updatedAt: utcDateTime;
}

// Bad
model BaseModel {
  created: string;        // 型が不適切
  modified_at: string;    // 命名規則が異なる
}
```

### Boolean フィールド
- `is`, `has`, `can` などのプレフィックスを使用

```tsp
// Good
model User {
  isActive: boolean;
  hasVerifiedEmail: boolean;
  canCreateOrganization: boolean;
}

// Bad
model User {
  active: boolean;        // プレフィックスがない
  emailVerified: boolean; // 一貫性がない
}
```

### 配列フィールド

- 複数形を使用
- 要素の型を明確に

```tsp
// Good
model Organization {
  users: OrganizationUser[];
  aliases: OrganizationAlias[];
}

// Bad
model Organization {
  user: OrganizationUser[];      // 単数形は避ける
  aliasList: OrganizationAlias[]; // Listサフィックスは不要
}
```
