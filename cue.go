package config

import (
	"embed"
	_ "embed"
	"fmt"
	"io"
	iofs "io/fs"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/load"

	"github.com/cantara/bragi/sbragi"
	"gopkg.in/yaml.v2"
)

//go:embed schema/*
var fs embed.FS

type rt struct {
	CurrentDirectory string `json:"currentDirectory"`
}

func LoadDirs() {
	wd, err := os.Getwd()
	if err != nil {
		sbragi.WithError(err).Fatal("while getting wd")
	}
	files := []string{}
	systems := []string{}
	err = filepath.WalkDir(wd, func(path string, d iofs.DirEntry, err error) error {
		if path == wd {
			return nil
		}
		sbragi.Info("walking", "path", path, "base", strings.ToLower(filepath.Base(path)))
		if d.IsDir() {

			switch strings.ToLower(filepath.Base(path)) {
			case "services":
				fallthrough
			case "roles":
				//dirs = append(dirs, path)
				err = filepath.WalkDir(path, func(path string, d iofs.DirEntry, err error) error {
					if d.IsDir() {
						return nil
					}
					if filepath.Ext(path) != ".cue" {
						return nil
					}
					sbragi.Info("walking", "path", path, "base", strings.ToLower(filepath.Base(path)), "ext", filepath.Ext(path))
					files = append(files, path)
					return nil
				})
			case "files":
			default:
				systems = append(systems, path)
			}
			return filepath.SkipDir
		}
		if filepath.Ext(path) != ".cue" {
			return nil
		}
		sbragi.Info("walking", "path", path, "base", strings.ToLower(filepath.Base(path)), "ext", filepath.Ext(path))
		files = append(files, path)
		return nil
	})
	if err != nil {
		sbragi.WithError(err).Fatal("walked dir")
	}
	var cfg Environment
	for _, system := range systems {
		files := files
		err = filepath.WalkDir(system, func(path string, d iofs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}
			if filepath.Ext(path) != ".cue" {
				return nil
			}
			sbragi.Info("walking", "path", path, "base", strings.ToLower(filepath.Base(path)), "ext", filepath.Ext(path))
			files = append(files, path)
			return nil
		})
		if err != nil {
			sbragi.WithError(err).Fatal("walked system dir", "root", system)
		}
		//TODO: get all subdirs with the exception of the files dir
		if err := Load(wd, files, fs, &cfg); err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		data, _ := yaml.Marshal(cfg)
		fmt.Printf("%s\n", data)
	}
	// This is a placeholder for anything that the program might actually do
	// with the configuration.
}

type Environment struct {
	Name       string   `json:"name"`
	NerthusURL string   `json:"nerthus_url"`
	VisualeURL string   `json:"visuale_url"`
	Systems    []System `json:"systems"`
}
type Artifact struct {
	ID    string `json:"id"`
	Group string `json:"group"`
}
type Requirements struct {
	RAM              string `json:"ram"`
	Disk             string `json:"disk"`
	CPU              int    `json:"cpu"`
	PropertiesName   string `json:"properties_name"`
	WebserverPortKey string `json:"webserver_port_key"`
	NotClusterAble   bool   `json:"not_cluster_able"`
	IsFrontend       bool   `json:"is_frontend"`
	Roles            []any  `json:"roles"`
	Services         []any  `json:"services"`
}
type ServiceInfo struct {
	Name         string       `json:"name"`
	ServiceType  string       `json:"service_type"`
	HealthType   string       `json:"health_type"`
	APIPath      string       `json:"api_path"`
	Artifact     Artifact     `json:"artifact"`
	Requirements Requirements `json:"requirements"`
}
type Service struct {
	Name       string      `json:"name"`
	Definition ServiceInfo `json:"definition"`
}
type Cluster struct {
	Name     string    `json:"name"`
	Node     Node      `json:"node"`
	Services []Service `json:"services"`
	Internal bool      `json:"internal"`
}
type Node struct {
	Os   string `json:"os"`
	Arch string `json:"arch"`
	Size string `json:"size"`
}
type System struct {
	Name          string    `json:"name"`
	Domain        string    `json:"domain"`
	RoutingMethod string    `json:"routing_method"`
	Cidr          string    `json:"cidr"`
	Zone          string    `json:"zone"`
	Clusters      []Cluster `json:"clusters"`
}

func Load(dir string, files []string, fs iofs.FS, dest any) error {
	sbragi.Info("loading", "dir", dir, "files", files)
	overlay := make(map[string]load.Source)
	err := iofs.WalkDir(fs, ".", func(path string, d iofs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".cue" {
			return nil
		}
		file, err := fs.Open(path)
		if err != nil {
			return err
		}
		bytes, err := io.ReadAll(file)
		if err != nil {
			return err
		}
		path = filepath.Join(dir, path)
		overlay[path] = load.FromBytes(bytes)
		files = append(files, path)
		return nil
	})
	if err != nil {
		return err
	}
	configInst := load.Instances(files, &load.Config{
		Dir:     dir,
		Package: "*",
		Overlay: overlay,
	})[0]

	sbragi.Info("loaded instances")
	if err := configInst.Err; err != nil {
		return fmt.Errorf("cannot load configuration from %q: %v", configInst.Root, err)
	}
	ctx := cuecontext.New()
	configVal := ctx.BuildInstance(configInst)
	fields, err := configVal.Fields()
	sbragi.WithError(err).Trace("built instance and got fields")
	//a, d := configVal.Struct()
	//i := a.Fields()
	for fields.Next() {
		fmt.Println(fields.Label(), fields.Value())
	}
	if err := configVal.Validate(cue.All()); err != nil {
		return fmt.Errorf("invalid configuration from %q: %v", dir, errors.Details(err, nil))
	}

	if err := configVal.Decode(dest); err != nil {
		return fmt.Errorf("cannot decode final configuration: %v", errors.Details(err, nil))
	}
	return nil
}
