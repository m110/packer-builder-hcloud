package hcloud

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/packer/packer"
	"github.com/mitchellh/multistep"
)

type stepCreateSSHKey struct {
	PrivateKeyFile string
}

func (s *stepCreateSSHKey) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)

	privateKeyBytes, err := ioutil.ReadFile(s.PrivateKeyFile)
	if err != nil {
		ui.Error(err.Error())
		state.Put("error", fmt.Errorf("Error loading configured private key file: %s", err))
		return multistep.ActionHalt
	}

	state.Put("ssh_private_key", string(privateKeyBytes))

	return multistep.ActionContinue
}

func (s *stepCreateSSHKey) Cleanup(state multistep.StateBag) {
}
