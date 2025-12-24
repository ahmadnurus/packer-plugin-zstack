# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

packer {
  required_plugins {
    zstack = {
      version = ">=v0.0.1"
      source  = "github.com/ahmadnurus/zstack"
    }
  }
}

variable "access_key" {
  type = string
}
variable "secret_key" {
  type = string
}
variable "host" {
  type = string
}

source "zstack" "example" {
  access_key_id        = var.access_key
  access_secret_key    = var.secret_key
  auth_type            = "access_key"
  communicator         = "ssh"
  host                 = var.host
  image_name           = "test-packer"
  image_backup_storage = "be67d461b2e8487e8aa5b61ff4853d34"
  instance_offering    = "cadd5ed01a8c431080628d4dfde4fdc5"
  instance_network     = "a60a408253ac44b198c3753c9b8b921f"
  source_image         = "b87e2724243943e4b6916231ab11b3a5"
  ssh_username         = "ubuntu"
  ssh_port             = 22
}

build {
  sources = ["source.zstack.example"]
  provisioner "shell" {
    inline = [
      "echo Hello From ${source.type} ${source.name}"
    ]
  }
}
