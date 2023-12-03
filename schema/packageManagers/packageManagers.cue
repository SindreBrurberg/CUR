package config

#PackageManagers: or([ for k, v in packageManagers {k}])
packageManagers?: {
	[string]: {
		syntax: string
		local?: string
		requires: [...#Features]
	}
}
