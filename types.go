package config

type Root struct {
	OS              map[string]OS             `json:"os,omitempty"`
	Features        map[string]Feature        `json:"features,omitempty"`
	PackageManagers map[string]PackageManager `json:"packageManagers,omitempty"`
	Packages        map[string]Package        `json:"packages,omitempty"`
	Name            string                    `json:"name"`
	NerthusURL      string                    `json:"nerthus_url"`
	VisualeURL      string                    `json:"visuale_url"`
	Systems         []System                  `json:"systems"`
}
type Artifact struct {
	ID    string `json:"id"`
	Group string `json:"group"`
}
type Requirements struct {
	RAM              string   `json:"ram"`
	Disk             string   `json:"disk"`
	CPU              int      `json:"cpu"`
	PropertiesName   string   `json:"properties_name"`
	WebserverPortKey string   `json:"webserver_port_key"`
	NotClusterAble   bool     `json:"not_cluster_able"`
	IsFrontend       bool     `json:"is_frontend"`
	Features         []string `json:"features"`
	Packages         []string `json:"packages"`
	Services         []string `json:"services"`
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
type Task struct {
	Info    string `json:"info,omitempty"`
	Type    string `json:"type,omitempty"`
	Source  string `json:"source,omitempty"`
	Dest    string `json:"dest,omitempty"`
	File    string `json:"file,omitempty"`
	Url     string `json:"url,omitempty"`
	Manager string `json:"manager,omitempty"`
	Package string `json:"package,omitempty"`
	Service string `json:"service,omitempty"`
	Start   bool   `json:"start,omitempty"`
}
type Feature struct {
	Friendly string   `json:"friendly,omitempty"`
	Requires []string `json:"requires,omitempty"`
	Tasks    []Task   `json:"tasks,omitempty"`
	Custom   map[string][]Task
}
type PackageManager struct {
	Syntax string `json:"syntax,omitempty"`
	Local  string `json:"local,omitempty"`
}
type Package struct {
	Managers []string `json:"managers,omitempty"`
	Provides []string `json:"provides,omitempty"`
}
type OS struct {
	PackageManagers []string `json:"packageManagers,omitempty"`
	Provides        []string `json:"provides,omitempty"`
}
