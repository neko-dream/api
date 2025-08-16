# Kotohiro Web Push通知実装ガイド（フロントエンド向け）

このドキュメントでは、KotohiroのWeb Push通知機能をフロントエンドアプリケーションに実装する方法を説明します。

## 概要

Kotohiroは、Firebase Cloud Messaging（FCM）を使用したWeb Push通知をサポートしています。ユーザーは以下の機能を利用できます：

- デバイスの登録・管理
- 通知設定の管理
- トークセッションの開始/終了などのイベント通知の受信
- テスト通知の送信

## 技術スタック

- **Firebase Cloud Messaging (FCM)**: プッシュ通知配信サービス
- **Service Worker**: バックグラウンド通知の処理
- **REST API**: デバイス登録・管理用エンドポイント

## FCMの仕組み

### 1. FCMの基本的な流れ

```mermaid
sequenceDiagram
    participant User as ユーザー
    participant Browser as ブラウザ
    participant FCM as FCM Server
    participant Backend as Kotohiro Backend
    participant AWS as AWS Pinpoint

    User->>Browser: 通知権限を許可
    Browser->>FCM: トークンをリクエスト
    FCM-->>Browser: FCMトークンを返却
    Browser->>Backend: デバイス登録（トークン送信）
    Backend-->>Browser: デバイス情報を返却

    Note over Backend,AWS: 通知送信時
    Backend->>AWS: 通知送信リクエスト
    AWS->>FCM: FCMトークンに通知送信
    FCM->>Browser: プッシュ通知配信
    Browser->>User: 通知を表示
```

### 2. FCMトークンとは

FCMトークンは、特定のデバイス/ブラウザインスタンスを識別する一意の文字列です。このトークンを使って、サーバーから特定のデバイスに通知を送信できます。

- トークンは自動的に更新される可能性があります
- ブラウザのデータをクリアするとトークンが変わります
- 同じユーザーでも複数のデバイス/ブラウザで異なるトークンを持ちます

## 実装手順

### 1. Firebase設定

プロジェクトでFirebaseを初期化します：

```javascript
// Firebase設定（開発環境用）
const firebaseConfig = {
    apiKey: "AIzaSyBEQ-gqPCUMyyBZ-YP10t_0X0D01BGPCgQ",
    authDomain: "kotohiro-dev-c57a8.firebaseapp.com",
    projectId: "kotohiro-dev-c57a8",
    storageBucket: "kotohiro-dev-c57a8.firebasestorage.app",
    messagingSenderId: "1023394330741",
    appId: "1:1023394330741:web:d31f1b8176164e8903a03c"
};
```

### 2. VAPID Keyの取得

VAPID KeyはサーバーAPIから動的に取得します：

```javascript
// VAPID Keyを取得
async function getVapidKey() {
    const response = await fetch('/notifications/vapid-key');
    if (response.ok) {
        const data = await response.json();
        return data.vapid_key;
    }
    throw new Error('VAPID Key取得失敗');
}

// 使用例
const vapidKey = await getVapidKey();
```

### 3. Service Workerの設定

`firebase-messaging-sw.js`をプロジェクトのルートに配置します。
名前は適当で良いです。実装はテスト用実装を参考。

```javascript
importScripts('https://www.gstatic.com/firebasejs/11.1.0/firebase-app-compat.js');
importScripts('https://www.gstatic.com/firebasejs/11.1.0/firebase-messaging-compat.js');

firebase.initializeApp(firebaseConfig);

const messaging = firebase.messaging();

// バックグラウンドメッセージの処理
messaging.onBackgroundMessage((payload) => {
    const notificationTitle = payload.notification?.title || payload.data?.title || 'Kotohiro';
    const notificationOptions = {
        body: payload.notification?.body || payload.data?.body || '',
        icon: '/icon-192x192.png',
        badge: '/badge-72x72.png',
        data: payload.data || {}
    };

    return self.registration.showNotification(notificationTitle, notificationOptions);
});
```

