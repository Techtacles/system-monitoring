package dashboard

import "embed"

// WebAssets embeds the web directory
//
//go:embed web/*
var WebAssets embed.FS
