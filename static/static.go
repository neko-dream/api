package static

import "embed"

//go:embed oas/*
var Oas embed.FS

//go:embed admin-ui/*
var AdminUI embed.FS