### 4. 通知権限の取得とFCMトークンの取得

```javascript
import { initializeApp } from 'firebase/app';
import { getMessaging, getToken, onMessage } from 'firebase/messaging';

const app = initializeApp(firebaseConfig);
const messaging = getMessaging(app);

// 通知権限をリクエスト
async function requestNotificationPermission() {
    const permission = await Notification.requestPermission();
    if (permission === 'granted') {
        console.log('通知権限が許可されました');
        return true;
    }
    return false;
}

// FCMトークンを取得
async function getFCMToken() {
    try {
        // VAPID Keyを取得
        const vapidKey = await getVapidKey();

        // Service Workerを登録
        const registration = await navigator.serviceWorker.register('/firebase-messaging-sw.js');

        // トークンを取得
        const token = await getToken(messaging, {
            vapidKey: vapidKey,
            serviceWorkerRegistration: registration
        });

        return token;
    } catch (error) {
        console.error('FCMトークン取得エラー:', error);
        return null;
    }
}
```

### 5. デバイスの登録

FCMトークンを取得したら、APIでデバイスを登録します：

```javascript
async function registerDevice(fcmToken) {
    const formData = new FormData();
    formData.append('device_token', fcmToken);
    formData.append('platform', 'web'); // 'ios', 'android', 'web'
    formData.append('device_name', navigator.userAgent); // オプション
    formData.append('app_version', '1.0.0'); // オプション
    formData.append('os_version', navigator.platform); // オプション

    const response = await fetch('/notifications/devices', {
        method: 'POST',
        credentials: 'include', // 認証クッキーを含める
        body: formData
    });

    if (response.ok) {
        const device = await response.json();
        console.log('デバイス登録成功:', device);
        return device;
    }
    throw new Error('デバイス登録失敗');
}
```

### 6. フォアグラウンドでの通知受信

```javascript
// フォアグラウンドでメッセージを受信
onMessage(messaging, (payload) => {
    console.log('メッセージ受信:', payload);

    // カスタム通知を表示
    if (Notification.permission === 'granted') {
        const title = payload.notification?.title || payload.data?.title || 'Kotohiro';
        const body = payload.notification?.body || payload.data?.body || '';

        new Notification(title, {
            body: body,
            icon: '/icon-192x192.png'
        });
    }
});
```

## API エンドポイント

### VAPID Key取得

#### VAPID公開鍵取得

```
GET /notifications/vapid-key

レスポンス:
{
    "vapid_key": ""
}
```

### デバイス管理

#### デバイス登録/更新

```
POST /notifications/devices
Content-Type: multipart/form-data

パラメータ:
- device_token: string (必須) - FCMトークン
- platform: "ios" | "android" | "web" (必須)
- device_name: string (オプション)
- app_version: string (オプション)
- os_version: string XXXXXXXX(オプション)

レスポンス:
{
    "device_id": "uuid",
    "user_id": "uuid",
    "platform": "web",
    "enabled": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
}
```

#### デバイス一覧取得

```
GET /notifications/devices

レスポンス:
{
    "devices": [
        {
            "device_id": "uuid",
            "user_id": "uuid",
            "platform": "web",
            "enabled": true,
            "device_name": "Chrome on Windows",
            "last_active_at": "2024-01-01T00:00:00Z",
            "created_at": "2024-01-01T00:00:00Z",
            "updated_at": "2024-01-01T00:00:00Z"
        }
    ]
}
```

#### デバイストークン存在確認

```
GET /notifications/devices/exists?device_token={token}

レスポンス:
{
    "exists": true
}
```

#### デバイス削除

```
DELETE /notifications/devices/{deviceId}

レスポンス: 204 No Content
```

### 通知設定管理

#### 通知設定取得

```
GET /notifications/preferences

レスポンス:
{
    "push_notification_enabled": true
}
```

#### 通知設定更新

