package pkg

import "embed"

//go:embed static
var Static embed.FS

//go:embed views
var Views embed.FS
