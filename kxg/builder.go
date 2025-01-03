package kxg

import (
	"encoding/xml"
	"math"
	"os"
	"slices"
	"strconv"

	l "github.com/micst/kxgctl/kxg/logging"

	"github.com/micst/kxgctl/kxg/kxml"
	"github.com/micst/kxgctl/kxg/yaml"
)

func BuildTree() {
	main_group_ids := yaml.AttributeGroupMembers{}
	main_groups := make(map[yaml.AttributeKey]*MainGroup)
	main_group_num := 0
	context := Data.Contexts.GetCurrentContext()
	root_group_cfg := Data.Attributes.GetAttributeGroup(context.MainAttribute, context.RootGroup)
	// iterate over all main groups
	for _, main_group_id := range root_group_cfg.Members {
		middle_group_ids := yaml.AttributeGroupMembers{}
		middle_groups := make(map[yaml.AttributeKey]*MiddleGroup)
		middle_group_set := Data.Attributes.GetMiddleGroupForAttributeValue(context.MiddleAttribute, context.MainAttribute, main_group_id.Key)
		middle_group_num := 0
		// iterate over middle groups for main group
		for _, middle_group_id := range middle_group_set.Members {
			middle_group_cfg := Data.Attributes.GetAttributeValue(context.MiddleAttribute, middle_group_id.Key)
			middle_group_ids = append(middle_group_ids, middle_group_id)
			middle_groups[middle_group_id.Key] = &MiddleGroup{
				AddressNames: []string{},
				Addresses:    make(map[string]*Address),
				MaxAddress:   -1,
				Numbers:      []int{},
				Meta: Meta{
					Binding: BindingValues{
						Main:   main_group_id.Key,
						Middle: middle_group_id.Key,
					},
					Name:        middle_group_cfg.Text,
					Description: middle_group_cfg.Description,
					Number:      middle_group_num,
					Security:    "",
				},
			}
			middle_group_num++
		}
		main_group_cfg := Data.Attributes.GetAttributeValue(context.MainAttribute, main_group_id.Key)
		main_group_ids = append(main_group_ids, main_group_id)
		main_groups[main_group_id.Key] = &MainGroup{
			MiddleGroups: middle_groups,
			GroupDef: yaml.AttributeGroup{
				GroupId: middle_group_set.GroupId,
				Members: middle_group_ids,
			},
			Meta: Meta{
				Binding: BindingValues{
					Main:   main_group_id.Key,
					Middle: "",
				},
				Name:        main_group_cfg.Text,
				Description: main_group_cfg.Description,
				Number:      main_group_num,
				Security:    "",
			},
		}
		main_group_num++
	}
	Data.Tree = GroupTree{
		GroupDef: yaml.AttributeGroup{
			GroupId: context.RootGroup,
			Members: main_group_ids,
		},
		MainGroups: main_groups,
	}
}