```
PUT /notifications/preferences
Content-Type: multipart/form-data

パラメータ:
- push_notification_enabled: boolean (オプション)

レスポンス:
{
    "push_notification_enabled": true
}
```

## 通知データ構造

プッシュ通知のペイロードは以下の構造を持ちます：

```javascript
{
    "notification": {
        "title": "通知タイトル",
        "body": "通知本文",
        "icon": "/icon-192x192.png"
    },
    "data": {
        // カスタムデータ
        "action": "open_talk_session", // アクションタイプ
        "talk_session_id": "uuid",      // トークセッションID
        "url": "/talk-sessions/uuid"    // 遷移先URL
    }
}
```

### アクションタイプ

- `open_talk_session`: トークセッション開始通知
- `open_talk_session_results`: トークセッション結果通知
- `test`: テスト通知

## テスト通知送信

### テストエンドポイント

開発・デバッグ時に通知機能をテストするためのエンドポイントが用意されています。

**注意**: このエンドポイントは開発環境（development）とローカル環境（local）でのみ利用可能です。

#### テスト通知送信

```
POST /notifications/test
Content-Type: multipart/form-data

パラメータ:
- title: string (オプション) - 通知タイトル（デフォルト: "Kotohiro テスト通知"）
- body: string (オプション) - 通知本文（デフォルト: "これはテスト通知です。時刻: HH:MM:SS"）
- device_id: string (オプション) - 特定のデバイスに送信する場合のデバイスID

レスポンス:
{
    "devices_count": 2,    // 送信対象デバイス数
    "success_count": 2     // 送信成功数
}
```

### テスト通知の使い方

```javascript
async function sendTestNotification() {
    const formData = new FormData();
    formData.append('title', 'カスタムタイトル');
    formData.append('body', 'カスタム本文');
    // 特定のデバイスにのみ送信する場合
    // formData.append('device_id', 'device-uuid');

    const response = await fetch('/notifications/test', {
        method: 'POST',
        credentials: 'include',
        body: formData
    });

    if (response.ok) {
        const result = await response.json();
        console.log(`${result.success_count}/${result.devices_count}台に送信成功`);
    }
}
```

### テスト通知のペイロード

テスト通知は以下のようなペイロードで送信されます：

```javascript
{
    "notification": {
        "title": "Kotohiro テスト通知",
        "body": "これはテスト通知です。時刻: 15:30:45"
    },
    "data": {
        "type": "test",
        "timestamp": "2024-01-01T15:30:45Z"
    }
}
```

## ベストプラクティス

### 1. エラーハンドリング

```javascript
try {
    const token = await getFCMToken();
    if (token) {
        await registerDevice(token);
    }
} catch (error) {
    console.error('通知設定エラー:', error);
    // ユーザーに適切なフィードバックを提供
}
```

### 2. トークンの更新

FCMトークンは定期的に更新される可能性があるため、以下を実装してください：

```javascript
import { onTokenRefresh } from 'firebase/messaging';

// トークン更新の監視
onTokenRefresh(messaging, async () => {
    const newToken = await getToken(messaging);
    if (newToken) {
        await registerDevice(newToken);
    }
});
```

### 3. Service Workerの更新

Service Workerのキャッシュ戦略を適切に設定してください：

```javascript
self.addEventListener('install', (event) => {
    self.skipWaiting();
});

self.addEventListener('activate', (event) => {
    event.waitUntil(clients.claim());
});
```

### 4. 認証状態の確認

デバイス登録前に必ず認証状態を確認してください：

```javascript
async function checkAuthAndRegisterDevice() {
    const authResponse = await fetch('/auth/token/info', {
        credentials: 'include'
    });

    if (!authResponse.ok) {
        // ログインページへリダイレクト
        return;
    }

    // デバイス登録処理
    const token = await getFCMToken();
    if (token) {
        await registerDevice(token);
    }
}
```

## トラブルシューティング

### 通知が届かない場合

1. **通知権限の確認**
   ```javascript
   console.log('通知権限:', Notification.permission);
   ```

