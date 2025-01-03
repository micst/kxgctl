# kxgctl

> [!CAUTION]
>
> This project is not affiliated, associated, authorized, endorsed by,
> or in any way officially connected with the KNX Association or any
> subsidiaries or its affiliates.

[**K**]n[**X**] [**G**]roup address [**C**]on[**T**]ro[**L**], a command line tool for configuring KNX group addresses.

kxgctl provides built in syntax help. Run

- `kxgctl` to get a list of all available
commands and flags
- `kxgctl <command> --help` to get help for a command
- `kxgctl <command> <subcommand> -h` to get help for a subcommand
- ...

## Quickstart

Initialize a new workspace with demo files

```bash
kxgctl init
```

> [!TIP]
>
> **OPTIONAL: Edit Workspace**
>
> Edit `kxg_attributes_01_location.yaml` to suit your needs. This could mean
> 
> - add missing rooms to `values`
> - rename rooms if needed
> - assign correct `parent` values to the rooms (i.e. put rooms on correct level)
> - edit group `group-location-root` so it contains only the values you want
> 
> Edit `kxg_devices.yaml` to suit your needs. Make sure you only use *location* attribute
> keys that are in your location group.
>
> See section about YAML syntax for more information about
> how to edit configuration files.

Generate KNX group addresses

```bash
kxgctl generate
```

## Configuration overview

Some definitions:

- **configuration files** are `.yaml` files that are used when running kxgctl
- **workspace** is a directory where *configuration files* are stored
- **library** is a special *workspace* where common *configuration files* can be stored

The configuration structure can roughly be described as follows

- **Attributes** define address (main/middle) groups and address prefixes
- **Templates** bundle address definitions for reusability
- **Devices** configure template instances
- **Contexts** define which attributes and attribute groups should be used for which purpose

This leads to a configuration concept, that allows

 - reusing *configuration files* in several *workspace*s
 - simple adding of new *configuration files*
 - reduction of configuration code by providing reasonable defaults everywhere possible
 - keeping full flexibility by overriding almost any configuration value individually

For more information see the sections "Configuration YAML syntax" and "Binding concept".

## Workspace

> [!TIP]
>
> If you are new to kgxctl and have absolutely no clue how workspaces work, just
> run `kxgctl init` to generate a demo workspace and start editing.

kxgctl configuration filenames start with `kxg` and have the extension `.yaml`.

- **attribute files** start with `kxg_attributes`
- **template files** start with `kxg_templates`
- **device files** start with `kxg_devices`
- the **context file** is called `kxg.yaml`

Each **attribute file** defines exactly one attribute with it's values and groups.
If two attribues with the same name are defined within one workspace, kxgctl will fail.

Each **template file** may contain an arbitrary number of template definitions. Templates
can be collected in one file or split up in several files. If two templates with same
name are defined within one workspace, kxgctl will fail.

Each **device file** defines groups of template instances. Device files usually
bind templates to a location attribute, as a default location cannot reasonably
be defined in a template definition.

The **context file** contains a list of all available contexts. A context defines
which attributes and attribute groups are used for main group, middle groups
and as address prefix. If several contexts with the same context name exist,
kxgctl will fail.

See the **Configuration YAML syntax** section for detailed explanations how configuration
files are defined.

## Libraries

Per default, kxgctl tries to load default *configuration files* from a *library*.
The directory of the library is derived in the following order:

1. if `--skip-library` is set, no library will be loaded
2. if `--library` is used, the directory given here will be used
3. if `KXGCTL_LIBRARY` is set, the directory defined in the environment variable will be used
4. if all above is not set, the directory `~/.kxgctl` will be used

If the directory derived from the above rules does not exist, no library will be loaded.
Also, if the `--library` flag or the `KXGCTL_LIBRARY` environment variable are used and the
directory is not found, no fallback to other directories is done. No library will be loaded
in this case. 

If no library exists yet, you can initialize one with default files.
Run

```bash
kxgctl init library
```

The `--library` flag and `KXGCTL_LIBRARY` environment variable will apply
for `kxgctl init library` as well.

> [!TIP]
>
> **Overriding library items**
>
> While within a workspace **attributes**  and **templates** may only be defined once, items
> from a library may be overriden from configuration files in a workspace. So if a library
> contains some default *cover* template, another *cover* template in the workspace will
> replace the one from the library.
>
> Overrides like this will be shown in the log if `--verbose` is used.

## Configuration YAML syntax

### Attribute files

