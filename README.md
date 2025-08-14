# CloudAvenue SDK V2

> [!CAUTION]
> This project is in early development and may change significantly.

CloudAvenue SDK V2 is a Go client library for interacting with the Orange CloudAvenue APIs.  
It provides authentication, client configuration, and helpers to access CloudAvenue services.

## Example Usage

Below is a minimal example showing how to initialize a CloudAvenue SDK client and perform a simple request using the `vdc` client:

```go
package main

import (
    "context"
    "fmt"

    "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/vdc/v1"
    "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
    "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

func main() {
    ctx := context.Background()

    // Initialize the main client with your organization and credentials
    client, err := cav.NewClient(
        "cav01ev01ocb0001234",
        cav.WithCloudAvenueCredential("your_username", "your_password"),
    )
    if err != nil {
        fmt.Println("Error creating client:", err)
        return
    }
    defer client.Close()

    // Create an VDC client
    ev, err := vdc.New(client)
    if err != nil {
      fmt.Println("Error creating VDC client:", err)
      return
    }
    
    vdcCreated, err := ev.CreateVDC(ctx, types.ParamsCreateVDC{
       Name:                "test-vdc",
       Description:         "Test VDC",
       ServiceClass:        "ECO",
       BillingModel:        "PAYG",
       DisponibilityClass:  "ONE-ROOM",
       StorageBillingModel: "PAYG",
       Vcpu:                5,
       Memory:              50,
       StorageProfiles: []types.ParamsCreateVDCStorageProfile{
         {
            Class:   "silver",
            Limit:   100,
            Default: true,
         },
       },
    })
    if err != nil {
     fmt.Println("Error creating VDC:", err)
     return
    }

    fmt.Printf("VDC created: %+v\n", vdcCreated)
}
```

## License

This project is licensed under the [Mozilla Public License 2.0](LICENSE).

---

Orange CloudAvenue

---

## Documentation & Contribution

- [Contribution Guide](./CONTRIBUTING.md)
- [Coding Guidelines](./GUIDELINE.md)
- [Project Architecture](./ARCHITECTURE.md)
