package yaml

import (
	"os"

	"gopkg.in/yaml.v2"

	l "github.com/micst/kxgctl/kxg/logging"
)

type DeviceGroup struct {
	Bindings AttributeBindings `yaml:"bindings"`
	Devices  []Template        `yaml:"devices"`
}
type DeviceGroups []DeviceGroup

type Devices struct {
	device_definitions []Template
}

func (d *Devices) init() {
	if d.device_definitions == nil {
		d.device_definitions = make([]Template, 0)
	}
}

func (d *Devices) LoadYaml(file_name string) {
	d.init()
	if yamlFile, err := os.ReadFile(file_name); err == nil {
		c := struct {
			Groups DeviceGroups `yaml:"groups"`
		}{}
		if err := yaml.Unmarshal(yamlFile, &c); err == nil {
			for _, group := range c.Groups {
				for _, device := range group.Devices {
					device.Bindings = append(group.Bindings, device.Bindings...)
					d.device_definitions = append(d.device_definitions, device)
				}
			}
		} else {
			l.Error("Could not unmarshal file \"" + file_name + "\"")
		}
	} else {
		l.Error("Cound not read file \"" + file_name + "\"")
	}
}

func (d *Devices) GetDeviceCount() int {
	return len(d.device_definitions)
}

func (d *Devices) GetDevice(pos int) Template {
	d.init()
	return d.device_definitions[pos]
}
