// Firebase Messaging Service Worker
importScripts('https://www.gstatic.com/firebasejs/11.1.0/firebase-app-compat.js');
importScripts('https://www.gstatic.com/firebasejs/11.1.0/firebase-messaging-compat.js');

// ログ転送機能
let logForwardingEnabled = false;

// 元のconsole.logを保存
const originalLog = console.log;

// console.logをオーバーライド
console.log = function(...args) {
    // 元のconsole.logを呼ぶ
    originalLog.apply(console, args);

    // ログ転送が有効な場合、メインページに送信
    if (logForwardingEnabled) {
        self.clients.matchAll().then(clients => {
            clients.forEach(client => {
                client.postMessage({
                    type: 'SW_LOG',
                    message: args.join(' ')
                });
            });
        });
    }
};

// メインページからのメッセージを受信
self.addEventListener('message', (event) => {
    if (event.data && event.data.type === 'ENABLE_LOG_FORWARDING') {
        logForwardingEnabled = true;
        console.log('ログ転送が有効になりました');
    }
});

console.log('Service Worker初期化開始');

firebase.initializeApp({
    apiKey: "AIzaSyBEQ-gqPCUMyyBZ-YP10t_0X0D01BGPCgQ",
    authDomain: "kotohiro-dev-c57a8.firebaseapp.com",
    projectId: "kotohiro-dev-c57a8",
    storageBucket: "kotohiro-dev-c57a8.firebasestorage.app",
    messagingSenderId: "1023394330741",
    appId: "1:1023394330741:web:d31f1b8176164e8903a03c"
});

const messaging = firebase.messaging();
console.log('Firebase Messaging初期化完了');

// バックグラウンドメッセージのハンドリング
messaging.onBackgroundMessage((payload) => {
    console.log('[firebase-messaging-sw.js] Received background message ', JSON.stringify(payload, null, 2));

    // PinpointからのGCMメッセージはdataフィールドに通知内容が含まれる
    const notificationTitle = payload.notification?.title || payload.data["pinpoint.notification.title"] || 'Kotohiro';
    const notificationOptions = {
        body: payload.notification["pinpoint.notification.body"] || payload.data?.body || '',
        icon: payload.notification["pinpoint.notification.icon"] || payload.data?.icon || '/icon-192x192.png',
        badge: '/badge-72x72.png',
        data: payload.data || {},
        actions: [
            {
                action: 'open',
                title: '開く'
            },
            {
                action: 'close',
                title: '閉じる'
            }
        ]
    };

    // カスタムデータに基づいて通知をカスタマイズ
    if (payload.data) {
        // talk_session_idがある場合
        if (payload.data.talk_session_id) {
            notificationOptions.tag = `talk-session-${payload.data.talk_session_id}`;
            notificationOptions.data.url = `/talk-sessions/${payload.data.talk_session_id}`;
        }

        // アクションタイプによる分岐
        switch (payload.data.action) {
            case 'open_talk_session':
                notificationOptions.requireInteraction = true;
                break;
            case 'open_talk_session_results':
                notificationOptions.vibrate = [200, 100, 200];
                break;
        }
    }

    return self.registration.showNotification(notificationTitle, notificationOptions);
});

// Pushイベントの直接処理（デバッグ用）
self.addEventListener('push', (event) => {
    console.log('[SW] Push event received:', event);
    if (event.data) {
        try {
            const data = event.data.json();
            console.log('[SW] Push data:', JSON.stringify(data, null, 2));

            // FirebaseのonBackgroundMessageが処理しない場合のフォールバック
            // dataのみのメッセージの場合、手動で通知を表示
            if (data && !data.notification && data.data) {
                const notificationTitle = data.data["pinpoint.notification.title"] || 'Kotohiro';
                const notificationOptions = {
                    body: data.data["pinpoint.notification.body"] || '',
                    icon: data.data.icon || '/icon-192x192.png',
                    badge: '/badge-72x72.png',
                    data: data.data || {}
                };

                console.log('[SW] Showing notification manually:', notificationTitle, notificationOptions);
                event.waitUntil(
                    self.registration.showNotification(notificationTitle, notificationOptions)
                );
            }
        } catch (e) {
            console.log('[SW] Push text:', event.data.text());
        }
    }
});

// 通知クリックイベント
self.addEventListener('notificationclick', (event) => {
    console.log('[firebase-messaging-sw.js] Notification click received.');

    event.notification.close();

    // URLを開く
    if (event.notification.data && event.notification.data.url) {
        event.waitUntil(
            clients.openWindow(event.notification.data.url)
        );
    } else if (event.action === 'open') {
        event.waitUntil(
            clients.openWindow('/')
        );
    }
});
