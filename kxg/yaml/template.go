package yaml

import (
	"os"

	"gopkg.in/yaml.v2"

	l "github.com/micst/kxgctl/kxg/logging"
)

type TemplateName string
type TemplateNames []TemplateName

type TemplateAddress struct {
	Id          string            `yaml:"id"`
	Name        string            `yaml:"name"`
	Description *string           `yaml:"description",omitempty`
	DataType    *string           `yaml:"datatype",omitempty`
	Security    *string           `yaml:"security",omitempty`
	Central     *bool             `yaml:"central",omitempty`
	Disabled    *bool             `yaml:"disabled",omitempty`
	Bindings    AttributeBindings `yaml:"bindings"`
}
type TemplateAddresses []TemplateAddress

type Template struct {
	Name         TemplateName      `yaml:"template_name"`
	StartAddress *int              `yaml:"start_address",omitempty`
	DeviceName   *string           `yaml:"device_name",omitempty`
	Addresses    TemplateAddresses `yaml:"addresses"`
	Bindings     AttributeBindings `yaml:"bindings"`
	CanOverride  bool
	File         string
}

type Templates struct {
	template_definitions map[TemplateName]Template
}

func (t *TemplateName) String() string {
	return string(*t)
}

func (t *Templates) init() {
	if t.template_definitions == nil {
		t.template_definitions = make(map[TemplateName]Template)
	}
}

func (t *Templates) LoadYaml(file_name string, from_lib bool) {
	t.init()
	if yamlFile, err := os.ReadFile(file_name); err == nil {
		c := struct {
			Templates []Template `yaml:"templates"`
		}{}
		if err := yaml.Unmarshal(yamlFile, &c); err == nil {
			for _, template := range c.Templates {
				valid := true
				exist_file := ""
				if exist_tpl, exists := t.template_definitions[template.Name]; exists {
					valid = exist_tpl.CanOverride
					exist_file = exist_tpl.File
					if valid {
						l.Debug("overriding template \"" + exist_tpl.Name.String() + "\"")
					}
				}
				if valid {
					template.CanOverride = from_lib
					template.File = file_name
					t.template_definitions[template.Name] = template
				} else {
					l.Error("template \"" + template.Name.String() + "\" from file \"" + exist_file + "\" redefined in file " + file_name)
					l.Error("not adding")
				}
			}
		} else {
			l.Error("Unmarshal file " + file_name)
		}
	} else {
		l.Error("yamlFile.Get err")
	}
}

func (t *Templates) GetTemplateNames() TemplateNames {
	t.init()
	types := make(TemplateNames, len(t.template_definitions))
	i := 0
	for k := range t.template_definitions {
		types[i] = k
		i++
	}
	return types
}

func (t *Templates) GetTemplate(name TemplateName) Template {
	t.init()
	if template, exists := t.template_definitions[name]; exists {
		return template
	}
	return Template{}
}

func (t *Templates) TemplateExists(name TemplateName) bool {
	if _, exists := t.template_definitions[name]; exists {
		return true
	}
	l.Debug("template \"" + name.String() + "\" does not exist")
	return false
}
