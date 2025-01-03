package yaml

import (
	"os"
	"slices"

	"gopkg.in/yaml.v2"

	l "github.com/micst/kxgctl/kxg/logging"
)

type ContextName string
type ContextNames []ContextName

type ContextValue struct {
	Name              ContextName       `yaml:"name"`
	MainAttribute     AttributeName     `yaml:"main_attribute"`
	MiddleAttribute   AttributeName     `yaml:"middle_attribute"`
	RootGroup         AttributeGroupKey `yaml:"root_group"`
	LocationAttribute AttributeName     `yaml:"location_attribute"`
}
type ContextValues []ContextValue

type Contexts struct {
	Contexts       ContextValues `yaml:"contexts"`
	CurrentContext ContextName   `yaml:"current-context"`
	File           string
}

func (cn *ContextName) String() string {
	return string(*cn)
}

func (d *Contexts) init() {
	if d.Contexts == nil {
		d.Contexts = make(ContextValues, 0)
	}
}

func (d *Contexts) LoadYaml(file_name string) {
	d.init()
	if yamlFile, err := os.ReadFile(file_name); err == nil {
		if err := yaml.Unmarshal(yamlFile, &d); err == nil {
			d.File = file_name
		} else {
			l.Error("Could not unmarshal file \"" + file_name + "\"")
		}
	} else {
		l.Error("Could not read file \"" + file_name + "\"")
	}
}

func (d *Contexts) GetContextNames() ContextNames {
	d.init()
	res := ContextNames{}
	for _, c := range d.Contexts {
		res = append(res, c.Name)
	}
	return res
}

func (d *Contexts) GetContext(name ContextName) ContextValue {
	d.init()
	for _, c := range d.Contexts {
		if c.Name == name {
			return c
		}
	}
	l.Error("context " + name.String() + " not found")
	return ContextValue{}
}

func (d *Contexts) GetCurrentContext() ContextValue {
	d.init()
	if len(d.Contexts) > 0 {
		for _, c := range d.Contexts {
			if c.Name == d.CurrentContext {
				return c
			}
		}
	}
	return ContextValue{}
}

func (d *Contexts) ContextExists(name ContextName) bool {
	for _, ctx := range d.Contexts {
		if ctx.Name == name {
			return true
		}
	}
	l.Debug("context \"" + name.String() + "\" does not exist")
	return false
}

func (d *ContextValue) Validate() bool {
	res := true
	res = res && d.MainAttribute != ""
	res = res && d.MiddleAttribute != ""
	res = res && d.LocationAttribute != ""
	res = res && d.RootGroup != ""
	return res
}

func (d *Contexts) Validate() bool {
	res := true
	names := ContextNames{}
	for _, ctx := range d.Contexts {
		if slices.Contains(names, ctx.Name) {
			l.Debug("context \"" + ctx.Name.String() + "\" exists multiple times")
			res = false
		}
		names = append(names, ctx.Name)
	}
	res = res && d.ContextExists(d.CurrentContext)
	return res
}
