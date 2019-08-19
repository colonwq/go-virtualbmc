A go based virtbmc client

This is designed to be API compatable with the vbmc client.
The command prints the JSON string returned from the server

Requirements:
 - developed on Fedora 30
 - uses golang-github-pebbe-zmq4 for messaging
 - developed against python3-virtualbmc-1.5. Should work with 1.4

Usage:
List configured IPMI interfaces
   go run go-vbmc.go list

Add a IPMI interface
   go run go-vbmc.go add VM-name

Show an IPMI interface
   go run go-vbmc.go show VM-name

Start an IPMI interface
   go run go-vbmc.go start VM-name

Stop an IPMI interface
   go run go-vbmc.go stop VM-name

Delete an IPMI interface
   go run go-vbmc.go Delete VM-name

The list, add and show sub-commands support the -h parameter for additonal options

TODO:
- add support for pretty printing the JSON reply string
- add support for formatting options from the list and show sub-commands
- set the return code the rc value in the JSON string

Example usage:

$ go run go-vbmc.go list
json reply( 55 ):  [{"rc": 0, "header": ["Domain name", "Status", "Address", "Port"], "rows": []}]

$ go run go-vbmc.go add -port 6230 fedora30-server
json reply( 202 ):  [{"rc": 0, "msg": []}] .

$ go run go-vbmc.go show fedora30-server
json reply( 70 ):  [{"rc": 0, "header": ["Property", "Value"], "rows": [["username", "admin"], ["password", "password"], ["address", "::"], ["port", 6230], ["domain_name", "fedora30-server"], ["libvirt_uri", "qemu:///system"], ["libvirt_sasl_username", ""], ["libvirt_sasl_password", ""], ["active", "False"], ["status", "down"]]}]

$ go run go-vbmc.go start fedora30-server
Retrun string:  [{"rc": 0, "msg": []}]

$ go run go-vbmc.go show fedora30-server
json reply( 70 ):  [{"rc": 0, "header": ["Property", "Value"], "rows": [["username", "admin"], ["password", "password"], ["address", "::"], ["port", 6230], ["domain_name", "fedora30-server"], ["libvirt_uri", "qemu:///system"], ["libvirt_sasl_username", ""], ["libvirt_sasl_password", ""], ["active", "True"], ["status", "running"]]}]

$ go run go-vbmc.go stop fedora30-server
Retrun string:  [{"rc": 0, "msg": []}]

$ go run go-vbmc.go show fedora30-server
json reply( 70 ):  [{"rc": 0, "header": ["Property", "Value"], "rows": [["username", "admin"], ["password", "password"], ["address", "::"], ["port", 6230], ["domain_name", "fedora30-server"], ["libvirt_uri", "qemu:///system"], ["libvirt_sasl_username", ""], ["libvirt_sasl_password", ""], ["active", "False"], ["status", "down"]]}]

$ go run go-vbmc.go list
json reply( 55 ):  [{"rc": 0, "header": ["Domain name", "Status", "Address", "Port"], "rows": [["fedora30-server", "down", "::", 6230]]}]

$ go run go-vbmc.go delete fedora30-server
Retrun string:  [{"rc": 0, "msg": []}]

$ go run go-vbmc.go list
json reply( 55 ):  [{"rc": 0, "header": ["Domain name", "Status", "Address", "Port"], "rows": []}]
