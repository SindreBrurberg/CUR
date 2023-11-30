package config

systems: [#Oplysningen]

#Oplysningen: #System & {
	name:           "oplysningen"
	domain:         "exoreaction.dev"
	routing_method: "host"
	cidr:           "10.0.0.0/24"
	zone:           "catalyst_one.infra"
	clusters: [#nerthus, #visuale]
}

#nerthus: #Cluster & {
	name: "nerthus"
	node: {
		size: "t4g.medium"
	}
	services: [{
		name:       "nerthus"
		definition: #ServiceDefinition & #nerthusSD
	}]
}

#visuale: #Cluster & {
	name: "visuale"
	node: {
		size: "t3.medium"
	}
	services: [{
		name:       "visuale"
		definition: #ServiceDefinition & #visualeSD
	}]
}
