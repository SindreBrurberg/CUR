package main

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
type Feature struct {
	Name     string `json:"name,omitempty"`
	Friendly string `json:"friendly,omitempty"`
	Tasks    []Task `json:"tasks,omitempty"`
}
type Task struct {
	Info    string         `json:"info,omitempty"`
	Type    string         `json:"type,omitempty"`
	Source  string         `json:"source,omitempty"`
	Dest    string         `json:"dest,omitempty"`
	File    string         `json:"file,omitempty"`
	Url     string         `json:"url,omitempty"`
	Manager PackageManager `json:"manager,omitempty"`
	Package Package        `json:"package,omitempty"`
	Service string         `json:"service,omitempty"`
	Start   bool           `json:"start,omitempty"`
}
type Requirements struct {
	RAM              string    `json:"ram"`
	Disk             string    `json:"disk"`
	CPU              int       `json:"cpu"`
	PropertiesName   string    `json:"properties_name"`
	WebserverPortKey string    `json:"webserver_port_key"`
	NotClusterAble   bool      `json:"not_cluster_able"`
	IsFrontend       bool      `json:"is_frontend"`
	Features         []Feature `json:"features"`
	Packages         []Package `json:"packages"`
	Services         []string  `json:"services"`
}
type Package struct {
	Name     string   `json:"name"`
	Managers []string `json:"managers,omitempty"`
	Provides []string `json:"provides,omitempty"`
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
	Os   OS     `json:"os"`
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
type PackageManager struct {
	Name   string `json:"name,omitempty"`
	Syntax string `json:"syntax,omitempty"`
	Local  string `json:"local,omitempty"`
}
type OS struct {
	PackageManagers []PackageManager `json:"package_managers,omitempty"`
	Provides        []string         `json:"provides,omitempty"`
}
