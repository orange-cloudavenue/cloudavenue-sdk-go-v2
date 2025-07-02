package mock

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/consoles"
)

const (
	mockOrg = "cav01ev01ocb0001234"
)

var (
	pathPrefix = map[cav.SubClientName]string{
		cav.ClientVmware:         "/cloudapi",
		cav.ClientCerberus:       "/api/customers",
		cav.ClientNetbackup:      "/netbackup",
		cav.SubClientName("ihm"): "/ihm",
		cav.SubClientName("s3"):  "/s3",
	}
)

func NewClient() (cav.Client, error) {
	// Mock implementation for testing purposes

	// Get All endpoints available in the endpoint package
	// Create an handler for each endpoint
	// Each handler should return a mock response
	// This is a placeholder for the actual implementation

	endpoints := cav.GetEndpointsUncategorized()

	mux := chi.NewRouter()

	for _, ep := range endpoints {
		switch ep.Method {
		case cav.MethodGET:
			if ep.MockResponseFuncIsDefined() {
				mux.Get(buildPath(ep.SubClient, ep.PathTemplate), ep.GetMockResponseFunc())
				continue
			}

			mux.Get(buildPath(ep.SubClient, ep.PathTemplate), cav.GetDefaultMockResponseFunc(ep))
		case cav.MethodPOST:
			mux.Post(buildPath(ep.SubClient, ep.PathTemplate), func(w http.ResponseWriter, _ *http.Request) {
				// Return a mock response
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Mock response"))
			})
		}
	}

	hts := httptest.NewServer(mux)

	log.Default().Println("Mock server started at", hts.URL)

	nC, err := cav.NewClient(
		mockOrg,
		cav.WithCustomEndpoints(consoles.Services{
			IHM: consoles.Service{
				Enabled:  true,
				Endpoint: hts.URL + "/ihm",
			},
			APIVCD: consoles.Service{
				Enabled:  true,
				Endpoint: hts.URL + "/cloudapi",
			},
			APICerberus: consoles.Service{
				Enabled:  true,
				Endpoint: hts.URL + "/api/customers",
			},
			S3: consoles.Service{
				Enabled:  true,
				Endpoint: hts.URL + "/s3",
			},
			Netbackup: consoles.Service{
				Enabled:  true,
				Endpoint: hts.URL + "/netbackup",
			},
		}),
		cav.WithCloudAvenueCredential("mockuser", "mockpassword"),
	)
	if err != nil {
		return nil, err
	}

	return nC, nil
}

func buildPath(subClient cav.SubClientName, path string) string {
	if !strings.HasPrefix(path, pathPrefix[subClient]) {
		return pathPrefix[subClient] + path
	}
	return path
}

func SetMockResponse(ep *cav.Endpoint, mockResponseData any, mockResponseStatusCode *int) {
	if ep.MockResponseFuncIsDefined() {
		log.Default().Println("Mock response already defined for endpoint", ep.Name)
		return
	}

	ep.SetMockResponse(mockResponseData, mockResponseStatusCode)
	log.Default().Printf("Mock response set for endpoint %s with status code %d", ep.Name, mockResponseStatusCode)
}

func CleanMockResponses() {
	endpoints := cav.GetEndpointsUncategorized()
	for _, ep := range endpoints {
		if ep.MockResponseFuncIsDefined() {
			ep.CleanMockResponse()
			log.Default().Printf("Mock response cleaned for endpoint %s", ep.Name)
		}
	}
}

var GetEndpoint = cav.GetEndpoint
