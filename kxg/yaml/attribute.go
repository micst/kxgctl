package yaml

import (
	"os"
	"slices"

	"gopkg.in/yaml.v2"

	l "github.com/micst/kxgctl/kxg/logging"
)

type AttributeName string
type AttributeNames []AttributeName

type AttributeKey string
type AttributeKeys []AttributeKey

type AttributeGroupKey string
type AttributeGroupKeys []AttributeGroupKey

type AttributeGroupMember struct {
	Key       AttributeKey       `yaml:"id"`
	SubGroups AttributeSubGroups `yaml:"middle_groups"`
}
type AttributeGroupMembers []AttributeGroupMember

type AttributeGroup struct {
	GroupId AttributeGroupKey     `yaml:"id"`
	Members AttributeGroupMembers `yaml:"members"`
}
type AttributeGroups []AttributeGroup

type AttributeSubGroup struct {
	SubGroupAttributeName AttributeName     `yaml:"name"`
	SubGroupId            AttributeGroupKey `yaml:"id"`
}
type AttributeSubGroups []AttributeSubGroup

type AttributeValue struct {
	ValueId     AttributeKey       `yaml:"id"`
	Text        string             `yaml:"text"`
	Description string             `yaml:"description"`
	Icon        string             `yaml:"icon"`
	HaName      string             `yaml:"ha_name"`
	Parent      AttributeKey       `yaml:"parent"`
	Root        bool               `yaml:"root"`
	SubGroups   AttributeSubGroups `yaml:"middle_groups"`
}
type AttributeValues []AttributeValue

type AttributeConfig struct {
	Name        AttributeName      `yaml:"name"`
	SubGroups   AttributeSubGroups `yaml:"middle_groups"`
	Values      AttributeValues    `yaml:"values"`
	Groups      AttributeGroups    `yaml:"groups"`
	CanOverride bool
	File        string
	Position    int
}

type AttributeBinding struct {
	Name AttributeName `yaml:"attribute_name"`
	Key  AttributeKey  `yaml:"attribute_key"`
}
type AttributeBindings []AttributeBinding

type Attributes struct {
	attribute_definitions map[AttributeName]AttributeConfig
}

func (a *AttributeName) String() string {
	return string(*a)
}

func (a *AttributeKey) String() string {
	return string(*a)
}

func (a *AttributeGroupKey) String() string {
	return string(*a)
}

func (a *Attributes) init() {
	if a.attribute_definitions == nil {
		a.attribute_definitions = make(map[AttributeName]AttributeConfig)
	}
}

func (a *Attributes) getAttributeValueConfig(name AttributeName, value_id AttributeKey) (AttributeValue, bool) {
	a.init()
	res := AttributeValue{}
	valid := false
	if att, ok := a.attribute_definitions[name]; ok {
		for _, value := range att.Values {
			if value.ValueId == value_id {
				res = value
				valid = true
			}
		}
	}
	return res, valid
}

func (a *Attributes) LoadYaml(file_name string, from_lib bool) {
	a.init()
	if yamlFile, err := os.ReadFile(file_name); err == nil {
		c := AttributeConfig{}
		if err := yaml.Unmarshal(yamlFile, &c); err == nil {
			valid := true
			pos := len(a.attribute_definitions)
			if a, exists := a.attribute_definitions[c.Name]; exists {
				valid = a.CanOverride
				if valid {
					pos = a.Position
					l.Debug("overriding attribute \"" + a.Name.String() + "\"")
				}
			}
			if valid {
				c.File = file_name
				c.CanOverride = from_lib
				c.Position = pos
				a.attribute_definitions[c.Name] = c
			} else {
				l.Error("cannot load attribute file \"" + file_name + "\"")
				l.Error("attribute with name \"" + c.Name.String() + "\" already exists in file \"" + a.attribute_definitions[c.Name].File + "\"")
				os.Exit(1)
			}
		} else {
			l.Error("Could not unmarshal file " + file_name)
		}
	} else {
		l.Error("Could not load file " + file_name)
	}
}

func (a *Attributes) GetAttributeNames() AttributeNames {
	a.init()
	types := make(AttributeNames, len(a.attribute_definitions))
	i := 0
	for k := range a.attribute_definitions {
		types[i] = k
		i++
	}
	return types
}

