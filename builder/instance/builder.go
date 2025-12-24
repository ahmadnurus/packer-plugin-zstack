package instance

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/multistep/commonsteps"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/terraform-zstack-modules/zstack-sdk-go/pkg/client"
)

const BuilderID = "packer.zstack"

var _ packer.Builder = &Builder{}

type Builder struct {
	config Config
	runner multistep.Runner
}

// ConfigSpec implements packer.Builder.
func (b *Builder) ConfigSpec() hcldec.ObjectSpec {
	return b.config.FlatMapstructure().HCL2Spec()
}

// Prepare implements packer.Builder.
func (b *Builder) Prepare(raws ...interface{}) ([]string, []string, error) {
	var warns []string
	var errs error

	warns, errs = b.config.Prepare(raws...)
	if errs != nil {
		return nil, warns, errs
	}

	return nil, warns, errs
}

// Run implements packer.Builder.
func (b *Builder) Run(ctx context.Context, ui packer.Ui, hook packer.Hook) (packer.Artifact, error) {
	zsconf := client.DefaultZSConfig(b.config.APIHost)
	zsconf.ProxyFunc(func(r *http.Request) (*url.URL, error) {
		r.Header.Set("User-Agent", "packer-zstack")
		return r.URL, nil
	})

	switch b.config.AuthType {
	case "access_key":
		zsconf = zsconf.AccessKey(b.config.AccessKeyID, b.config.AccessKeySecret)
	case "account":
		zsconf = zsconf.LoginAccount(b.config.AccountName, b.config.AccountPassword)
	case "account_user":
		zsconf = zsconf.LoginAccountUser(b.config.AccountName, b.config.AccountName, b.config.AccountPassword)
	default:
		return nil, fmt.Errorf("unsupported authentication type - %s", b.config.AuthType)
	}

	if zsconf == nil {
		fmt.Println("is nil")
	}

	state := new(multistep.BasicStateBag)
	state.Put("config", &b.config)
	state.Put("zsconf", zsconf)
	state.Put("hook", hook)
	state.Put("ui", ui)

	steps := []multistep.Step{
		&communicator.StepSSHKeyGen{
			CommConf:            &b.config.Communicator,
			SSHTemporaryKeyPair: b.config.Communicator.SSHTemporaryKeyPair,
		},
		&StepCreateInstance{},
		&communicator.StepConnect{
			Config:    &b.config.Communicator,
			Host:      communicator.CommHost(b.config.Communicator.Host(), "instance_ip"),
			SSHConfig: b.config.Communicator.SSHConfigFunc(),
		},
		&commonsteps.StepProvision{},
		&commonsteps.StepCleanupTempKeys{
			Comm: &b.config.Communicator,
		},
		&StepCreateImage{},
	}

	b.runner = commonsteps.NewRunner(steps, b.config.Common, ui)
	b.runner.Run(ctx, state)
	if err, ok := state.GetOk("error"); ok {
		return nil, err.(error)
	}

	if _, ok := state.GetOk("image"); !ok {
		log.Println("Failed to find image in state. skipping..")
		return nil, nil
	}

	artifact := &Artifact{
		Image:     state.Get("image_uuid").(string),
		ImageName: state.Get("image_name").(string),
		zsconf:    zsconf,
	}

	return artifact, nil
}
