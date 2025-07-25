module github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cmd/devref

go 1.24.5

replace (
	github.com/orange-cloudavenue/cloudavenue-sdk-go-v2 => ../..
	github.com/orange-cloudavenue/common-go => ../../../common-go
	github.com/orange-cloudavenue/common-go/print => ../../../common-go/print
	github.com/orange-cloudavenue/common-go/validators => ../../../common-go/validators
)

require github.com/orange-cloudavenue/cloudavenue-sdk-go-v2 v0.0.0-00010101000000-000000000000

require (
	github.com/brianvoe/gofakeit/v7 v7.3.0 // indirect
	github.com/creasty/defaults v1.8.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.27.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/orange-cloudavenue/common-go/extractor v1.0.0 // indirect
	github.com/orange-cloudavenue/common-go/generator v1.2.0 // indirect
	github.com/orange-cloudavenue/common-go/internal/regex v0.0.0-20250729195615-a2902a82caeb // indirect
	github.com/orange-cloudavenue/common-go/urn v1.0.1 // indirect
	github.com/orange-cloudavenue/common-go/validators v1.0.0 // indirect
	github.com/scaleway/scaleway-sdk-go v1.0.0-beta.34 // indirect
	golang.org/x/crypto v0.40.0 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	resty.dev/v3 v3.0.0-beta.3 // indirect
)