2. **Service Workerの登録確認**
   ```javascript
   navigator.serviceWorker.getRegistrations().then(registrations => {
       console.log('Service Workers:', registrations);
   });
   ```

3. **FCMトークンの確認**
   - トークンが正しく取得できているか
   - デバイスがAPIに正しく登録されているか

4. **ブラウザコンソールでエラーを確認**
   - Service Workerのエラー
   - ネットワークエラー

### デバッグモード

開発時は以下のようにログを追加して問題を特定してください：

```javascript
// Service Worker内でのログ転送
self.addEventListener('message', (event) => {
    if (event.data && event.data.type === 'ENABLE_LOG_FORWARDING') {
        // ログをメインページに転送
        console.log = function(...args) {
            self.clients.matchAll().then(clients => {
                clients.forEach(client => {
                    client.postMessage({
                        type: 'SW_LOG',
                        message: args.join(' ')
                    });
                });
            });
        };
    }
});
```

## セキュリティ考慮事項

1. **HTTPSの使用**: Web Push通知はHTTPS環境でのみ動作します
2. **トークンの保護**: FCMトークンを安全に管理し、不要になったら削除してください
3. **認証の確認**: デバイス登録・削除時は必ず認証状態を確認してください
4. **CSPの設定**: Content Security PolicyでFirebaseドメインを許可してください

## テストページ

開発環境には、Web Push通知機能をテストするための専用ページが用意されています。

### アクセス方法

```
http://localhost:3000/static/test-webpush.html
```

### テストページの機能

1. **認証状態の確認**
   - 現在のログイン状態を表示
   - 開発用ログイン機能

2. **通知権限の管理**
   - 現在の通知権限状態を表示
   - 通知権限のリクエスト

3. **FCMトークンの取得**
   - Service Workerの登録
   - FCMトークンの取得と表示

4. **デバイス管理**
   - デバイスの登録
   - 登録済みデバイス一覧の表示

5. **テスト通知送信**
   - カスタムタイトル・本文でテスト通知を送信
   - リアルタイムでログを表示

### 使用手順

1. テストページにアクセス
2. 開発用ログインでユーザー認証
3. 通知権限をリクエスト
4. FCMトークンを取得
5. デバイスを登録
6. テスト通知を送信

### デバッグ情報

テストページでは、以下のデバッグ情報がリアルタイムで表示されます：

- Service Workerの登録状態
- FCMトークンの取得ログ
- API通信のログ
- 通知受信のログ
- エラー情報

## サーバー設定

### 環境変数

Web Push通知を動作させるために、以下の環境変数を設定する必要があります：

```bash
# FCM VAPID公開鍵
FCM_VAPID_KEY=BOeHeua3epcKCcuSGQ-Q0wNtf-rHKv-AF_E8b3oSL3fEasNkFOTYFO-fy7VgOtmvFFnaY9cJZ-ZwYNvKef4pYmk

# AWS Pinpoint設定（プッシュ通知送信用）
PINPOINT_APPLICATION_ID=your-pinpoint-app-id
AWS_ACCESS_KEY_ID=your-aws-access-key
AWS_SECRET_ACCESS_KEY=your-aws-secret-key

# 環境設定（local/development/production）
ENV=development
```

### VAPID Keyの生成

新しいVAPID Keyを生成する場合は、Firebaseコンソールから取得するか、以下のツールを使用できます：

1. Firebaseコンソール → プロジェクト設定 → Cloud Messaging → ウェブプッシュ証明書
2. または[web-push](https://www.npmjs.com/package/web-push)ライブラリを使用して生成

## 参考リンク

- [Firebase Cloud Messaging ドキュメント](https://firebase.google.com/docs/cloud-messaging/js/client)
- [Web Push API](https://developer.mozilla.org/docs/Web/API/Push_API)
- [Service Worker API](https://developer.mozilla.org/docs/Web/API/Service_Worker_API)
