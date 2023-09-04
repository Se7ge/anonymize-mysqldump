package embed

import (
	_ "embed"
)

// DefaultConfig contains the contents of the default config file.
//
//go:embed files/config.default.json
var DefaultConfig string
