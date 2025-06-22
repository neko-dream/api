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
   */
  @tag("organization")
  @extension("x-ogen-operation-group", "Organization")
  @route("/organizations")
  @post
  @summary("組織作成（運営ユーザーのみ）")
  op establishOrganization(
    @multipartBody body: {
      name: HttpPart<string>;
      /**
       * 組織コード（ログイン時に使用）
       * 4文字以上127文字以下の英数字 _- のみ使用可能
       */
      code: HttpPart<string>;
      /**
       * 組織タイプ
       * - 1: 通常（基本これ）
       */
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
   * - 10: SuperAdmin
   * - 20: Owner
   * - 30: Admin
   * - 40: Member
   */
  @tag("organization")
  @extension("x-ogen-operation-group", "Organization")
  @route("/organizations/invite")
  @post
  @summary("組織ユーザー招待（運営ユーザーのみ）")
  op inviteOrganization(
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
  @route("/organizations/invite_user")
  @post
  @summary("組織にユーザーを追加")
  op inviteOrganizationForUser(
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
  @route("/organization/{code}/validate")
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
    aliases: OrganizationAlias[];
  }> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("organization")
  @extension("x-ogen-operation-group", "Organization")
  @route("/organizations/aliases")
  @post
  @summary("組織エイリアス作成")
  op createOrganizationAlias(
    @multipartBody body: {
      aliasName: HttpPart<string>;
    },
  ): Body<OrganizationAlias> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("organization")
  @extension("x-ogen-operation-group", "Organization")
  @route("/organizations/aliases/{aliasID}")
  @delete
  @summary("組織エイリアス削除")
  op deleteOrganizationAlias(
    @path aliasID: string,
  ): Body<{}> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };
}
