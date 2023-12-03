package config

#Features: or([ for k, v in features {k}])
features?: {
	[string]: {
		info?: string
		tasks?: [...#Tasks]
		custom?: [#OS]: [...#Tasks]
		requires: [...#Features]
	}
}

#Tasks: #Install | #InstallLocal | #InstallExternal | #Enable | #Download | #Link | #Delete | #Schedule

#Task: {
	info?: string
	type:  string
}

#Install: {
	#Task
	type:    "install"
	package: string
}

#InstallLocal: {
	#Task
	type:    "install_local"
	file:    string
	manager: string
}

#InstallExternal: {
	#Task
	type:    "install_external"
	url:     #Url
	manager: string
}

#Enable: {
	#Task
	type:    "enable"
	service: string
	start:   bool | *true
}

#Download: {
	#Task
	type:   "download"
	source: #Url
	dest:   string
}

#Link: {
	#Task
	type:   "link"
	source: string
	dest:   string
}

#Delete: {
	#Task
	type: "delete"
	file: string
}

#Schedule: {
	#Task
	type:     "schedule"
	cronTime: =~#"^(\*\/)?([1-5]?[0-9]|\*) (\*\/)?(2[0-3]|1[0-9]|[0-9]|\*) (\*\/)?(3[01]|[12][0-9]|[1-9]|\*) (\*\/)?(1[0-2]|[1-9]|\*) (\*\/)?([0-6]|\*)$"#
	command:  string
}
