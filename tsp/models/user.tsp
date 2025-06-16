import "@typespec/http";

using Http;

namespace kotohiro {
  model User {
    @minLength(3)
    @maxLength(25)
    @pattern("^[A-Za-z0-9]$")
    displayID: string;

    displayName: string;
    iconURL?: string | null;
  }

  model UserDemographics {
    /**
     * 生年月日
     */
    @summary("20001010") dateOfBirth?: integer | null;

    /**
     * 性別
     */
    @summary("性別") gender?: string | null;

    /**
     * 市町村
     */
    @summary("市町村") city?: string | null;

    /**
     * 都道府県
     */
    @summary("都道府県") prefecture?: string | null;
  }

  model UserGroupPosition {
    posX: numeric;
    posY: numeric;
    displayID: string;
    displayName: string;
    iconURL?: string | null;
    groupName: string;
    groupID: integer;

    /**
     * 境界ポイントのインデックス
     */
    perimeterIndex?: integer;
  }
}