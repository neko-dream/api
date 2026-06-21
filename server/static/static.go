package static

import "embed"

//go:embed oas/*
var Oas embed.FS

//go:embed admin-ui/*
var AdminUI embed.FS

//go:embed test-webpush.html firebase-messaging-sw.js
var WebPushTest embed.FS
