# Mythic Beasts Client

[![Go Reference](https://pkg.go.dev/badge/github.com/paultibbetts/mythicbeasts-client-go.svg)](https://pkg.go.dev/github.com/paultibbetts/mythicbeasts-client-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/paultibbetts/mythicbeasts-client-go)](https://goreportcard.com/report/github.com/paultibbetts/mythicbeasts-client-go)
[![Test Status](https://github.com/paultibbetts/mythicbeasts-client-go/actions/workflows/tests.yaml/badge.svg?branch=main)](https://github.com/paultibbetts/mythicbeasts-client-go/actions/workflows/tests.yaml)

mythicbeasts-client-go is a Go client for the Mythic Beasts [Raspberry Pi](https://www.mythic-beasts.com/support/api/raspberry-pi), [VPS](https://www.mythic-beasts.com/support/api/vps), and [Proxy](https://www.mythic-beasts.com/support/api/proxy) APIs.

## Installation

```bash
go get github.com/paultibbetts/mythicbeasts-client-go@latest
```

## Usage

### Quick start

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/paultibbetts/mythicbeasts-client-go"
	"github.com/paultibbetts/mythicbeasts-client-go/pi"
)

func main() {
	c, err := mythicbeasts.NewClient("YOUR_API_KEYID", "YOUR_API_SECRET")
	if err != nil {
		log.Fatal(err)
	}

	please := pi.CreateRequest{
		Model:      4,
		Memory:     4098,
		DiskSize:   10,
		OSImage:    "rpi-bookworm-arm64",
		SSHKey:     "ssh-ed25519 ... code@paultibbetts.uk",
		WaitForDNS: true,
	}

	ctx := context.Background()
	piServer, err := c.Pi().Create(ctx, "example-pi", please)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Pi IPv6: %s", piServer.IP)
}
```

### Authentication

Create a new [API key](https://www.mythic-beasts.com/customer/api-users) and construct a new client using your API key ID and secret:

```go
c, err := mythicbeasts.NewClient("YOUR_API_KEYID", "YOUR_API_SECRET")
if err != nil {
	// handle error
}
```

You can manage your API tokens on [the Mythic Beasts site](https://www.mythic-beasts.com/customer/api-users).

### Server status values

`Server.Status` carries the power state. The API documents no set of values, so
this is what the client has observed and it may not be exhaustive:

| `Status` | Meaning |
| --- | --- |
| `running` | The server is up |
| `paused` | A transient of about two seconds during power-on |
| `shut down` | Powered off, whether dormant or freshly stopped |

Note the space in `shut down`. The `shutdown` power action
(`vps.PowerActionShutdown`) has no space: one is the action you send, the other
the state that comes back, and they are different strings.

Two rules follow from the transient. To wait for a shutdown, wait for `Status`
to leave `running`. To wait for a boot, wait for `running` itself rather than
for the absence of anything else, since `paused` is a server on its way up.

`running` appears about ten seconds before the server accepts connections, so a
first connect after a power-on needs retrying whatever the API says.

### Waiting after a power request

`RebootWithGrace` and `ShutdownWithGrace` pause for a fixed period, two minutes
by default, and that pause is unconditional: it runs to completion however
quickly the server stops or comes back, and neither call checks. To know that a
server actually stopped, poll `Get` until `Status` leaves `running`.

### Making a VPS dormant

A dormant VPS keeps its storage and IP addresses but is otherwise decommissioned. The transition is a forced power off and discards whatever is in RAM, so shut a running server down gracefully first:

```go
if _, err := c.VPS().ShutdownWithGrace(ctx, "example-vps", 0); err != nil {
	// handle error
}

if _, err := c.VPS().MakeDormant(ctx, "example-vps"); err != nil {
	// handle error
}
```

`ShutdownWithGrace` requests an ACPI shutdown and waits `vps.DefaultShutdownGracePeriod` when the grace period is zero. Reactivating a dormant server needs a product code from `GetProducts`:

```go
if _, err := c.VPS().Reactivate(ctx, "example-vps", "VPSX1"); err != nil {
	// handle error
}
```

### Deleting an absent server

Deleting a VPS or Pi that the API does not have returns `ErrNotFound`. Callers that treat an already-absent server as success match with `errors.Is(err, vps.ErrNotFound)`.

## Versioning

This project is pre-1.0 and minor releases may include breaking changes.

[Semantic versioning](https://semver.org/) is used for tags and a v1.0.0 will signal a stable API.

## Contributing

Contributions are welcome. This project uses [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).

## License

MIT 2025 Paul Tibbetts.
