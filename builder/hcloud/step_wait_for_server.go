package hcloud

import (
	"context"
	"time"

	"github.com/hashicorp/packer/packer"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/mitchellh/multistep"
)

type stepWaitForServer struct{}

func (s *stepWaitForServer) Run(state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(*hcloud.Client)
	ui := state.Get("ui").(packer.Ui)

	serverData := state.Get("server_data").(hcloud.ServerCreateResult)
	serverID := serverData.Server.ID

	ctx := context.Background()

	ui.Say("Waiting for the server to be running...")

	waiter := NewWaiter(client, 2*time.Minute)
	err := waiter.WaitForServer(ctx, serverID, hcloud.ServerStatusRunning)
	if err != nil {
		ui.Error(err.Error())
		state.Put("error", err)
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *stepWaitForServer) Cleanup(state multistep.StateBag) {
}
