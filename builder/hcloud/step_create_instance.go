package hcloud

import (
	"context"

	"fmt"

	"time"

	"github.com/hashicorp/packer/packer"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/mitchellh/multistep"
)

type stepCreateInstance struct {
	instanceID int
}

func (s *stepCreateInstance) Run(state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(*hcloud.Client)
	config := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)

	ctx := context.Background()

	serverType, _, err := client.ServerType.Get(ctx, config.ServerType)
	if err != nil {
		ui.Error(err.Error())
		state.Put("error", err)
		return multistep.ActionHalt
	}

	sourceImage, _, err := client.Image.Get(ctx, config.SourceImage)
	if err != nil {
		ui.Error(err.Error())
		state.Put("error", err)
		return multistep.ActionHalt
	}

	sshKey, _, err := client.SSHKey.Get(ctx, config.SSHKey)
	if err != nil {
		ui.Error(err.Error())
		state.Put("error", err)
		return multistep.ActionHalt
	}

	ui.Say("Creating new server")

	serverData, _, err := client.Server.Create(ctx, hcloud.ServerCreateOpts{
		Name:       fmt.Sprintf("packer-hcloud-%s", time.Now().Unix()),
		ServerType: serverType,
		Image:      sourceImage,
		SSHKeys:    []*hcloud.SSHKey{sshKey},
		// TODO
		// Location:
		// Datacenter
		// UserData
	})

	if err != nil {
		ui.Error(err.Error())
		state.Put("error", err)
		return multistep.ActionHalt
	}

	state.Put("server_data", serverData)
	s.instanceID = serverData.Server.ID

	ui.Say(fmt.Sprintf("Crated server %d", s.instanceID))

	return multistep.ActionContinue
}

func (s *stepCreateInstance) Cleanup(state multistep.StateBag) {
	client := state.Get("client").(*hcloud.Client)
	ui := state.Get("ui").(packer.Ui)

	if s.instanceID <= 0 {
		return
	}

	ui.Say(fmt.Sprintf("Waiting for server %d to be destroyed...", s.instanceID))

	ctx := context.Background()

	server, _, err := client.Server.GetByID(ctx, s.instanceID)

	_, err = client.Server.Delete(ctx, server)
	if err != nil {
		ui.Error(err.Error())
		state.Put("error", err)
		return
	}
}
