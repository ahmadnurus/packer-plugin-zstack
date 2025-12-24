package instance

import (
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/packer"

	"github.com/terraform-zstack-modules/zstack-sdk-go/pkg/client"
	"github.com/terraform-zstack-modules/zstack-sdk-go/pkg/errors"
	"github.com/terraform-zstack-modules/zstack-sdk-go/pkg/param"
)

var _ packer.Artifact = &Artifact{}

type Artifact struct {
	Image     string
	ImageName string

	zsconf *client.ZSConfig
}

// BuilderId implements packer.Artifact.
func (a *Artifact) BuilderId() string { return BuilderID }

// Destroy implements packer.Artifact.
func (a *Artifact) Destroy() error {
	client := client.NewZSClient(a.zsconf)
	if _, err := client.Login(); err != nil {
		if err != errors.ErrNotSupported {
			return err
		}
	}

	return client.DeleteImage(a.Image, param.DeleteModeEnforcing)
}

// Files implements packer.Artifact.
func (a *Artifact) Files() []string { return nil }

// Id implements packer.Artifact.
func (a *Artifact) Id() string { return a.Image }

// State implements packer.Artifact.
func (a *Artifact) State(name string) interface{} { panic("unimplemented") }

// String implements packer.Artifact.
func (a *Artifact) String() string {
	return fmt.Sprintf("The image %s was created with uuid: %s", a.ImageName, a.Image)
}
