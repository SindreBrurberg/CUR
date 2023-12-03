package config

import (
	"embed"
	_ "embed"
	"encoding/json"
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
)

//go:embed schema/*
var FS embed.FS

type rt struct {
	CurrentDirectory string `json:"currentDirectory"`
}

func LoadDirs() (cfgs []Root) {
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
			case "packages":
				fallthrough
			case "packageManagers":
				fallthrough
			case "features":
				fallthrough
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
	//var cfg root
	cfgs = make([]Root, len(systems))
	for i, system := range systems {
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
		if err := Load(wd, files, FS, &cfgs[i]); err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		//data, _ := yaml.Marshal(cfg)
		data, _ := json.MarshalIndent(cfgs[i], "", "    ")
		fmt.Printf("%s\n", data)
	}
	// This is a placeholder for anything that the program might actually do
	// with the configuration.
	return
}

func Load(dir string, files []string, fs iofs.FS, dest any) error {
	sbragi.Info("loading", "dir", dir, "files", files)
	overlay := make(map[string]load.Source)
	err := iofs.WalkDir(fs, ".", func(path string, d iofs.DirEntry, err error) error {
		sbragi.Info("reading fs", "path", path)
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
