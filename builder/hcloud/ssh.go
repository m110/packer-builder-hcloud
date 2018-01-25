package hcloud

import (
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/mitchellh/multistep"
	"golang.org/x/crypto/ssh"
)

func commHost(state multistep.StateBag) (string, error) {
	serverData := state.Get("server_data").(hcloud.ServerCreateResult)
	return serverData.Server.PublicNet.IPv4.IP.String(), nil
}

func sshConfig(state multistep.StateBag) (*ssh.ClientConfig, error) {
	serverData := state.Get("server_data").(hcloud.ServerCreateResult)

	return &ssh.ClientConfig{
		// TODO make this configurable
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password(serverData.RootPassword),
		},
	}, nil
}
