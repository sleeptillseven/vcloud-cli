# a cli for the vcloud api

## usage

set the following env vars
* VCD_USER
* VCD_PASSWORD
* VCD_ORG

explore the possiblities of the cli by using the help.  

the command structure of the vcloud-cli:  
`vcloud-cli --network query allocatedips`

`query` --> root command  
`allocatedips` --> sub command  
`--network` --> argument for the last command

at every level you can use the help:    
* `vcloud-cli query --help`
* `vcloud-cli query allocatedips --help`
