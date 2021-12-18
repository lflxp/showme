package pkg

import "embed"

//go:embed static/*
var Static embed.FS
