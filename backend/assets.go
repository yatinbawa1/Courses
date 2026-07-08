package backend

import "embed"

//go:embed dist/*
var FrontendFS embed.FS
