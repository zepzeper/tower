package data

import "net/url"

type RequestHandler interface {
	Request(method, endpoint string, params url.Values, response interface{}) error
}

