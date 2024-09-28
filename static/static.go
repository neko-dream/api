package static

import "embed"

//go:embed *.yaml
var Static embed.FS
