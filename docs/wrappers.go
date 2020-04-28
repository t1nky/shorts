package docs

import (
	"shorts/helper"
	"shorts/models"
)

// "result": "error" and error text returns in the response
// swagger:response ResponseError
type ResponseErrorWrapper struct {
	// in: body
	Body helper.ResponseError
}

// "result": "error" and error text returns in the response
// swagger:response ResponseError
type UnauthorizedResponseWrapper struct {
	// in: body
	Body helper.ResponseError
}

// "result": "ok" returns in the response
// swagger:response ResponseOK
type ResponseOKWrapper struct {
	// in: body
	Body helper.ResponseOK
}

// Redirect path returns in the Location header
// swagger:response RedirectResponse
type RedirectResponseWrapper struct {
	// Full link
	Location string
}

// Information about a user
// swagger:response UserResponse
type UserResponseWrapper struct {
	// in: body
	Body models.UserResponse
}

// Information about a new short link
// swagger:response AddShortResponse
type AddShortResponseWrapper struct {
	// in: body
	Body models.ShortlinkResponse
}

// Information about a short link
// swagger:response ShortlinkResponse
type ShortlinkResponseWrapper struct {
	// in: body
	Body struct {
		Data   models.Shortlink `json:"data"`
		Result string           `json:"result"`
	}
}

// List of short links
// swagger:response ShortlinksResponse
type ShortlinksResponseWrapper struct {
	// in: body
	Body models.ShortlinksResponse
}

// List of top domains (up to 20)
// swagger:response TopDomainsResponse
type TopDomainsResponseWrapper struct {
	// in: body
	Body struct {
		// items.minimum: 1
		// items.maximum: 20
		Data   []models.TopDomainsResponse `json:"data"`
		Result string                      `json:"result"`
	}
}

// Information about uses in following format: "result": { "Day1": { Hour1: { Minute1: uses, Minute2: uses, ... } } }
// swagger:response ShortlinksGraphResponse
type ShortlinksGraphResponseWrapper struct {
	// in: body
	Body struct {
		// No idea how to make it generate proper example/schema with nested map types
		Data   map[string]interface{} `json:"data"`
		Result string                 `json:"result"`
	}
}

// Path parameters for deleting short link
// swagger:parameters deleteShortlink
type DeleteShortlinkParameterWrapper struct {
	// in: path
	// required: true
	ID int `json:"id"`
}

// Path parameters for getting short link
// swagger:parameters getShortlink
type GetShortlinkParameterWrapper struct {
	// in: path
	// required: true
	ID int `json:"id"`
}

// Path parameters for redirecting short link
// swagger:parameters redirectByShortlink
type RedirectShortlinkParameterWrapper struct {
	// in: path
	// required: true
	Short string `json:"short"`
}
