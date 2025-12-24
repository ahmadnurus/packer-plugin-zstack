<!--
  Include a short description about the builder. This is a good place
  to call out what the builder does, and any requirements for the given
  builder environment. See https://www.packer.io/docs/builder/null
-->

The zstack packer builder is able to create custom images for use 
with zstack vm instance based on existing images.

<!-- Builder Configuration Fields -->

**Required**

<!-- Code generated from the comments of the Config struct in builder/instance/config.go; DO NOT EDIT MANUALLY -->

- `auth_type` (string) - Authentication type used to connect zstack API. Can be any of
  access_key, account, account_user.

- `image_backup_storage` (string) - The backup storage uuid to store the created image.

- `instance_offering` (string) - The instance offering uuid to be used by the instance.

- `instance_network` (string) - The instance network uuid to be used by the instance.

- `source_image` (string) - The source image uuid to use to create the new image from.

<!-- End of code generated from the comments of the Config struct in builder/instance/config.go; -->



<!--
  Optional Configuration Fields

  Configuration options that are not required or have reasonable defaults
  should be listed under the optionals section. Defaults values should be
  noted in the description of the field
-->

**Optional**

<!-- Code generated from the comments of the Config struct in builder/instance/config.go; DO NOT EDIT MANUALLY -->

- `access_key_id` (string) - The access key id used to authenticate to the zstack API.

- `access_secret_key` (string) - The access key secret used to authenticate to the zstack API.

- `account_name` (string) - The account name used to authenticate to the zstack API.

- `account_username` (string) - The account user name used to authenticate to the zstack API.

- `password` (string) - The password used to authenticate to the zstack API.

- `host` (string) - The zstack API Host to connect.

- `image_name` (string) - The unique name of the created image. Defaults to
  `packer-{{timestamp}}`.

- `image_description` (string) - The description of the created image.

- `instance_name` (string) - A name to give the launched instance.
  Defaults to `packer-{{uuid}}`.

- `instance_root_disk` (int64) - The instance root disk size in GB. This defaults to source image size.
  NOTE: This value would become the created image size.

<!-- End of code generated from the comments of the Config struct in builder/instance/config.go; -->



<!--
  A basic example on the usage of the builder. Multiple examples
  can be provided to highlight various build configurations.

-->
### Example Usage


```hcl
 source "scaffolding" "example" {
   mock = "bird"
 }

 build {
   sources = ["source.scaffolding.example"]
 }
```
