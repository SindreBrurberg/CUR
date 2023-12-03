package main

import (
	"encoding/json"
	"fmt"
	iofs "io/fs"
	"os"
	"path/filepath"
	"strings"

	config "github.com/SindreBrurberg/CUR"
	"github.com/cantara/bragi/sbragi"
)

func main() {
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
	cfgs := make([]config.Root, len(systems))
system:
	for num, system := range systems {
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
		if err := config.Load(wd, files, config.FS, &cfgs[num]); err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		//data, _ := yaml.Marshal(cfg)
		sys := make([]System, len(cfgs[num].Systems))
		for i, s := range cfgs[num].Systems {
			sys[i].Name = s.Name
			sys[i].Cidr = s.Cidr
			sys[i].Zone = s.Zone
			sys[i].Domain = s.Domain
			sys[i].RoutingMethod = s.RoutingMethod
			sys[i].Clusters = make([]Cluster, len(s.Clusters))
			for j, c := range s.Clusters {
				sys[i].Clusters[j].Name = c.Name
				os, ok := cfgs[num].OS[c.Node.Os]
				if !ok {
					sbragi.Error("os is not present", "name", c.Node.Os)
					continue system
				}
				sbragi.Info("cluste", "os", os, "oses", cfgs[num].OS)
				pms := make([]PackageManager, len(os.PackageManagers))
				for i, name := range os.PackageManagers {
					pm, ok := cfgs[num].PackageManagers[name]
					if !ok {
						sbragi.Error("os is not present", "name", name)
						continue system
					}
					pms[i] = PackageManager{
						Name:   name,
						Syntax: pm.Syntax,
						Local:  pm.Local,
					}
				}
				sys[i].Clusters[j].Node = Node{
					Os: OS{
						PackageManagers: pms,
						Provides:        os.Provides,
					},
					Arch: c.Node.Arch,
					Size: c.Node.Size,
				}
				sbragi.Info("cluste", "node", sys[i].Clusters[j].Node)
				sys[i].Clusters[j].Internal = c.Internal
				sys[i].Clusters[j].Services = make([]Service, len(c.Services))
				for k, serv := range c.Services {
				servReq:
					for _, name := range serv.Definition.Requirements.Services {
						nameLow := strings.ToLower(name)
						for _, clust := range s.Clusters {
							for _, serv := range clust.Services {
								if nameLow == strings.ToLower(serv.Name) {
									continue servReq
								}
							}
						}
						sbragi.Error("missing service requirement in system", "name", name)
					}
					f := make([]Feature, len(serv.Definition.Requirements.Features))
					required := 0
					for i, name := range serv.Definition.Requirements.Features {
						featDef, ok := cfgs[num].Features[name]
						if !ok {
							sbragi.Error("feature is not present", "name", name)
							continue system
						}
						for _, reqName := range featDef.Requires { //This strat will add duplicates, Shouls change from list to a ordered graph or something similar
							def, ok := cfgs[num].Features[reqName]
							if !ok {
								sbragi.Error("feature is not present", "name", reqName)
								continue system
							}
							var tasks []Task
							if cust, ok := def.Custom[c.Node.Os]; ok {
								tasks = make([]Task, len(cust))
								for tn, task := range cust {
									tasks[tn] = confTaskToTask(task, cfgs[num])
								}
							} else {
								tasks = make([]Task, len(def.Tasks))
								for tn, task := range def.Tasks {
									tasks[tn] = confTaskToTask(task, cfgs[num])
								}
							}
							f[i+required] = Feature{
								Name:     reqName,
								Friendly: def.Friendly,
								Tasks:    tasks,
							}
							required++
							f = append(f, Feature{})
						}
						var tasks []Task
						if cust, ok := featDef.Custom[c.Node.Os]; ok {
							tasks = make([]Task, len(cust))
							for tn, task := range cust {
								tasks[tn] = confTaskToTask(task, cfgs[num])
							}
						} else {
							tasks = make([]Task, len(featDef.Tasks))
							for tn, task := range featDef.Tasks {
								tasks[tn] = confTaskToTask(task, cfgs[num])
							}
						}
						feat := Feature{
							Name:     name,
							Friendly: featDef.Friendly,
							Tasks:    tasks,
						}
						f[i+required] = feat
					}
					p := make([]Package, len(serv.Definition.Requirements.Packages))
					for i, name := range serv.Definition.Requirements.Packages {
						def, ok := cfgs[num].Packages[name]
						if !ok {
							sbragi.Error("package is not present", "name", name)
							continue system
						}
						p[i] = Package{
							Name:     name,
							Managers: def.Managers,
							Provides: def.Provides,
						}
					}
					sys[i].Clusters[j].Services[k].Name = serv.Name
					sys[i].Clusters[j].Services[k].Definition = ServiceInfo{
						Name:        serv.Definition.Name,
						ServiceType: serv.Definition.ServiceType,
						HealthType:  serv.Definition.HealthType,
						APIPath:     serv.Definition.APIPath,
						Artifact:    Artifact(serv.Definition.Artifact),
						Requirements: Requirements{
							RAM:              serv.Definition.Requirements.RAM,
							Disk:             serv.Definition.Requirements.Disk,
							CPU:              serv.Definition.Requirements.CPU,
							PropertiesName:   serv.Definition.Requirements.PropertiesName,
							WebserverPortKey: serv.Definition.Requirements.WebserverPortKey,
							NotClusterAble:   serv.Definition.Requirements.NotClusterAble,
							IsFrontend:       serv.Definition.Requirements.IsFrontend,
							Features:         f,
							Packages:         p,
							Services:         serv.Definition.Requirements.Services,
						},
					}
				}
			}
		}
		cfg := Environment{
			Name:       cfgs[num].Name,
			NerthusURL: cfgs[num].NerthusURL,
			VisualeURL: cfgs[num].VisualeURL,
			Systems:    sys,
		}
		data, _ := json.MarshalIndent(cfg, "", "    ")
		fmt.Printf("%s\n", data)
		for _, s := range cfg.Systems {
			for _, c := range s.Clusters {
				fmt.Println("system script:")
				for _, serv := range c.Services {
					fmt.Println("service script:")
					for _, pack := range serv.Definition.Requirements.Packages {
						for _, manager := range c.Node.Os.PackageManagers {
							if contains(pack.Managers, manager.Name) < 0 {
								continue
							}
							fmt.Println(strings.ReplaceAll(manager.Syntax, "<package>", pack.Name))
							break
						}
					}
					for _, feat := range serv.Definition.Requirements.Features {
						for _, task := range feat.Tasks {
							if task.Info != "" {
								fmt.Printf("//%s (%s)\n", task.Info, task.Type)
							}
							switch task.Type {
							case "install":
								if task.Manager.Name != "" {
									fmt.Println(strings.ReplaceAll(task.Manager.Syntax, "<package>", task.Package.Name))
									continue
								}
								for _, manager := range c.Node.Os.PackageManagers {
									if contains(task.Package.Managers, manager.Name) < 0 {
										continue
									}
									fmt.Println(strings.ReplaceAll(manager.Syntax, "<package>", task.Package.Name))
									break
								}
							case "install_external":
								fmt.Println(strings.ReplaceAll(task.Manager.Syntax, "<package>", task.Url))
							case "install_local":
								if task.Manager.Name != "" {
									fmt.Println(strings.ReplaceAll(task.Manager.Local, "<file>", task.File))
									continue
								}
								for _, manager := range c.Node.Os.PackageManagers {
									if contains(task.Package.Managers, manager.Name) < 0 {
										continue
									}
									fmt.Println(strings.ReplaceAll(manager.Local, "<file>", task.File))
									break
								}
							case "download":
								fmt.Printf("curl \"%s\" > \"%s\"\n", task.Source, task.Dest)
							case "link":
								fmt.Printf("ln -s \"%s\" \"%s\"\n", task.Dest, task.Source)
							case "delete":
								fmt.Printf("rm \"%s\"\n", task.File)
							case "enable":
								fmt.Printf("systemctl enable \"%s\"", task.Service)
								if task.Start {
									fmt.Println(" --now")
								}
							case "schedule":
							default:
								sbragi.Info("task not suported", "type", task.Type)
							}
						}
					}
				}
			}
		}
	}
	// This is a placeholder for anything that the program might actually do
	// with the configuration.
	fmt.Println("Configs:")
	for _, cfg := range cfgs {
		fmt.Println("    Packages:")
		for k := range cfg.Packages {
			fmt.Println("        ", k)
		}
		fmt.Println("    Features:")
		for k := range cfg.Features {
			fmt.Println("        ", k)
		}
	}
}

