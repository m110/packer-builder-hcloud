package hcloud

import (
	"context"
	"time"

	"github.com/hashicorp/packer/packer"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/mitchellh/multistep"
)

type stepShutdown struct{}

func (s *stepShutdown) Run(state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(*hcloud.Client)
	ui := state.Get("ui").(packer.Ui)

	ctx := context.Background()

	serverData := state.Get("server_data").(hcloud.ServerCreateResult)
	serverID := serverData.Server.ID

	ui.Say("Shutting down server...")

	server, _, err := client.Server.GetByID(ctx, serverID)
	if err != nil {
		ui.Error(err.Error())
		state.Put("error", err)
		return multistep.ActionHalt
	}

	_, _, err = client.Server.Shutdown(ctx, server)
	if err != nil {
		ui.Error(err.Error())
		state.Put("error", err)
		return multistep.ActionHalt
	}

	waiter := NewWaiter(client, 30*time.Second)
	err = waiter.WaitForServer(ctx, serverID, hcloud.ServerStatusOff)
	if err != nil {
		ui.Error(err.Error())
		state.Put("error", err)
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *stepShutdown) Cleanup(multistep.StateBag) {
}
