package httpclient

import "resty.dev/v3"

func NewHTTPClient() *resty.Client {
	return resty.New()
}
