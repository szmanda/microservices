// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"github.com/labstack/echo/v4"
)

// NipRequest defines model for NipRequest.
type NipRequest struct {
	// Nip The NIP code to check.
	Nip string `json:"nip"`
}

// NipResponse defines model for NipResponse.
type NipResponse struct {
	// Apartment Apartment number.
	Apartment *string `json:"apartment"`

	// Building Building number.
	Building *string `json:"building"`

	// City City name.
	City *string `json:"city"`

	// LongName Full name of the company.
	LongName *string `json:"longName"`

	// Province Province name.
	Province *string `json:"province"`

	// ShortName Short name of the company.
	ShortName *string `json:"shortName"`

	// Street Street name.
	Street *string `json:"street"`

	// TaxId Tax ID of the company.
	TaxId *string `json:"taxId"`

	// Zip ZIP code
	Zip *string `json:"zip"`
}

// PrintJob defines model for PrintJob.
type PrintJob struct {
	// Copies The number of copies to print.
	Copies int32 `json:"copies"`

	// DocumentName The name of the document to print.
	DocumentName string `json:"document_name"`

	// Orientation The orientation of the print job.
	Orientation string `json:"orientation"`

	// PaperSize The paper size for printing.
	PaperSize string   `json:"paper_size"`
	Printers  []string `json:"printers"`
}

// PrintStatus defines model for PrintStatus.
type PrintStatus struct {
	// Message The status of the print job.
	Message *string `json:"message,omitempty"`

	// Service Service name
	Service *string `json:"service,omitempty"`
}

// Response defines model for Response.
type Response struct {
	// Message Response message
	Message *string `json:"message,omitempty"`

	// Service Service name
	Service *string `json:"service,omitempty"`
}

// PostNipCheckerJSONRequestBody defines body for PostNipChecker for application/json ContentType.
type PostNipCheckerJSONRequestBody = NipRequest

// PostPrintJSONRequestBody defines body for PostPrint for application/json ContentType.
type PostPrintJSONRequestBody = PrintJob

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get hello message.
	// (GET /hello)
	GetHello(ctx echo.Context) error
	// Check the NIP code.
	// (POST /nip_checker)
	PostNipChecker(ctx echo.Context) error
	// Submit a print job.
	// (POST /print)
	PostPrint(ctx echo.Context) error
	// Get the status of a print job.
	// (GET /print/status)
	GetPrintStatus(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetHello converts echo context to params.
func (w *ServerInterfaceWrapper) GetHello(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetHello(ctx)
	return err
}

// PostNipChecker converts echo context to params.
func (w *ServerInterfaceWrapper) PostNipChecker(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostNipChecker(ctx)
	return err
}

// PostPrint converts echo context to params.
func (w *ServerInterfaceWrapper) PostPrint(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostPrint(ctx)
	return err
}

// GetPrintStatus converts echo context to params.
func (w *ServerInterfaceWrapper) GetPrintStatus(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetPrintStatus(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/hello", wrapper.GetHello)
	router.POST(baseURL+"/nip_checker", wrapper.PostNipChecker)
	router.POST(baseURL+"/print", wrapper.PostPrint)
	router.GET(baseURL+"/print/status", wrapper.GetPrintStatus)

}
