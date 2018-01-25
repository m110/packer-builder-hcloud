package main

import (
	"log"

	"github.com/hashicorp/packer/packer/plugin"
	"github.com/m110/packer-builder-hcloud/builder/hcloud"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		log.Fatal(err)
	}

	server.RegisterBuilder(new(hcloud.Builder))
	server.Serve()
}
