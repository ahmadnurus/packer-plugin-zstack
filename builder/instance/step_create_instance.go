package instance

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/terraform-zstack-modules/zstack-sdk-go/pkg/client"
	"github.com/terraform-zstack-modules/zstack-sdk-go/pkg/errors"
	"github.com/terraform-zstack-modules/zstack-sdk-go/pkg/param"
)

var _ multistep.Step = &StepCreateInstance{}

// StepCreateInstance represents a Packer build step that creates vm instance.
type StepCreateInstance struct{}

// Cleanup implements multistep.Step.
func (s *StepCreateInstance) Cleanup(state multistep.StateBag) {
	ui := state.Get("ui").(packer.Ui)
	config := state.Get("config").(*Config)
	zsconf := state.Get("zsconf").(*client.ZSConfig)

	client := client.NewZSClient(zsconf)
	_, _ = client.Login()

	if _, ok := state.GetOk("instance_uuid"); !ok {
		return
	}

	ui.Say("Deleting instance...")
	if err := client.DestroyVmInstance(state.Get("instance_uuid").(string), param.DeleteModeEnforcing); err != nil {
		ui.Errorf("Error deleting instance %s with uuid: %s. Please delete it manually.", config.InstanceName, state.Get("instance_uuid").(string))
		return
	}

	ui.Message("Instance has been deleted!")
	state.Put("instance_uuid", "")
	state.Put("instance_name", "")
}

// Run implements multistep.Step.
func (s *StepCreateInstance) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	config := state.Get("config").(*Config)
	zsconf := state.Get("zsconf").(*client.ZSConfig)

	client := client.NewZSClient(zsconf)
	if _, err := client.Login(); err != nil {
		if err != errors.ErrNotSupported {
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
	}

	ui.Say("Creating instance...")
	systemtags := make([]string, 0)
	sshkey := "sshkey::" + string(config.Communicator.SSHPublicKey)
	systemtags = append(systemtags, sshkey)
	createInstanceParam := param.CreateVmInstanceDetailParam{
		Name:                 config.InstanceName,
		ImageUUID:            config.SourceImage,
		InstanceOfferingUUID: config.InstanceOffering,
		L3NetworkUuids:       []string{config.InstanceNetwork},
	}

	if config.InstanceRootDisk > 1 {
		size := config.InstanceRootDisk * 1024 * 1024 * 1024
		createInstanceParam.RootDiskSize = &size
	}

	instance, err := client.CreateVmInstance(param.CreateVmInstanceParam{
		BaseParam: param.BaseParam{
			SystemTags: systemtags,
		},
		Params: createInstanceParam,
	})

	if err != nil {
		err := fmt.Errorf("error creating instance: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Message("Instance has been created!")
	state.Put("instance_uuid", instance.UUID)
	state.Put("instance_root_volume_uuid", instance.RootVolumeUUID)
	if config.Common.PackerDebug {
		ui.Message(fmt.Sprintf("Created instance %s with uuid: %s", instance.Name, instance.UUID))
	}

	if len(instance.VMNics) == 0 {
		ui.Errorf("Instance created without vmnic")
		return multistep.ActionHalt
	}

	ui.Message("Using first instance vmnic IP address as host communicator - IP Adress: " + instance.VMNics[0].IP)
	state.Put("instance_ip", instance.VMNics[0].IP)
	return multistep.ActionContinue
}
