package kxg

import (
	"github.com/micst/kxgctl/kxg/kxml"
	l "github.com/micst/kxgctl/kxg/logging"
	"github.com/micst/kxgctl/kxg/yaml"
)

var (
	Data = struct {
		Attributes yaml.Attributes
		Templates  yaml.Templates
		Devices    yaml.Devices
		Contexts   yaml.Contexts
		Tree       GroupTree
		Document   kxml.GroupAddressDocument
	}{}

	Ws  Workspace
	Lib Workspace

	Args = struct {
		Force       bool
		SkipLibrary bool
		SkipVerify  bool
		Dry         bool
		XmlOut      string
		ContextName string
	}{}
)

func LoadLibrary() {
	if Args.SkipLibrary {
		l.Debug("skipping library loading")
		return
	}
	l.Debug("loading library from \"" + Lib.Directory + "\"")
	Lib.Load(true)
}

func LoadWorkspace() {
	l.Debug("loading workspace from \"" + Ws.Directory + "\"")
	Ws.Load(false)
}
