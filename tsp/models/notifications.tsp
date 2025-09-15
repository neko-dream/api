import "@typespec/http";

using Http;
using OpenAPI;

namespace kotohiro {
  @doc("デバイス情報")
  model Device {
    device_id: string;
    user_id: string;
    platform: "ios" | "android" | "web";
    device_name?: string;
    enabled: boolean;
    last_active_at?: string;
    created_at: string;
    updated_at: string;
  }

  @doc("通知設定")
  model NotificationPreferences {
    push_notification_enabled: boolean;
  }
}
