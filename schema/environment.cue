package config

import "list"

#Url: =~#"^[^:/\.\s][^:/\s]*\.[^:/\.\s]+$"#

name:        string
nerthus_url: #Url
visuale_url: #Url
systems:     [...#System] & list.MinItems(1)
