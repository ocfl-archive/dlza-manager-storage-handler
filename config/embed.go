package config

import "embed"

//go:embed storagehandler.toml
var ConfigFS embed.FS