`<attribute>`:
```yaml
name: string            # mandatory, name of the attribute. Usually "location",
                        # "craft" or "function"
values:                 # mandatory, list of elements that can be formed to 
  - <attribute_value>   # groups in .groups
  - ...
groups:                 # mandatory, defines actual members of KNX group address
  - <attribute_group>   # main/middle groups
  - ...
middle_groups:             # optional, only required if this attribute is used as
  - <middle_group_binding> # main group attribute
  - ...
```
`<attribute_value>`:
```yaml
id: string              # mandatory, identifier of the value
text: string            # mandatory, text that is used for naming groups
                        # or prefix addresses
description: string     # optional, text used for description in ETS
parent: string          # optional, may reference another attribute_value.id
                        # of same attribute
root: bool              # optional, if true, text concatenation is skipped
                        # even if .parent is set
icon: string            # optional, reserved
ha_name: string         # optional, reserved
middle_groups:          # optional
  - <middle_group_binding> 
```
`<attribute_group>`:
```yaml
id: string         # mandatory, identifier of attribute group
members:           # mandatory, actual members (<attribute_value>) that
  - <group_member> # are grouped together
  - ...            #
                   # - when used as main group, up to 32 members are allowed
                   # - when used as middle group, up to 8 members are allowed
```
`<group_member>`:
```yaml
id: string                 # mandatory, identifier of <attribute_value>
middle_groups:             # optional, may used for redefining middle group
  - <middle_group_binding> # binding for this group member
  - ...                    # Most probably not required.
```
`<middle_group_binding>`:
```yaml
name: string # attribute name for middle group
id: string   # attribute_group.id from attribute .name to use as middle group
```

### Template files

```yaml
templates:
  - <template>
  - ...
```
`<template>`:
```yaml
name: string          # mandatory, identifier of the template
device_name: string   # optional, a default device name that is used if no
                      # device name is configured in template instantiation.
                      # In general, no device_name at all is required. It is
                      # used when generating the address name and left out if it
                      # does not exist.
bindings:
  - <address_binding> # optional, defines binding of all address values
  - ...               # in template to specific attribute values
addresses:
  - <address>         # mandatory, without any address, a template does not make
  - ...               # any sense.
```
`<address>`:
```yaml
id: string            # mandatory, identifier of the address within the template
                      # This is used when overriding address value members
                      # in a template instantiation in a device file.
name: string          # mandatory, name of the address, used for address
                      # name generation
description: string   # optional, text for description field of address
datatype: string      # optional, force data type for group address,
                      # see KNX datatype documentation for details 
security: string      # optional, set security configuration of generated address,
                      # may have values "", "On" and "Off"
central: bool         # optional, defaults to false
disabled: bool        # optional, defaults to false
bindings:
  - <address_binding> # optional, allows overriding <template>.bindings
  - ...               # for single addresses
```
`<address_binding>`:
```yaml
attribute_name: string # mandatory, attribute selector for the address
attribute_key: string  # mandatory, <attribute_value>.id of the
                       # attribute to bind to
```

### Device files

```yaml
groups:
  - <device_group>
  - ...
```
`<device_group>`:
```yaml
bindings:
  - <address_binding> # optional, define address bindings for a group of
  - ...               # template instances. Most commonly this is used to
                      # bind a group of template instances to "location"
                      # attribute, i.e. to collect all devices that
                      # belong to a common room or area.
devices:              # mandatory, template instances
  - <device>
  - ...
```
`<device>`:
```yaml
template_name: string # mandatory, the <template>.name to be instantiated
device_name: string   # optional, set/override device name, used for
                      # naming of address
start_address: string # optional, numerical address value for the first
                      # address of the device This can be used to force
                      # template addresses start at specific values.
                      # Usefull for leaving some addresses unused/reserved
                      # for future use.
bindings:
  - <address_binding> # optional, override  
  - ...
addresses:
  - <address>         # optional, override template address configuration
  - ...
```

### Context files

```yaml
contexts:                  # optional, list of contexts
  - <context>
  - 
current-context: string    # optional, name of the context that will be
                           # used by default
```
`<context>`:
```yaml
name: string               # mandatory, name of context
main_attribute: string     # mandatory, name of the attribute that will be
                           # used for main groups
middle_attribute: string   # mandatory, name of the attribute that will be
                           # used for middle groups
root_group: string         # mandatory, attribute_group.id that will be
                           # used as main group members
location_attribute: string # mandatory, name of the attribute that will be
                           # used as address name prefix
```

> [!NOTE]
>
> If no context is defined, a **default context** will be generated
> and named `default`. The members of the default context will be chosen
> as follows:
> - main attribute: first loaded attribute
> - middle attribute: second loaded attribute
> - root group: first group defined in first loaded attribute
> - location: attribute with name "location"

## Binding concept

There are two types of bindings in kxgctl config:

- middle_groups
- address_bindings 
