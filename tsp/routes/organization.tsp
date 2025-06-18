import "@typespec/http";
import "@typespec/openapi";
import "../models/organization.tsp";
import "../config/service.tsp";

using Http;
using OpenAPI;

namespace kotohiro {
  @tag("organization")
  @extension("x-ogen-operation-group", "Organization")
  @route("/organizations")
  @get
  @summary("所属組織一覧")
  op getOrganizations(): Body<{
    organizations: Organization[];
  }> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  /**
   * 組織を作成できる。
   * これを作れるユーザーはDBを直接叩いて作るしかない。
   *
   * OrgType
   * - 1: 通常
   * - 2: 自治体
   * - 3: 議員
   */
  @tag("organization")
  @extension("x-ogen-operation-group", "Organization")
  @route("/organizations")
  @post
  @summary("組織作成（運営ユーザーのみ）")
  op establishOrganization(
    @multipartBody body: {
      name: HttpPart<string>;
      code: HttpPart<string>;
      orgType: HttpPart<numeric>;
    },
  ): Body<{}> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  /**
   * Role
   * - 1: Member
   * - 2: Admin
   * - 3: Owner
   */
  @tag("organization")
  @extension("x-ogen-operation-group", "Organization")
  @route("/organizations/{organizationID}/invite")
  @post
  @summary("組織ユーザー招待（運営ユーザーのみ）")
  op inviteOrganization(
    @path organizationID: string,
    @multipartBody body: {
      email: HttpPart<string>;
      role: HttpPart<numeric>;
    },
  ): Body<{}> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("organization")
  @extension("x-ogen-operation-group", "Organization")
  @route("/organizations/{organizationID}/invite_user")
  @post
  @summary("組織にユーザーを追加")
  op inviteOrganizationForUser(
    @path organizationID: string,
    @multipartBody body: {
      displayID: HttpPart<string>;
      role: HttpPart<numeric>;
    },
  ): Body<{
    success: boolean;
  }> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("organization")
  @extension("x-ogen-operation-group", "Organization")
  @route("/auth/organization/{code}/validate")
  @get
  @summary("組織コード検証")
  @useAuth([])
  op validateOrganizationCode(@path code: string): Body<{
    valid: boolean;
    organization?: Organization;
  }> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("organization")
  @extension("x-ogen-operation-group", "Organization")
  @route("/organizations/{organizationID}/aliases")
  @get
  @summary("組織エイリアス一覧取得")
  op getOrganizationAliases(@path organizationID: string): Body<{
    aliases: {
      aliasID: string;
      aliasName: string;
      createdAt: string;
    }[];
  }> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("organization")
  @extension("x-ogen-operation-group", "Organization")
  @route("/organizations/{organizationID}/aliases")
  @post
  @summary("組織エイリアス作成")
  op createOrganizationAlias(
    @path organizationID: string,
    @multipartBody body: {
      aliasName: HttpPart<string>;
    },
  ): Body<{
    aliasID: string;
    aliasName: string;
    createdAt: string;
  }> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("organization")
  @extension("x-ogen-operation-group", "Organization")
  @route("/organizations/{organizationID}/aliases/{aliasID}")
  @delete
  @summary("組織エイリアス削除")
  op deleteOrganizationAlias(
    @path organizationID: string,
    @path aliasID: string,
  ): Body<{}> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };
}