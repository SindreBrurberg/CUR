package config

#visualeSD: {
	name:         "visuale"
	service_type: "H2A"
	health_type:  "go"
	api_path:     "/health"
	artifact: {
		id:    "visuale"
		group: "no/cantara"
	}
	requirements: {
		ram:                "2GB"
		disk:               "30GB"
		cpu:                2
		properties_name:    "local_override.propperties"
		webserver_port_key: "webserver.port"
		not_cluster_able:   true
		is_frontend:        true
		features: [
			"zulu",
			"java17",
			"cron",
		]
		packages: []
		services: ["nerthus"]
	}
}
