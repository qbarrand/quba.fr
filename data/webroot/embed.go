package webroot

import "embed"

//go:embed *
var WebRoot embed.FS
