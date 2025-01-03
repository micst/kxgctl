package kxg

import (
	"embed"
	"io"
	"os"
	"path/filepath"
	"strings"

	l "github.com/micst/kxgctl/kxg/logging"
	"github.com/micst/kxgctl/kxg/yaml"
)

//go:generate pwsh -NoProfile -Command "Copy-Item -Path ../examples/ex1/* -Recurse -Destination ./example"
//go:embed example/*.yaml
var InitData embed.FS

const (
	DefaultPrefix          = "kxg"
	DefaultSeparator       = "_"
	DefaultInfixAttributes = "attributes"
	DefaultInfixTemplates  = "templates"
	DefaultInfixDevices    = "devices"
	DefaultContextFileName = DefaultPrefix + ".yaml"
)

type WorkspaceFiles struct {
	Attributes []string
	Templates  []string
	Devices    []string
	Contexts   string
}

type Workspace struct {
	Directory string
	Exists    bool
	Files     WorkspaceFiles
}

func (w *Workspace) prepare() {
	if _, err := os.Stat(w.Directory); err == nil {
		w.Exists = true
		w.Files.Attributes = GetConfigFiles(w.Directory, DefaultInfixAttributes)
		w.Files.Templates = GetConfigFiles(w.Directory, DefaultInfixTemplates)
		w.Files.Devices = GetConfigFiles(w.Directory, DefaultInfixDevices)
		if _, err := os.Stat(filepath.Join(w.Directory, DefaultContextFileName)); err == nil {
			w.Files.Contexts = DefaultContextFileName
		}
	} else {
		l.Debug("no workspace found in \"" + w.Directory + "\"")
		w.Exists = false
	}
}

func (w *Workspace) Load(as_lib bool) {
	w.prepare()
	if w.Exists {
		for _, file_name := range w.Files.Attributes {
			file_path := filepath.Join(w.Directory, file_name)
			Data.Attributes.LoadYaml(file_path, as_lib)
		}
		for _, file_name := range w.Files.Templates {
			file_path := filepath.Join(w.Directory, file_name)
			Data.Templates.LoadYaml(file_path, as_lib)
		}
		if !as_lib {
			for _, file_name := range w.Files.Devices {
				file_path := filepath.Join(w.Directory, file_name)
				Data.Devices.LoadYaml(file_path)
			}
		}
		if w.Files.Contexts != "" {
			file_path := filepath.Join(w.Directory, w.Files.Contexts)
			Data.Contexts.LoadYaml(file_path)
		}
		// make sure default context is always appended
		Data.Contexts.Contexts = append(
			Data.Contexts.Contexts,
			Data.Attributes.DefaultContext(),
		)
		// if no current context set, choose the 1st one
		if Data.Contexts.CurrentContext == "" {
			Data.Contexts.CurrentContext = Data.Contexts.Contexts[0].Name
		}
		if Args.ContextName != "" {
			Data.Contexts.CurrentContext = yaml.ContextName(Args.ContextName)
		}
	} else {
		if !as_lib {
			l.Error("no workspace found in \"" + w.Directory + "\", could not load")
		}
	}
	l.Debug3("loading workspace finished from \"" + w.Directory + "\"")
}

func (w *Workspace) ResetDirectory(force bool) error {
	dir := w.Directory
	if fileInfo, err := os.Stat(dir); err == nil {
		if !fileInfo.IsDir() {
			l.Error("workspace at \"" + dir + "\" exists but is not a directory")
			return os.ErrExist
		} else {
			files := GetConfigFiles(dir, "")
			if len(files) > 0 {
				l.Info("resetting existing workspace at \"" + dir + "\"")
				l.Debug("removing files:")
				for _, file_name := range files {
					if force {
						l.Debug(" - " + file_name)
						os.Remove(file_name)
					} else {
						l.Error("found configuration files when resetting workspace directory")
						l.Error("rerun with --force to delete existing files")
						return os.ErrExist
					}
				}
			} else {
				l.Debug("workspace in \"" + dir + "\" already empty")
			}
		}
	} else {
		l.Debug("creating new workspace directory at \"" + dir + "\"")
		os.Mkdir(dir, os.ModeDir)
	}
	return nil
}

func (w *Workspace) CopyFromResources(force bool) error {
	if err := w.ResetDirectory(force); err != nil {
		return err
	}
	dir := "example"
	if d, err := InitData.ReadDir(dir); err == nil {
		l.Debug("creating files:")
		for _, entry := range d {
			file := dir + "/" + entry.Name()
			if content, err := InitData.ReadFile(file); err == nil {
				ofile := filepath.Join(w.Directory, entry.Name())
				l.Debug(" - " + ofile)
				if f, err := os.Create(ofile); err == nil {
					defer f.Close()
					f.Write(content)
				} else {
					l.Error("could not create file " + ofile)
				}
			} else {
				l.Error("could not read init file " + file)
			}
		}
	} else {
		l.Error("could not list init directory " + dir)
	}
	return nil
}

func (w *Workspace) CopyFromDirectory(directory string, force bool) error {
	src := Workspace{
		Directory: directory,
	}
	src.prepare()
	if !src.Exists {
		l.Error("no workspace/library at src directory found: \"" + directory + "\"")
	}
	if err := w.ResetDirectory(force); err != nil {
		return err
	}
	files := src.Files.Attributes
	files = append(files, src.Files.Templates...)
	files = append(files, src.Files.Devices...)
	if src.Files.Contexts != "" {
		files = append(files, src.Files.Contexts)
	}
	l.Debug("copying workspace files to \"" + w.Directory + "\"")
	for _, file_name := range files {
		file_path_in := filepath.Join(src.Directory, file_name)
		if src, src_err := os.Open(file_path_in); src_err == nil {
			defer src.Close()
			file_path_out := filepath.Join(w.Directory, file_name)
			if dst, dst_err := os.Create(file_path_out); dst_err == nil {
				defer dst.Close()
				if _, write_err := io.Copy(dst, src); write_err != nil {
					l.Error(" - could not copy data to \"" + dst.Name())
					return os.ErrInvalid
				}
			}
		}
	}
	return nil
}

func GetConfigFiles(dir string, infix string) []string {
	res := []string{}
	if files, err := os.ReadDir(dir); err == nil {
		prefix := DefaultPrefix
		if infix != "" {
			prefix += DefaultSeparator + infix
		}
		for _, file := range files {
			file_name := file.Name()
			if strings.HasSuffix(file_name, ".yaml") {
				if strings.HasPrefix(file_name, prefix) {
					res = append(res, file_name)
				}
			}
		}
	}
	return res
}
