# CloudAvenue SDK V2

> [!CAUTION]
> This project is in early development and may change significantly.

CloudAvenue SDK V2 is a Go client library for interacting with the Orange CloudAvenue APIs.  
It provides authentication, client configuration, and helpers to access CloudAvenue services.

## Example Usage

Below is a minimal example showing how to initialize a CloudAvenue SDK client and perform a simple request using the `org` API:

```go
package main

import (
    "context"
    "fmt"

    "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/api/org/v1"
    "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
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

    // Create an org API client
    orgClient, err := org.New(client)
    if err != nil {
        fmt.Println("Error creating org client:", err)
        return
    }

    // Perform a demo request (replace with your own URN)
    _, err = orgClient.DemoRequest(ctx, "urn:vcloud:org:9bf2eb9d-78fb-476b-a15a-1f4b7da4d132")
    if err != nil {
        fmt.Println("Error making request:", err)
        return
    }

    fmt.Println("Request completed successfully")
}
```

## License

This project is licensed under the [Mozilla Public License 2.0](LICENSE).

---

Orange CloudAvenue
