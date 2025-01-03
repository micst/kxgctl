package kxg

import "github.com/micst/kxgctl/kxg/yaml"

type BindingValues struct {
	Main     yaml.AttributeKey
	Middle   yaml.AttributeKey
	Location yaml.AttributeKey
}

type Meta struct {
	Binding     BindingValues
	Number      int
	Name        string
	Description string
	Security    string
}

type Address struct {
	Meta       Meta
	DeviceName string
	DataType   string
	Central    bool
	Disabled   bool
}
type Addresses map[string]*Address

type MiddleGroup struct {
	Meta         Meta
	MaxAddress   int
	Numbers      []int
	AddressNames []string
	Addresses    Addresses
}
type MiddleGroups map[yaml.AttributeKey]*MiddleGroup

type MainGroup struct {
	Meta         Meta
	GroupDef     yaml.AttributeGroup
	MiddleGroups MiddleGroups
}
type MainGroups map[yaml.AttributeKey]*MainGroup

type GroupTree struct {
	GroupDef   yaml.AttributeGroup
	MainGroups MainGroups
}

func (b *BindingValues) Set(main yaml.AttributeName, middle yaml.AttributeName, location yaml.AttributeName, bindings yaml.AttributeBindings) {
	for _, binding := range bindings {
		if binding.Name == main {
			b.Main = binding.Key
		}
		if binding.Name == middle {
			b.Middle = binding.Key
		}
		if binding.Name == location {
			b.Location = binding.Key
		}
	}
}

func (b *BindingValues) Key() string {
	return string(b.Main) + "-" + string(b.Middle)
}

func (a *Address) Set(cfg yaml.TemplateAddress) {
	if a.Meta.Name == cfg.Name || a.Meta.Name == "" {
		a.Meta.Name = cfg.Name
		if cfg.DataType != nil {
			a.DataType = *cfg.DataType
		}
		if cfg.Central != nil {
			a.Central = *cfg.Central
		}
		if cfg.Security != nil {
			a.Meta.Security = *cfg.Security
		}
		if cfg.Description != nil {
			a.Meta.Description = *cfg.Description
		}
		if cfg.Disabled != nil {
			a.Disabled = *cfg.Disabled
		}
	}
}
