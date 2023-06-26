# Mainflux Go SDK

Go SDK, a Go driver for Mainflux `mfxkit` API.

## Installation

Import `"github.com/mainflux/mfxkit/sdk/go"` into your Go package.

```go
import "github.com/mainflux/mfxkit/pkg/sdk/go"
```

Then call SDK Go functions to interact with the system.

## API Reference

```go
// NewSDK returns new mainflux SDK instance.
//
// conf := sdk.Config{
//      MFxkitURL: "http://localhost:9099",
//      MsgContentType: sdk.CTJSON,
//      TLSVerification: false,
//  }
//
//  sdk := sdk.NewSDK(conf)
func NewSDK(conf Config)

// Ping sends a ping request to Mainflux.
//
// Mainflux responds with a greeting message.
//
// example:
//
//  greeting, err := sdk.Ping("my-secret")
//  if err != nil {
//      fmt.Println(err)
//  }
//  fmt.Println(greeting)
Ping(secret string) (string, error)

// Health sends a health request to Mainflux.
//
// Mainflux responds with service status.
//
// example:
//
//  health, err := sdk.Health()
//  if err != nil {
//      fmt.Println(err)
//  }
//  fmt.Println(health)
Health() (HealthInfo, error)
```
