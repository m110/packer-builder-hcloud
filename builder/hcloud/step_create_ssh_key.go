package hcloud

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/packer/packer"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/mitchellh/multistep"
	"golang.org/x/crypto/ssh"
)

type stepCreateSSHKey struct {
	keyID int
}

func (s *stepCreateSSHKey) Run(state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(*hcloud.Client)
	ui := state.Get("ui").(packer.Ui)

	ui.Say("Creating temporary ssh key")

	priv, err := rsa.GenerateKey(rand.Reader, 2048)

	privDER := x509.MarshalPKCS1PrivateKey(priv)
	privBLK := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	state.Put("private_key", string(pem.EncodeToMemory(&privBLK)))

	pub, err := ssh.NewPublicKey(&priv.PublicKey)
	if err != nil {
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	publicKey := string(ssh.MarshalAuthorizedKey(pub))
	keyName := fmt.Sprintf("packer-hcloud-%d", time.Now().Unix())

	ctx := context.Background()

	key, _, err := client.SSHKey.Create(ctx, hcloud.SSHKeyCreateOpts{
		Name:      keyName,
		PublicKey: publicKey,
	})
	if err != nil {
		err := fmt.Errorf("error creating temporary ssh key: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	s.keyID = key.ID
	state.Put("ssh_key_id", key.ID)

	log.Printf("Created temporary ssh key: %d (%s)", key.ID, keyName)

	return multistep.ActionContinue
}

func (s *stepCreateSSHKey) Cleanup(state multistep.StateBag) {
	if s.keyID == 0 {
		return
	}

	client := state.Get("client").(*hcloud.Client)
	ui := state.Get("ui").(packer.Ui)

	ctx := context.Background()

	key, _, err := client.SSHKey.GetByID(ctx, s.keyID)
	if err != nil {
		ui.Error(fmt.Sprintf("error getting ssh key: %s", err))
		return
	}

	ui.Say("Deleting temporary ssh key...")
	_, err = client.SSHKey.Delete(ctx, key)
	if err != nil {
		ui.Error(fmt.Sprintf("Error deleting ssh key: %s", err))
	}
}
