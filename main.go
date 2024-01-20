package main

import (
	"fmt"

	"github.com/mhashemm/torrent/upnp"
)

func main() {
	// _, err := upnp.AddPortMapping(upnp.AddPortMappingMsg{
	// 	NewProtocol:               "TCP",
	// 	NewRemoteHost:             struct{}{},
	// 	NewExternalPort:           6969,
	// 	NewInternalPort:           6969,
	// 	NewEnabled:                1,
	// 	NewPortMappingDescription: "hey man",
	// 	NewLeaseDuration:          0,
	// })
	_, err := upnp.DeletePortMapping(upnp.DeletePortMappingMsg{
		NewExternalPort: 6969,
		NewProtocol:     "TCP",
	})
	fmt.Println(err)
}
