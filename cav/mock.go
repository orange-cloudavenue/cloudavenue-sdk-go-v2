package cav

import (
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/auth"
	subclient "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/internal/subClient"
)

// TODO move to internal/subClient/mock.go
// WithMock sets the credential for the client.
func WithMock() ClientOption {
	return func(s *settings) error {
		if s.SubClients[mock] == nil {
			s.SubClients[mock] = subclient.Clients[mock]
		}

		s.SubClients[mock].SetConsole(s.Console)
		s.SubClients[mock].SetCredential(auth.NewMockAuth(map[string]string{
			"X-Mock": "mock",
		}))

		return nil
	}
}