func (a *Attributes) GetAttributeKeys(name AttributeName) AttributeKeys {
	a.init()
	res := AttributeKeys{}
	if a, ok := a.attribute_definitions[name]; ok {
		for _, value := range a.Values {
			res = append(res, value.ValueId)
		}
	}
	return res
}

func (a *Attributes) GetAttributeGroupKeys(name AttributeName) AttributeGroupKeys {
	a.init()
	res := AttributeGroupKeys{}
	if a, ok := a.attribute_definitions[name]; ok {
		for _, group := range a.Groups {
			res = append(res, group.GroupId)
		}
	}
	return res
}

func (a *Attributes) GetAttributeValue(name AttributeName, value_id AttributeKey) AttributeValue {
	a.init()
	if value, ok := a.getAttributeValueConfig(name, value_id); ok {
		res := AttributeValue{}
		res.ValueId = value_id
		res.Description = value.Description
		res.HaName = value.HaName
		res.Icon = value.Icon
		// build final text
		cVal := value
		text := ""
		for {
			if len(text) > 0 {
				text = "-" + text
			}
			text = cVal.Text + text
			if cVal.Root || len(cVal.Parent) == 0 {
				break
			}
			cVal, ok = a.getAttributeValueConfig(name, cVal.Parent)
		}
		res.Text = text
		// build subgroup element
		res.SubGroups = value.SubGroups
		middle_group_attr_names := make(AttributeNames, len(res.SubGroups))
		for _, grp := range res.SubGroups {
			middle_group_attr_names = append(middle_group_attr_names, grp.SubGroupAttributeName)
		}
		default_middle_groups := a.attribute_definitions[name].SubGroups
		for _, grp := range default_middle_groups {
			if !slices.Contains(middle_group_attr_names, grp.SubGroupAttributeName) {
				res.SubGroups = append(res.SubGroups, grp)
			}
		}
		return res
	}
	return AttributeValue{}
}

func (a *Attributes) GetAttributeGroup(name AttributeName, group_id AttributeGroupKey) AttributeGroup {
	a.init()
	res := AttributeGroup{}
	if a, ok := a.attribute_definitions[name]; ok {
		for _, group := range a.Groups {
			if group.GroupId == group_id {
				//res.AttributeName = name
				res.GroupId = group_id
				res.Members = group.Members
			}
		}
	}
	return res
}

func (a *Attributes) GetAttributeGroupSize(name AttributeName, group_id AttributeGroupKey) int {
	a.init()
	if a, ok := a.attribute_definitions[name]; ok {
		for _, group := range a.Groups {
			if group.GroupId == group_id {
				return len(group.Members)
			}
		}
	}
	return 0
}

func (a *Attributes) GetMiddleGroupForAttributeValue(middle_group_attribute AttributeName, attr_name AttributeName, value_id AttributeKey) AttributeGroup {
	attr := a.attribute_definitions[attr_name]
	val := a.GetAttributeValue(attr_name, value_id)
	grps := append(val.SubGroups, attr.SubGroups...)
	for _, group := range grps {
		if group.SubGroupAttributeName == middle_group_attribute {
			return a.GetAttributeGroup(group.SubGroupAttributeName, group.SubGroupId)
		}
	}
	return AttributeGroup{}
}

func (a *Attributes) GetMiddleGroupSizeForAttributeValue(middle_group_attribute AttributeName, main_name AttributeName, main_key AttributeKey) int {
	attr := a.attribute_definitions[main_name]
	val := a.GetAttributeValue(main_name, main_key)
	grps := append(val.SubGroups, attr.SubGroups...)
	for _, group := range grps {
		if group.SubGroupAttributeName == middle_group_attribute {
			return a.GetAttributeGroupSize(group.SubGroupAttributeName, group.SubGroupId)
		}
	}
	return 0
}

