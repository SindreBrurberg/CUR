package config

features: zulu: {
	tasks: [
		#InstallExternal & {
			info:    "Install zulu repo"
			url:     "https://cdn.azul.com/zulu/bin/zulu-repo-1.0.0-1.noarch.rpm"
			manager: "yum"
		},
	]
}
