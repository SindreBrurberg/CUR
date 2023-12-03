package config

import "list"

#Uri: =~#"^[^:/\.\s][^:/\s]*\.[^:/\.\s]+$"#
#Url: =~#"^(https:\/\/)?([^:/\.\s][^:/\s]*\.)+[^:\/\.\s]+(\/[^:\/\s]+)*?$"#

name:        string
nerthus_url: #Uri
visuale_url: #Uri
systems:     [...#System] & list.MinItems(1)
