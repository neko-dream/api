import "@typespec/http";
import "@typespec/openapi";
import "../config/service.tsp";
import "../models/notifications.tsp";

using Http;
using OpenAPI;

namespace kotohiro {
  @tag("notifications")
  @extension("x-ogen-operation-group", "Notifications")
  @route("/notifications/vapid-key")
  @get
  @summary("VAPID公開鍵取得")
  op getVapidKey(): Body<{
    vapid_key: string;
  }>;

  @tag("notifications")
  @extension("x-ogen-operation-group", "Notifications")
  @route("/notifications/devices")
  @post
  @summary("デバイス登録/更新")
  op registerDevice(
    /**
     * FCMトークンまたはAPNSトークン
     */
    @multipartBody body: {
      device_token: HttpPart<string>;
      platform: HttpPart<"ios" | "android" | "web">;
      device_name?: HttpPart<string>;
      app_version?: HttpPart<string>;
      os_version?: HttpPart<string>;
    }
  ): Body<Device> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 401;
    @body body: {};
  };

  @tag("notifications")
  @extension("x-ogen-operation-group", "Notifications")
  @route("/notifications/devices")
  @get
  @summary("デバイス一覧取得")
  op getDevices(): Body<{
    devices: Device[];
  }> | {
    @statusCode statusCode: 401;
    @body body: {};
  };

  @tag("notifications")
  @extension("x-ogen-operation-group", "Notifications")
  @route("/notifications/devices/exists")
  @get
  @summary("デバイストークンが登録されているか確認")
  op checkDeviceExists(
    /**
     * FCMトークンまたはAPNSトークン
     */
    @query device_token: string,
  ): {
    @statusCode statusCode: 200;
    @body body: {
      exists: boolean;
    };
  } | {
    @statusCode statusCode: 401;
    @body body: {};
  } | {
    @statusCode statusCode: 404;
    @body body: {};
  };


  @tag("notifications")
  @extension("x-ogen-operation-group", "Notifications")
  @route("/notifications/devices/{deviceId}")
  @delete
  @summary("デバイス削除")
  op deleteDevice(
    @path deviceId: string
  ): {
    @statusCode statusCode: 204;
  } | {
    @statusCode statusCode: 401;
    @body body: {};
  } | {
    @statusCode statusCode: 404;
    @body body: {};
  };

  @tag("notifications")
  @extension("x-ogen-operation-group", "Notifications")
  @route("/notifications/preferences")
  @get
  @summary("通知設定取得")
  op getNotificationPreferences(): Body<NotificationPreferences> | {
    @statusCode statusCode: 401;
    @body body: {};
  };

  @tag("notifications")
  @extension("x-ogen-operation-group", "Notifications")
  @route("/notifications/preferences")
  @put
  @summary("通知設定更新")
  op updateNotificationPreferences(
    @multipartBody body: {
      push_notification_enabled: HttpPart<boolean | null>;
    }
  ): Body<NotificationPreferences> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 401;
    @body body: {};
  };

  @tag("notifications")
  @extension("x-ogen-operation-group", "Notifications")
  @route("/notifications/test")
  @post
  @summary("テスト通知送信")
  op sendTestNotification(
    @multipartBody body: {
      title: HttpPart<string>;
      body: HttpPart<string>;
      device_id?: HttpPart<string>;
    }
  ): Body<{
    devices_count: int32;
    success_count: int32;
  }> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 401;
    @body body: {};
  };
}
