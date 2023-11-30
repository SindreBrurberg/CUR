package config

#nerthusSD: #ServiceDefinition & {
	name:         "nerthus"
	service_type: "H2A"
	health_type:  "go"
	api_path:     "/health"
	artifact: {
		id:    "nerthus2"
		group: "no/cantara/gotools"
	}
	requirements: {
		ram:                "2GB"
		disk:               "30GB"
		cpu:                2
		properties_name:    ".env"
		webserver_port_key: "webserver.port"
		not_cluster_able:   true
		is_frontend:        true
		roles: []
		services: []
	}
}
