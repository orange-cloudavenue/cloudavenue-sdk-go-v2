module github.com/orange-cloudavenue/cloudavenue-sdk-go-v2

go 1.26.6

// replace (
// 	github.com/orange-cloudavenue/common-go => ../common-go
// 	github.com/orange-cloudavenue/common-go/extractor => ../common-go/extractor
// 	github.com/orange-cloudavenue/common-go/generator => ../common-go/generator
// 	github.com/orange-cloudavenue/common-go/print => ../common-go/print
// 	github.com/orange-cloudavenue/common-go/regex => ../common-go/regex
// 	github.com/orange-cloudavenue/common-go/strcase => ../common-go/strcase
// 	github.com/orange-cloudavenue/common-go/validators => ../common-go/validators
// )

require (
	github.com/aws/aws-sdk-go-v2 v1.45.1
	github.com/aws/aws-sdk-go-v2/service/s3 v1.110.0
	github.com/go-chi/chi/v5 v5.3.2
	github.com/orange-cloudavenue/common-go/extractor v1.0.1
	github.com/orange-cloudavenue/common-go/generator v1.4.0
	github.com/orange-cloudavenue/common-go/urn v1.4.0
	github.com/orange-cloudavenue/common-go/validators v1.2.0
	github.com/quasilyte/go-ruleguard/dsl v0.3.22
	github.com/stretchr/testify v1.11.0
	golang.org/x/sync v0.22.0
	resty.dev/v3 v3.0.0-rc.3
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.20 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.5.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.5.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.19 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.11.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.14.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.20.1 // indirect
	github.com/aws/smithy-go v1.28.1 // indirect
	github.com/brianvoe/gofakeit/v7 v7.15.0 // indirect
	github.com/creasty/defaults v1.8.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.15 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.30.3 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/leodido/go-urn v1.5.0 // indirect
	github.com/orange-cloudavenue/common-go/internal/regex v0.0.0-20250812202800-bd0d5bbb6c4a // indirect
	github.com/orange-cloudavenue/common-go/regex v1.2.0 // indirect
	github.com/orange-cloudavenue/common-go/strcase v1.0.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.55.0 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
