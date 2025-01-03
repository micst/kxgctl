package kxg

import (
	"strconv"

	l "github.com/micst/kxgctl/kxg/logging"
)

func Quote(e string) string {
	return "\"" + e + "\""
}

func ShowContexts() {
	l.Info("Contexts:")
	l.Info("active: " + Quote(Data.Contexts.CurrentContext.String()))
	for _, ctx := range Data.Contexts.Contexts {
		l.Info(" - name:     " + Quote(ctx.Name.String()))
		l.Debug("   main:     " + Quote(ctx.MainAttribute.String()))
		l.Debug("   group:    " + Quote(ctx.RootGroup.String()))
		l.Debug("   middle:   " + Quote(ctx.MiddleAttribute.String()))
		l.Debug("   location: " + Quote(ctx.LocationAttribute.String()))
	}
}

func ShowAttributes() {
	l.Info("Attributes:")
	for _, attr_name := range Data.Attributes.GetAttributeNames() {
		l.Info(" - name: " + Quote(attr_name.String()))
		l.Info("   groups:")
		for _, group_name := range Data.Attributes.GetAttributeGroupKeys(attr_name) {
			l.Info("     - name: " + Quote(group_name.String()))
			group := Data.Attributes.GetAttributeGroup(attr_name, group_name)
			l.Debug("       members: ")
			for _, group_member := range group.Members {
				l.Debug("       - " + Quote(group_member.Key.String()))
				// TODO: GROUP BINDINGS!
			}
		}
		l.Info("   values:")
		for _, attr_key := range Data.Attributes.GetAttributeKeys(attr_name) {
			l.Info("    - id: " + Quote(attr_key.String()))
			value := Data.Attributes.GetAttributeValue(attr_name, attr_key)
			l.Debug("      text: " + Quote(value.Text))
			// TODO: GROUP BINDINGS!
		}
		// TODO: GROUP BINDINGS!
	}
}

func ShowTemplates() {
	l.Info("Templates:")
	for _, tpl_name := range Data.Templates.GetTemplateNames() {
		l.Info(" - name: " + Quote(tpl_name.String()))
		tpl := Data.Templates.GetTemplate(tpl_name)
		l.Debug("   file: " + Quote(tpl.File))
		l.Debug("   bindings: " + Quote(tpl.File))
		for _, binding := range tpl.Bindings {
			l.Debug("    - attr: " + Quote(binding.Name.String()))
			l.Debug("      key:  " + Quote(binding.Key.String()))
		}
		l.Debug("   addresses: " + Quote(tpl.File))
		for _, address := range tpl.Addresses {
			l.Debug("    - id: " + Quote(address.Id))

			for _, binding := range address.Bindings {
				l.Debug("    - attr: " + Quote(binding.Name.String()))
				l.Debug("      key:  " + Quote(binding.Key.String()))
			}
		}
	}
}

func ShowDevices() {
	l.Info("Devices:")
	for device_idx := 0; device_idx < Data.Devices.GetDeviceCount(); device_idx++ {
		device := Data.Devices.GetDevice(device_idx)
		l.Info(" - template: " + Quote(device.Name.String()))
		l.Info("   start:    " + strconv.Itoa(*device.StartAddress))
	}
}
