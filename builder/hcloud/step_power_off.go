package hcloud

import (
	"context"
	"time"

	"github.com/hashicorp/packer/packer"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/mitchellh/multistep"
)

type stepPowerOff struct{}

func (s *stepPowerOff) Run(state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(*hcloud.Client)
	ui := state.Get("ui").(packer.Ui)

	ctx := context.Background()

	serverData := state.Get("server_data").(hcloud.ServerCreateResult)
	serverID := serverData.Server.ID

	server, _, err := client.Server.GetByID(ctx, serverID)
	if err != nil {
		ui.Error(err.Error())
		state.Put("error", err)
		return multistep.ActionHalt
	}

	if server.Status == hcloud.ServerStatusOff {
		return multistep.ActionContinue
	}

	ui.Say("Forcing power off on server...")

	_, _, err = client.Server.Poweroff(ctx, server)
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

func (s *stepPowerOff) Cleanup(multistep.StateBag) {
}
