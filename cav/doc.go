/*
CloudAvenue SDK Client Initialization Flow

	NewClient()
	    |
	    v
	Main client is initialized with:
	    - organization
	    - options (such as credential)
	    |
	    v
	SubClients are created (Vmware, Cerberus, S3, etc.)
	    |
	    v
	Each subClient receives:
	    - the console
	    - the credential (Auth interface)
	    |
	    v
	When an API call is made:
	    - The subClient uses credential.Headers() to set authentication headers
	    - The subClient may call credential.Refresh() to refresh the token if needed

Summary:
  - NewClient initializes the main client and injects authentication and configuration.
  - SubClients are instantiated and receive shared authentication.
  - Each subClient handles its own API calls using the provided credential.

Note:

	All subClients share the same credential instance, ensuring centralized authentication
*/
package cav
