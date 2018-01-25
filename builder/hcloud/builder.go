package hcloud

import (
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/mitchellh/multistep"
	"github.com/pkg/errors"
)

const BuilderID = "packer.hcloud"

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	Comm                communicator.Config `mapstructure:",squash"`

	Token      string `mapstructure:"token"`
	ServerType string `mapstructure:"server_type"`

	ImageName   string `mapstructure:"image_name"`
	SourceImage string `mapstructure:"source_image"`

	SSHKey string `mapstructure:"ssh_key"`

	ctx interpolate.Context
}

type Builder struct {
	config Config
}

func (b *Builder) Prepare(raws ...interface{}) ([]string, error) {
	err := config.Decode(&b.config, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &b.config.ctx,
	}, raws...)

	if err != nil {
		return nil, err
	}

	var errs *packer.MultiError
	errs = packer.MultiErrorAppend(errs, b.config.Comm.Prepare(&b.config.ctx)...)

	if b.config.Token == "" {
		errs = packer.MultiErrorAppend(errs, errors.New("Missing token"))
	}

	if b.config.ServerType == "" {
		errs = packer.MultiErrorAppend(errs, errors.New("Missing server type"))
	}

	if b.config.ImageName == "" {
		errs = packer.MultiErrorAppend(errs, errors.New("Missing image name"))
	}

	if b.config.SourceImage == "" {
		errs = packer.MultiErrorAppend(errs, errors.New("Missing source image"))
	}

	if len(errs.Errors) > 0 {
		return nil, errors.New(errs.Error())
	}

	return nil, nil
}

func (b *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {
	client := hcloud.NewClient(hcloud.WithToken(b.config.Token))

	state := new(multistep.BasicStateBag)
	state.Put("config", b.config)
	state.Put("client", client)
	state.Put("hook", hook)
	state.Put("ui", ui)

	steps := []multistep.Step{
		new(stepCreateInstance),
		new(stepWaitForInstance),
		&communicator.StepConnect{
			Config:    &b.config.Comm,
			Host:      commHost,
			SSHConfig: sshConfig,
		},
		new(common.StepProvision),
		new(stepCaptureImage),
	}

	runner := multistep.BasicRunner{Steps: steps}
	runner.Run(state)

	artifact := &Artifact{
		imageID:   "",
		imageName: "",
	}

	return artifact, nil
}

func (b *Builder) Cancel() {

}
