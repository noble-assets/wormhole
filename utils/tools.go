//go:build tools

// This file uses the recommended method for tracking developer tools in a Go module.
//
// REF: https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
package tools

import (
	_ "github.com/incu6us/goimports-reviser/v3"
	_ "mvdan.cc/gofumpt"
)
