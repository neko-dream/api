package static

import "embed"

//go:embed *.yaml
var Static embed.FS

//go:embed index.html
var IndexHTML embed.FS
