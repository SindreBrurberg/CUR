package config

#Service: {
	name:            string
	git?:            string
	branch?:         string
	webserver_port?: int
	properties?:     string
	dirs?: {[string]: string}
	files?: {[string]: string}
	definition: #ServiceDefinition
}

#ServiceDefinition: {
	name:         string
	service_type: string
	health_type:  string
	api_path:     string
	artifact:     #Artifact
	requirements: #Requirements
}

#Requirements: {
	ram:                =~#"\d+[TGMK]B"#
	disk:               =~#"\d+[TGMK]B"#
	cpu:                int
	properties_name:    string
	webserver_port_key: string
	not_cluster_able:   bool
	is_frontend:        bool
	features: [...#Features]
	packages: [...#Packages]
	services: [...string]
}

#Artifact: {
	id:             string
	group:          string
	release_repo?:  string
	snapshot_repo?: string
	user?:          string
	password?:      string
}
