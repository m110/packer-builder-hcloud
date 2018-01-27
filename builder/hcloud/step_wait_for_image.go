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

	ui.Say("Waiting for image to become available...")

	for {
		image, _, err := client.Image.GetByID(ctx, imageID)
		if err != nil {
			ui.Error(err.Error())
			state.Put("error", err)
			return multistep.ActionHalt
		}

		if image.Status == hcloud.ImageStatusAvailable {
			break
		}

		time.Sleep(3 * time.Second)
	}

	return multistep.ActionContinue
}

func (s *stepWaitForImage) Cleanup(state multistep.StateBag) {

}