func BuildAddresses() {
	ctx := Data.Contexts.GetCurrentContext()
	// iterate over all devices
	for device_idx := 0; device_idx < Data.Devices.GetDeviceCount(); device_idx++ {
		device_cfg := Data.Devices.GetDevice(device_idx)
		device_template := Data.Templates.GetTemplate(device_cfg.Name)
		binding_list_device := append(device_template.Bindings, device_cfg.Bindings...)
		device_name := ""
		if device_template.DeviceName != nil {
			device_name = *device_template.DeviceName
		}
		if device_cfg.DeviceName != nil {
			device_name = *device_cfg.DeviceName
		}
		// collect address configurations
		device_addresses := []Address{}
		for _, t_address := range device_template.Addresses {
			address := Address{
				DeviceName: device_name,
				Disabled:   false,
				Meta: Meta{
					Binding: BindingValues{},
					Number:  -1,
				},
			}
			address.Meta.Binding.Set(ctx.MainAttribute, ctx.MiddleAttribute, ctx.LocationAttribute,
				append(binding_list_device, t_address.Bindings...))
			address.Set(t_address)
			for _, d_address := range device_cfg.Addresses {
				address.Set(d_address)
			}
			device_addresses = append(device_addresses, address)
		}
		// initialize address indices for all group bindings
		ga_indices := map[string]int{}
		for _, address_cfg := range device_addresses {
			binding := address_cfg.Meta.Binding
			key := binding.Key()
			middle_group := Data.Tree.MainGroups[binding.Main].MiddleGroups[binding.Middle]
			// initialize with current next free address
			ga_indices[key] = middle_group.MaxAddress + 1
			// override if address given in device config
			if device_cfg.StartAddress != nil {
				start := *device_cfg.StartAddress
				if start > middle_group.MaxAddress {
					ga_indices[key] = *device_cfg.StartAddress
				} else {
					current := ga_indices[key]
					l.Debug("device: " + device_name + " address: " + address_cfg.Meta.Name + ", cannot set address")
					l.Debug(" - configured start address " + strconv.Itoa(start) + " is already taken, shifting to " + strconv.Itoa(current))
				}
			}
		}
		// create addresses
		for _, address_cfg := range device_addresses {
			binding := address_cfg.Meta.Binding
			key := binding.Key()
			middle_group := Data.Tree.MainGroups[binding.Main].MiddleGroups[binding.Middle]
			// address value handling
			ga := ga_indices[key]
			ga_indices[key] = ga + 1
			middle_group.Numbers = append(middle_group.Numbers, ga)
			middle_group.MaxAddress = slices.Max(middle_group.Numbers)
			slices.Max(middle_group.Numbers)
			if !address_cfg.Disabled {
				// compute address name
				address_separator := "-"
				location_attr_name := Data.Contexts.GetCurrentContext().LocationAttribute
				address_name := Data.Attributes.GetAttributeValue(location_attr_name, address_cfg.Meta.Binding.Location).Text
				if address_cfg.DeviceName != "" {
					address_name += address_separator
					address_name += address_cfg.DeviceName
				}
				address_name += address_separator
				address_name += address_cfg.Meta.Name
				// generate final address struct
				middle_group.AddressNames = append(middle_group.AddressNames, address_name)
				middle_group.Addresses[address_name] = &Address{
					DataType: address_cfg.DataType,
					Central:  address_cfg.Central,
					Meta: Meta{
						Binding:     BindingValues{},
						Number:      ga,
						Name:        address_name,
						Description: address_cfg.Meta.Description,
						Security:    address_cfg.Meta.Security,
					},
				}
			}
		}
	}
}

