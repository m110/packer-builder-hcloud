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

	ui.Say("Waiting for the server to become active...")

	ctx := context.Background()

	for {
		server, _, err := client.Server.GetByID(ctx, serverID)
		if err != nil {
			ui.Error(err.Error())
			state.Put("error", err)
			return multistep.ActionHalt
		}

		if server.Status == hcloud.ServerStatusRunning {
			break
		}

		time.Sleep(3 * time.Second)
	}

	return multistep.ActionContinue
}

func (s *stepWaitForServer) Cleanup(state multistep.StateBag) {

}
