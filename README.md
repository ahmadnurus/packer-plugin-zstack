# Packer Plugin for zstack

The zstack builder plugin can be used with HashiCorp Packer to create custom images.

## Build from source

1. Clone this GitHub repository locally.

2. Run this command from the root directory: 
```shell 
go build -ldflags="-X github.com/ahmadnurus/packer-plugin-zstack/version.VersionPrerelease=dev" -o packer-plugin-zstack
```

3. After you successfully compile, the `packer-plugin-zstack` plugin binary file is in the root directory. 

4. To install the compiled plugin, run the following command 
```shell
packer plugins install --path packer-plugin-zstack github.com/ahmadnurus/zstack
```
