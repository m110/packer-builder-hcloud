package hcloud

import (
	"context"
	"time"

	"github.com/hashicorp/packer/packer"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/mitchellh/multistep"
)

type stepWaitForImage struct{}

func (s *stepWaitForImage) Run(state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(*hcloud.Client)
	ui := state.Get("ui").(packer.Ui)

	imageID := state.Get("image_id").(int)

	ctx := context.Background()

	waiter := NewWaiter(client, 2*time.Minute)
	err := waiter.WaitForImage(ctx, imageID, hcloud.ImageStatusAvailable)
	if err != nil {
		ui.Error(err.Error())
		state.Put("error", err)
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *stepWaitForImage) Cleanup(state multistep.StateBag) {
}
