# packer-builder-hcloud
[Packer](https://packer.io/) builder plugin for [Hetzner Cloud](https://hetzner.cloud/).

This plugin can be used to build provisioned images (snapshots) for Hetzner Cloud.
Those can be later used for rapid deploying using [Terraform](https://terraform.io/) (check the [official plugin](https://github.com/hetznercloud/terraform-provider-hcloud).

# Building

You'll need [dep](https://github.com/golang/dep) to install dependencies.

Clone the repository and run:

```
dep ensure
go build -o ~/.packer.d/plugins/packer-builder-hcloud
```

# Example template

```
{
  "builders": [
      {
          "type": "hcloud",
          "token": "API_TOKEN",
          "server_type": "cx11",
          "source_image": "ubuntu-16.04",
          "image_name": "some-image",
      }
  ]
}
```

You can find server types and soure images querying [the API](https://docs.hetzner.cloud/#resources-server-types).

# Known issues

* For some reason, `ansible-remote` provider works only with paramiko connection.