func BuildDocument() {
	Data.Document = kxml.GroupAddressDocument{
		XMLName: xml.Name{Local: kxml.XmlNameDocument, Space: ""},
		Xmlns:   kxml.XmlNamespace,
	}
	range_main := int(math.Pow(2, 11))
	range_middle := int(math.Pow(2, 8))
	for _, main_group_id := range Data.Tree.GroupDef.Members {
		main_group := Data.Tree.MainGroups[main_group_id.Key]
		main_start_addr := main_group.Meta.Number * range_main
		main_end_addr := main_start_addr + range_main - 1
		main_start_addr_final := main_start_addr
		if main_start_addr_final == 0 {
			main_start_addr_final = 1
		}
		doc_main := kxml.GroupRangeMain{
			XMLName:      xml.Name{Local: kxml.XmlNameGroup},
			RangeStart:   strconv.Itoa(main_start_addr_final),
			RangeEnd:     strconv.Itoa(main_end_addr),
			Name:         main_group.Meta.Name,
			Description:  main_group.Meta.Description,
			MiddleGroups: []kxml.GroupRangeMiddle{},
		}
		for _, middle_group_id := range main_group.GroupDef.Members {
			middle_group := main_group.MiddleGroups[middle_group_id.Key]
			middle_start_addr := main_start_addr + (range_middle * middle_group.Meta.Number)
			middle_end_addr := middle_start_addr + range_middle - 1
			middle_start_addr_final := middle_start_addr
			if middle_start_addr_final == 0 {
				middle_start_addr_final = 1
			}
			doc_middle := kxml.GroupRangeMiddle{
				XMLName:     xml.Name{Local: kxml.XmlNameGroup},
				RangeStart:  strconv.Itoa(middle_start_addr_final),
				RangeEnd:    strconv.Itoa(middle_end_addr),
				Name:        middle_group.Meta.Name,
				Description: middle_group.Meta.Description,
			}
			for _, address_name := range middle_group.AddressNames {
				address := middle_group.Addresses[address_name]
				address_value := strconv.Itoa(main_group.Meta.Number)
				address_value += "/"
				address_value += strconv.Itoa(middle_group.Meta.Number)
				address_value += "/"
				address_value += strconv.Itoa(address.Meta.Number)
				doc_address := kxml.GroupAddress{
					XMLName:     xml.Name{Local: kxml.XmlNameAddress},
					Name:        address.Meta.Name,
					Address:     address_value,
					Description: address.Meta.Description,
					Security:    address.Meta.Security,
					DPTs:        address.DataType,
				}
				doc_middle.Addresses = append(doc_middle.Addresses, doc_address)
			}
			doc_main.MiddleGroups = append(doc_main.MiddleGroups, doc_middle)
		}
		Data.Document.MainGroups = append(Data.Document.MainGroups, doc_main)
	}
}

func Validate() bool {
	if Args.SkipVerify {
		l.Debug("skipping validation of config")
		return true
	}
	res := true
	res = res && Data.Contexts.Validate()
	// validate attributes
	res = res && Data.Attributes.Validate()
	// validate template bindings
	for _, template_name := range Data.Templates.GetTemplateNames() {
		template := Data.Templates.GetTemplate(template_name)
		for _, binding := range template.Bindings {
			res = res && Data.Attributes.AttributeKeyExists(binding.Name, binding.Key)
		}
		for _, address := range template.Addresses {
			for _, binding := range address.Bindings {
				res = res && Data.Attributes.AttributeKeyExists(binding.Name, binding.Key)
			}
		}
	}
	// validate device bindings
	for i := 0; i < Data.Devices.GetDeviceCount(); i++ {
		device := Data.Devices.GetDevice(i)
		res = res && Data.Templates.TemplateExists(device.Name)
		for _, binding := range device.Bindings {
			res = res && Data.Attributes.AttributeKeyExists(binding.Name, binding.Key)
		}
		for _, address := range device.Addresses {
			for _, binding := range address.Bindings {
				res = res && Data.Attributes.AttributeKeyExists(binding.Name, binding.Key)
			}
		}
	}
	// validate context bindings
	for _, ctx := range Data.Contexts.Contexts {
		res = res && Data.Attributes.AttributeExists(ctx.MainAttribute)
		res = res && Data.Attributes.AttributeExists(ctx.MiddleAttribute)
		res = res && Data.Attributes.AttributeGroupExists(ctx.MainAttribute, ctx.RootGroup)
		res = res && Data.Attributes.AttributeExists(ctx.LocationAttribute)
	}
	// validate group sizes
	ctx := Data.Contexts.GetCurrentContext()
	root_grp := Data.Attributes.GetAttributeGroup(ctx.MainAttribute, ctx.RootGroup)
	if len(root_grp.Members) > 32 {
		l.Debug("main group too large")
		res = false
	}
	for _, main_group_key := range root_grp.Members {
		if Data.Attributes.GetMiddleGroupSizeForAttributeValue(ctx.MiddleAttribute, ctx.MainAttribute, main_group_key.Key) > 8 {
			l.Debug("middle group too large")
			res = false
		}
	}
	if !res {
		l.Error("validation error")
		os.Exit(1)
	}
	return res
}
