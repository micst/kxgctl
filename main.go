package main

import (
	"github.com/micst/kxgctl/cmd"
)

//go:generate go generate ./kxg
//go:generate pwsh -NoProfile -File "tag.ps1"

func main() {
	cmd.Execute()
}
