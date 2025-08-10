package vdc

import (
	"fmt"

	"resty.dev/v3"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/common-go/extractor"
	"github.com/orange-cloudavenue/common-go/validators"
)

//go:generate endpoint-generator -path storage_profile_endpoints.go

func init() {
	// * ListStorageProfiles
	cav.Endpoint{
		DocumentationURL: "https://developer.broadcom.com/xapis/vmware-cloud-director-api/latest/doc/queries/orgVdcStorageProfile.html",
		Name:             "ListStorageProfiles",
		Description:      "List VDC Storage Profiles",
		Method:           cav.MethodGET,
		SubClient:        cav.ClientVmware,
		PathTemplate:     "/api/query/",
		QueryParams: []cav.QueryParam{
			{
				Name:        "filter",
				Description: "ID of the VDC to get storage profiles for.",
				ValidatorFunc: func(value string) error {
					return validators.New().Var(value, "urn=vdc")
				},
				TransformFunc: func(value string) (string, error) {
					// vdc-id require UUID format and not urn format
					v, err := extractor.ExtractUUID(value)
					if err != nil {
						return "", err
					}
					return fmt.Sprintf("vdc==%s", v), nil
				},
			},
			{
				Name:        "pageSize",
				Description: "The number of items per page.",
				Value:       "30",
			},
			{
				Name:        "format",
				Description: "The format of the response.",
				Value:       "records",
			},
			{
				Name:        "type",
				Description: "The type of object to query",
				Value:       "orgVdcStorageProfile",
			},
			{
				Name:        "sortAsc",
				Description: "Sort the results in ascending order.",
				Value:       "name",
			},
		},
		BodyResponseType: apiResponseListStorageProfiles{},
		RequestMiddlewares: []resty.RequestMiddleware{
			func(_ *resty.Client, req *resty.Request) error {
				// Set the Accept header to application/*+json;version=38.1
				req.SetHeader("Accept", "application/*+json;version=38.1")
				return nil
			},
		},
		// MockResponseFunc: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 	resp := &apiResponseListStorageProfiles{
		// 		StorageProfiles: make([]apiResponseListStorageProfile, 0),
		// 	}

		// 	// If QueryParam "filter" is set, return a filtered response
		// 	if r.URL.Query().Get("filter") != "" {
		// 		filter := r.URL.Query().Get("filter")

		// 		filterParts := strings.Split(filter, "==")

		// 		r := apiResponseListStorageProfile{}
		// 		generator.MustStruct(r)

		// 		r.ID = func() string {
		// 			if filterParts[0] == "vdc" {
		// 				return filterParts[1]
		// 			}
		// 			return ""
		// 		}()
		// 	} else {
		// 		generator.MustStruct(&resp)
		// 	}
		// 	// json encode
		// 	w.Header().Set("Content-Type", "application/json")
		// 	respJ, err := json.Marshal(resp)
		// 	if err != nil {
		// 		http.Error(w, err.Error(), http.StatusInternalServerError)
		// 		return
		// 	}
		// 	_, err = w.Write(respJ)
		// 	if err != nil {
		// 		http.Error(w, err.Error(), http.StatusInternalServerError)
		// 		return
		// 	}
		// }),
	}.Register()
}