func contains[T comparable](arr []T, v T) int {
	for i, el := range arr {
		if el != v {
			continue
		}
		return i
	}
	return -1
}

func confTaskToTask(task config.Task, cfg config.Root) Task {
	switch task.Type {
	case "install":
		var m PackageManager
		if pm, ok := cfg.PackageManagers[task.Manager]; ok {
			m = PackageManager{
				Name:   task.Manager,
				Syntax: pm.Syntax,
				Local:  pm.Local,
			}
		}
		def, ok := cfg.Packages[task.Package]
		if !ok {
			sbragi.Error("package is not present", "name", task.Package)
			//continue system
			return Task{}
		}
		sbragi.Info("install", "managers", def.Managers)
		return Task{
			Info:    task.Info,
			Type:    task.Type,
			Manager: m,
			Package: Package{
				Name:     task.Package,
				Managers: def.Managers,
				Provides: def.Provides,
			},
		}
	case "install_local":
		var m PackageManager
		if pm, ok := cfg.PackageManagers[task.Manager]; ok {
			m = PackageManager{
				Name:   task.Manager,
				Syntax: pm.Syntax,
				Local:  pm.Local,
			}
		}
		return Task{
			Info:    task.Info,
			Type:    task.Type,
			Manager: m,
			File:    task.File,
		}
	case "install_external":
		var m PackageManager
		if pm, ok := cfg.PackageManagers[task.Manager]; ok {
			m = PackageManager{
				Name:   task.Manager,
				Syntax: pm.Syntax,
				Local:  pm.Local,
			}
		}
		return Task{
			Info:    task.Info,
			Type:    task.Type,
			Manager: m,
			Url:     task.Url,
		}
	case "download":
		return Task{
			Info:   task.Info,
			Type:   task.Type,
			Source: task.Source,
			Dest:   task.Dest,
		}
	case "link":
		return Task{
			Info:   task.Info,
			Type:   task.Type,
			Source: task.Source,
			Dest:   task.Dest,
		}
	case "delete":
		return Task{
			Info: task.Info,
			Type: task.Type,
			File: task.File,
		}
	case "enable":
		return Task{
			Info:    task.Info,
			Type:    task.Type,
			Service: task.Service,
			Start:   task.Start,
		}
	case "schedule":
		return Task{
			Info: task.Info,
			Type: task.Type,
		}
		/*
			Type    string         `json:"type,omitempty"`
			Source  string         `json:"source,omitempty"`
			Dest    string         `json:"dest,omitempty"`
			File    string         `json:"file,omitempty"`
			Manager PackageManager `json:"manager,omitempty"`
			Package Package        `json:"package,omitempty"`
			Service string         `json:"service,omitempty"`
		*/
	default:
		sbragi.Warning("unsuported task", "type", task.Type)
	}
	return Task{}
}
