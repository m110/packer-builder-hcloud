package hcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer/packer"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/mitchellh/multistep"
)

type stepCaptureImage struct{}

func (s *stepCaptureImage) Run(state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(*hcloud.Client)
	config := state.Get("config").(Config)
	ui := state.Get("ui").(packer.Ui)

	serverData := state.Get("server_data").(hcloud.ServerCreateResult)
	serverID := serverData.Server.ID

	ui.Say("Waiting for the server to become active...")

	ctx := context.Background()

	server, _, err := client.Server.GetByID(ctx, serverID)
	if err != nil {
		ui.Error(err.Error())
		state.Put("error", err)
		return multistep.ActionHalt
	}

	result, _, err := client.Server.CreateImage(ctx, server, &hcloud.ServerCreateImageOpts{
		Type:        hcloud.ImageTypeSnapshot,
		Description: &config.ImageName,
	})
	if err != nil {
		ui.Error(err.Error())
		state.Put("error", err)
		return multistep.ActionHalt
	}

	ui.Say(fmt.Sprintf("Created image: %s ID: %s", result.Image.Name, result.Image.ID))

	return multistep.ActionContinue
}

func (s *stepCaptureImage) Cleanup(state multistep.StateBag) {
}
