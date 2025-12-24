//go:generate packer-sdc struct-markdown
//go:generate packer-sdc mapstructure-to-hcl2 -type Config
package instance

import (
	"fmt"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/hashicorp/packer-plugin-sdk/uuid"
)

type Config struct {
	ctx          interpolate.Context
	Common       common.PackerConfig `mapstructure:",squash"`
	Communicator communicator.Config `mapstructure:",squash"`

	// The access key id used to authenticate to the zstack API.
	AccessKeyID string `mapstructure:"access_key_id" required:"false"`
	// The access key secret used to authenticate to the zstack API.
	AccessKeySecret string `mapstructure:"access_secret_key" required:"false"`
	// The account name used to authenticate to the zstack API.
	AccountName string `mapstructure:"account_name" required:"false"`
	// The account user name used to authenticate to the zstack API.
	AccountUserName string `mapstructure:"account_username" required:"false"`
	// The password used to authenticate to the zstack API.
	AccountPassword string `mapstructure:"password" required:"false"`
	// The zstack API Host to connect.
	APIHost string `mapstructure:"host" required:"false"`
	// Authentication type used to connect zstack API. Can be any of
	// access_key, account, account_user.
	AuthType string `mapstructure:"auth_type" required:"true"`
	// The unique name of the created image. Defaults to
	// `packer-{{timestamp}}`.
	ImageName string `mapstructure:"image_name" required:"false"`
	// The description of the created image.
	ImageDescription string `mapstructure:"image_description" required:"false"`
	// The backup storage uuid to store the created image.
	ImageBackupStorage string `mapstructure:"image_backup_storage" required:"true"`
	// A name to give the launched instance.
	// Defaults to `packer-{{uuid}}`.
	InstanceName string `mapstructure:"instance_name" required:"false"`
	// The instance root disk size in GB. This defaults to source image size.
	// NOTE: This value would become the created image size.
	InstanceRootDisk int64 `mapstructure:"instance_root_disk" required:"false"`
	// The instance offering uuid to be used by the instance.
	InstanceOffering string `mapstructure:"instance_offering" required:"true"`
	// The instance network uuid to be used by the instance.
	InstanceNetwork string `mapstructure:"instance_network" required:"true"`
	// The source image uuid to use to create the new image from.
	SourceImage string `mapstructure:"source_image" required:"true"`
}

func (c *Config) Prepare(raws ...interface{}) ([]string, error) {
	var warns []string
	var errs error

	if err := config.Decode(c, &config.DecodeOpts{
		PluginType:         BuilderID,
		Interpolate:        true,
		InterpolateContext: &c.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{
				"run_command",
			},
		},
	}, raws...); err != nil {
		return nil, err
	}

	if c.APIHost == "" {
		host, ok := os.LookupEnv("ZSTACK_HOST")
		if !ok {
			return nil, fmt.Errorf("host field or ZSTACK_HOST environment variable is required")
		}

		c.APIHost = host
	}

	switch c.AuthType {
	case "access_key":
		if c.AccessKeyID == "" {
			accessKeyID, ok := os.LookupEnv("ZSTACK_ACCESS_KEY")
			if !ok {
				return nil, fmt.Errorf("access_key_id field or ZSTACK_ACCESS_KEY environment variable is required when authentication type is %s", c.AuthType)
			}

			c.AccessKeyID = accessKeyID
		}

		if c.AccessKeySecret == "" {
			accessKeySecret, ok := os.LookupEnv("ZSTACK_SECRET_KEY")
			if !ok {
				return nil, fmt.Errorf("access_secret_key field or ZSTACK_SECRET_KEY environment variable is required when authentication type is %s", c.AuthType)
			}

			c.AccessKeySecret = accessKeySecret
		}
	case "account":
		if c.AccountName == "" {
			accountName, ok := os.LookupEnv("ZSTACK_ACCOUNT_NAME")
			if !ok {
				return nil, fmt.Errorf("account_name field or ZSTACK_ACCOUNT_NAME environment variable is required when authentication type is %s", c.AuthType)
			}

			c.AccountName = accountName
		}

		if c.AccountPassword == "" {
			password, ok := os.LookupEnv("ZSTACK_ACCOUNT_PASSWORD")
			if !ok {
				return nil, fmt.Errorf("password field or ZSTACK_ACCOUNT_PASSWORD environment variable is required when authentication type is %s", c.AuthType)
			}

			c.AccountPassword = password
		}
	case "account_user":
		if c.AccountName == "" {
			accountName, ok := os.LookupEnv("ZSTACK_ACCOUNT_NAME")
			if !ok {
				return nil, fmt.Errorf("account_name field or ZSTACK_ACCOUNT_NAME environment variable is required when authentication type is %s", c.AuthType)
			}

			c.AccountName = accountName
		}

		if c.AccountUserName == "" {
			accountUserName, ok := os.LookupEnv("ZSTACK_ACCOUNT_USERNAME")
			if !ok {
				return nil, fmt.Errorf("account_username field or ZSTACK_ACCOUNT_USERNAME environment variable is required when authentication type is %s", c.AuthType)
			}

			c.AccountUserName = accountUserName
		}

		if c.AccountPassword == "" {
			password, ok := os.LookupEnv("ZSTACK_ACCOUNT_PASSWORD")
			if !ok {
				return nil, fmt.Errorf("password field or ZSTACK_ACCOUNT_PASSWORD environment variable is required when authentication type is %s", c.AuthType)
			}

			c.AccountPassword = password
		}
	default:
		return nil, fmt.Errorf("unsupported authentication type - %s", c.AuthType)
	}

	if c.ImageName == "" {
		img, err := interpolate.Render("packer-{{timestamp}}", nil)
		if err != nil {
			errs = packer.MultiErrorAppend(errs, fmt.Errorf("unable to parse image name: %s ", err))
		} else {
			c.ImageName = img
		}
	}

	if c.ImageDescription == "" {
		c.ImageDescription = "Created by Packer"
	}

	if c.InstanceName == "" {
		c.InstanceName = fmt.Sprintf("packer-%s", uuid.TimeOrderedUUID())
	}

	return warns, errs
}
