package web

import "embed"

//go:embed all:web/*
var Static embed.FS
