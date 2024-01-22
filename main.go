package main

import (
	"fmt"

	"github.com/mhashemm/torrent/bencode"
)

type Test struct {
	X *int `bencode:"x"`
}

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
	// _, err := upnp.DeletePortMapping(upnp.DeletePortMappingMsg{
	// 	NewExternalPort: 6969,
	// 	NewProtocol:     "TCP",
	// })
	x := 69
	bc, err := bencode.Marshal(Test{X: &x})
	fmt.Println(string(bc))
	fmt.Println(err)
}
