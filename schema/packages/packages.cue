package config

#Packages: or([ for k, v in packages {k}])
packages?: {
	[string]: {
		managers: [...string]
		provides: [...string]
	}
}
