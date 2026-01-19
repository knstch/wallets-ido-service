package middleware

import "github.com/go-kit/kit/endpoint"

type Middleware func(endpoint.Endpoint) endpoint.Endpoint
