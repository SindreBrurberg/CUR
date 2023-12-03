package config

features: python: {
	tasks: [
		#Install & {
			package: "python3"
		},
		#Install & {
			package: "python3-pip"
		},
	]
	custom: {
		_#ChromeRPM: "/tmp/google-chrome-stable_current_x86_64.rpm"
		"Amazon Linux 2023": [
			#Download & {
				info:   "Download package"
				source: "https://dl.google.com/linux/direct/google-chrome-stable_current_x86_64.rpm"
				dest:   _#ChromeRPM
			},
			#InstallLocal & {
				info:    "Import and install package"
				file:    _#ChromeRPM
				manager: "yum"
			},
			#Delete & {
				info: "Delete local package file"
				file: _#ChromeRPM
			},
		]
	}
}