func (a *Attributes) DefaultContext() ContextValue {
	if len(a.attribute_definitions) >= 2 {
		ctx := ContextValue{
			Name: "default",
		}
		for _, def := range a.attribute_definitions {
			if def.Position == 0 {
				ctx.MainAttribute = def.Name
				main_groups := a.GetAttributeGroupKeys(def.Name)
				if len(main_groups) > 0 {
					ctx.RootGroup = main_groups[0]
				} else {
					l.Error("could not auto detect root group")
				}
			}
			if def.Position == 1 {
				ctx.MiddleAttribute = def.Name
			}
			if def.Name == "location" {
				ctx.LocationAttribute = def.Name
			}
		}
		if ctx.Validate() {
			return ctx
		}
		l.Error("context invalid")
	} else {
		l.Error("not enough attributes loaded")
	}
	os.Exit(1)
	return ContextValue{}
}

func (a *Attributes) AttributeExists(name AttributeName) bool {
	if _, exists := a.attribute_definitions[name]; exists {
		return true
	} else {
		l.Error("attribute \"" + name.String() + "\" does not exist")
	}
	return false
}

func (a *Attributes) AttributeKeyExists(name AttributeName, key AttributeKey) bool {
	if attr, exists := a.attribute_definitions[name]; exists {
		for _, value := range attr.Values {
			if value.ValueId == key {
				return true
			}
		}
		l.Error("attribute key \"" + key.String() + "\" in attribute \"" + name.String() + "\" does not exist")
	} else {
		l.Error("attribute " + name.String() + " does not exist")
	}
	return false
}

func (a *Attributes) AttributeGroupExists(name AttributeName, key AttributeGroupKey) bool {
	if attr, exists := a.attribute_definitions[name]; exists {
		for _, value := range attr.Groups {
			if value.GroupId == key {
				return true
			}
		}
		l.Error("attribute group " + key.String() + " in attribute " + name.String() + " does not exist")
	} else {
		l.Error("attribute " + name.String() + " does not exist")
	}
	return false
}

func (c *AttributeConfig) Validate() bool {
	res := true
	keys := AttributeKeys{}
	for _, v := range c.Values {
		if v.ValueId == "" {
			l.Error("attribute " + string(c.Name) + " has at least one key with no id")
			res = false
		}
		if slices.Contains(keys, v.ValueId) {
			l.Error("the key " + string(v.ValueId) + " in attribute " + string(c.Name) + " exists at least twice")
			res = false
		}
		keys = append(keys, v.ValueId)
	}
	groups := AttributeGroupKeys{}
	for _, g := range c.Groups {
		if g.GroupId == "" {
			l.Error("attribute " + c.Name.String() + " has at least one group with no id")
			res = false
		}
		if slices.Contains(groups, g.GroupId) {
			l.Error("the group " + g.GroupId.String() + " in attribute " + string(c.Name) + " exists at least twice")
			res = false
		}
		groups = append(groups, g.GroupId)
		for _, member := range g.Members {
			if !slices.Contains(keys, member.Key) {
				l.Error("in attribute \"" + c.Name.String() + "\" the group \"" + g.GroupId.String() + "\" contains not existing key \"" + member.Key.String() + "\"")
				res = false
			}
		}
	}
	return res
}

func (a *Attributes) Validate() bool {
	res := true
	for _, val := range a.attribute_definitions {
		res = res && val.Validate()
		// verify subgroup bindings
		for _, subgroup := range val.SubGroups {
			res = res && a.AttributeGroupExists(subgroup.SubGroupAttributeName, subgroup.SubGroupId)
			/*
				if !res {
					l.Error("attribute \"" + val.Name.String() + "\" references subgroup that does not exist: attr:\"" + subgroup.SubGroupAttributeName.String() + "\", group:\"" + subgroup.SubGroupId.String() + "\"")
				}
			*/
		}
		// verify value bindings
		for _, value := range val.Values {
			for _, subgroup := range value.SubGroups {
				res = res && a.AttributeGroupExists(subgroup.SubGroupAttributeName, subgroup.SubGroupId)
				if !res {
					l.Error("attribute \"" + val.Name.String() + "\" value " + value.ValueId.String() + " references not existing attribute key of attribute \"" + subgroup.SubGroupAttributeName.String() + "\": \"" + subgroup.SubGroupId.String() + "\"")
				}
			}
		}
	}
	return res
}
