package hcloud

import (
	"context"
	"log"

	"time"

	"github.com/hashicorp/packer/packer"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/mitchellh/multistep"
)

type stepWaitForInstance struct{}

func (s *stepWaitForInstance) Run(state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(*hcloud.Client)
	ui := state.Get("ui").(packer.Ui)

	serverData := state.Get("server_data").(hcloud.ServerCreateResult)
	serverID := serverData.Server.ID

	ui.Say("Waiting for the server to become active...")

	ctx := context.Background()

	for {
		log.Printf("Checking server status...")

		server, _, err := client.Server.GetByID(ctx, serverID)
		if err != nil {
			ui.Error(err.Error())
			state.Put("error", err)
			return multistep.ActionHalt
		}

		if server.Status == hcloud.ServerStatusRunning {
			break
		}

		time.Sleep(2 * time.Second)
	}

	return multistep.ActionContinue
}

func (s *stepWaitForInstance) Cleanup(state multistep.StateBag) {

}
