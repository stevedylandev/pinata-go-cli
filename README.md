# Pinata Go CLI

Welcome to the Pinata Go CLI! This is a rewrite of our node.js cli but written in Go. This is still in active development so please let us know if you have any questions! :) 

## Installation

We are currently working on the build flow for binaries to make installation easier, but for now we recommend building from source.

To do this make sure you have [Go](https://go.dev/) installed on your computer and the following command returns a version:
```bash
go version
```

Then paste and run the following into your terminal:

```bash
git clone https://github.com/stevedylandev/pinata-go-cli && cd pinata-go-cli && go install -0 pinata
```

## Usage

With the CLI installed you will first need to authenticate it with your [Pinata JWT](https://docs.pinata.cloud/docs/api-keys)

```bash
pinata auth <your-jwt>
```

After its been authenticated you can now upload using the `upload` command or `u` for short, then pass in the path to the file or folder you want to upload.

```bash
pinata upload ~/Pictures/somefolder/image.png
```

## Contact 

If you have any questions please feel free to reach out to us! 

[team@pinata.cloud](mailto:team@pinata.cloud)
