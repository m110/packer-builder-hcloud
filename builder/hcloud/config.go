package hcloud

import (
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/template/interpolate"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	Comm                communicator.Config `mapstructure:",squash"`

	Token      string `mapstructure:"token"`
	ServerType string `mapstructure:"server_type"`

	ImageName   string `mapstructure:"image_name"`
	SourceImage string `mapstructure:"source_image"`

	Location   string `mapstructure:"location"`
	Datacenter string `mapstructure:"datacenter"`
	UserData   string `mapstructure:"user_data"`

	ctx interpolate.Context
}
