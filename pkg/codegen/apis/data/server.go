// Package data provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package data

import (
	"fmt"
	"net/http"

	"github.com/Jacobbrewer1/uhttp"
	"github.com/gorilla/mux"
	"github.com/oapi-codegen/runtime"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get all drivers for a season
	// (GET /championships/{year}/constructors)
	GetConstructorsChampionship(w http.ResponseWriter, r *http.Request, year PathYear, params GetConstructorsChampionshipParams)
	// Get all drivers for a season
	// (GET /championships/{year}/drivers)
	GetDriversChampionship(w http.ResponseWriter, r *http.Request, year PathYear, params GetDriversChampionshipParams)
	// Get all drivers
	// (GET /drivers)
	GetDrivers(w http.ResponseWriter, r *http.Request, params GetDriversParams)
	// Get all results for a season
	// (GET /races/{race_id}/results)
	GetRaceResults(w http.ResponseWriter, r *http.Request, raceId PathRaceId, params GetRaceResultsParams)
	// Get all seasons
	// (GET /seasons)
	GetSeasons(w http.ResponseWriter, r *http.Request, params GetSeasonsParams)
	// Get all races for a season
	// (GET /seasons/{year}/races)
	GetSeasonRaces(w http.ResponseWriter, r *http.Request, year PathYear, params GetSeasonRacesParams)
}

type RateLimiterFunc = func(http.ResponseWriter, *http.Request) error
type MetricsMiddlewareFunc = http.HandlerFunc
type ErrorHandlerFunc = func(http.ResponseWriter, *http.Request, error)

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	authz             ServerInterface
	handler           ServerInterface
	rateLimiter       RateLimiterFunc
	metricsMiddleware MetricsMiddlewareFunc
	errorHandlerFunc  ErrorHandlerFunc
}

// WithAuthorization applies the passed authorization middleware to the server.
func WithAuthorization(authz ServerInterface) ServerOption {
	return func(s *ServerInterfaceWrapper) {
		s.authz = authz
	}
}

// WithRateLimiter applies the rate limiter middleware to routes with x-global-rate-limit.
func WithRateLimiter(rateLimiter RateLimiterFunc) ServerOption {
	return func(s *ServerInterfaceWrapper) {
		s.rateLimiter = rateLimiter
	}
}

// WithErrorHandlerFunc sets the error handler function for the server.
func WithErrorHandlerFunc(errorHandlerFunc ErrorHandlerFunc) ServerOption {
	return func(s *ServerInterfaceWrapper) {
		s.errorHandlerFunc = errorHandlerFunc
	}
}

// WithMetricsMiddleware applies the metrics middleware to the server.
func WithMetricsMiddleware(middleware MetricsMiddlewareFunc) ServerOption {
	return func(s *ServerInterfaceWrapper) {
		s.metricsMiddleware = middleware
	}
}

// ServerOption represents an optional feature applied to the server.
type ServerOption func(s *ServerInterfaceWrapper)

// GetConstructorsChampionship operation middleware
func (siw *ServerInterfaceWrapper) GetConstructorsChampionship(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cw := uhttp.NewResponseWriter(w,
		uhttp.WithDefaultStatusCode(http.StatusOK),
		uhttp.WithDefaultHeader("X-Request-ID", uhttp.RequestIDFromContext(ctx)),
		uhttp.WithDefaultHeader(uhttp.HeaderContentType, uhttp.ContentTypeJSON),
	)

	defer func() {
		if siw.metricsMiddleware != nil {
			siw.metricsMiddleware(cw, r)
		}
	}()

	var err error

	// ------------- Path parameter "year" -------------
	var year PathYear

	err = runtime.BindStyledParameterWithOptions("simple", "year", mux.Vars(r)["year"], &year, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "year", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetConstructorsChampionshipParams

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", r.URL.Query(), &params.Limit)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "limit", Err: err})
		return
	}

	// ------------- Optional query parameter "last_val" -------------

	err = runtime.BindQueryParameter("form", true, false, "last_val", r.URL.Query(), &params.LastVal)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "last_val", Err: err})
		return
	}

	// ------------- Optional query parameter "last_id" -------------

	err = runtime.BindQueryParameter("form", true, false, "last_id", r.URL.Query(), &params.LastId)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "last_id", Err: err})
		return
	}

	// ------------- Optional query parameter "sort_by" -------------

	err = runtime.BindQueryParameter("form", true, false, "sort_by", r.URL.Query(), &params.SortBy)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "sort_by", Err: err})
		return
	}

	// ------------- Optional query parameter "sort_dir" -------------

	err = runtime.BindQueryParameter("form", true, false, "sort_dir", r.URL.Query(), &params.SortDir)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "sort_dir", Err: err})
		return
	}

	// ------------- Optional query parameter "name" -------------

	err = runtime.BindQueryParameter("form", true, false, "name", r.URL.Query(), &params.Name)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "name", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.handler.GetConstructorsChampionship(cw, r, year, params)
	}))

	handler.ServeHTTP(cw, r.WithContext(ctx))
}

// GetDriversChampionship operation middleware
func (siw *ServerInterfaceWrapper) GetDriversChampionship(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cw := uhttp.NewResponseWriter(w,
		uhttp.WithDefaultStatusCode(http.StatusOK),
		uhttp.WithDefaultHeader("X-Request-ID", uhttp.RequestIDFromContext(ctx)),
		uhttp.WithDefaultHeader(uhttp.HeaderContentType, uhttp.ContentTypeJSON),
	)

	defer func() {
		if siw.metricsMiddleware != nil {
			siw.metricsMiddleware(cw, r)
		}
	}()

	var err error

	// ------------- Path parameter "year" -------------
	var year PathYear

	err = runtime.BindStyledParameterWithOptions("simple", "year", mux.Vars(r)["year"], &year, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "year", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetDriversChampionshipParams

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", r.URL.Query(), &params.Limit)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "limit", Err: err})
		return
	}

	// ------------- Optional query parameter "last_val" -------------

	err = runtime.BindQueryParameter("form", true, false, "last_val", r.URL.Query(), &params.LastVal)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "last_val", Err: err})
		return
	}

	// ------------- Optional query parameter "last_id" -------------

	err = runtime.BindQueryParameter("form", true, false, "last_id", r.URL.Query(), &params.LastId)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "last_id", Err: err})
		return
	}

	// ------------- Optional query parameter "sort_by" -------------

	err = runtime.BindQueryParameter("form", true, false, "sort_by", r.URL.Query(), &params.SortBy)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "sort_by", Err: err})
		return
	}

	// ------------- Optional query parameter "sort_dir" -------------

	err = runtime.BindQueryParameter("form", true, false, "sort_dir", r.URL.Query(), &params.SortDir)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "sort_dir", Err: err})
		return
	}

	// ------------- Optional query parameter "name" -------------

	err = runtime.BindQueryParameter("form", true, false, "name", r.URL.Query(), &params.Name)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "name", Err: err})
		return
	}

	// ------------- Optional query parameter "tag" -------------

	err = runtime.BindQueryParameter("form", true, false, "tag", r.URL.Query(), &params.Tag)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "tag", Err: err})
		return
	}

	// ------------- Optional query parameter "team" -------------

	err = runtime.BindQueryParameter("form", true, false, "team", r.URL.Query(), &params.Team)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "team", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.handler.GetDriversChampionship(cw, r, year, params)
	}))

	handler.ServeHTTP(cw, r.WithContext(ctx))
}

// GetDrivers operation middleware
func (siw *ServerInterfaceWrapper) GetDrivers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cw := uhttp.NewResponseWriter(w,
		uhttp.WithDefaultStatusCode(http.StatusOK),
		uhttp.WithDefaultHeader("X-Request-ID", uhttp.RequestIDFromContext(ctx)),
		uhttp.WithDefaultHeader(uhttp.HeaderContentType, uhttp.ContentTypeJSON),
	)

	defer func() {
		if siw.metricsMiddleware != nil {
			siw.metricsMiddleware(cw, r)
		}
	}()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetDriversParams

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", r.URL.Query(), &params.Limit)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "limit", Err: err})
		return
	}

	// ------------- Optional query parameter "last_val" -------------

	err = runtime.BindQueryParameter("form", true, false, "last_val", r.URL.Query(), &params.LastVal)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "last_val", Err: err})
		return
	}

	// ------------- Optional query parameter "last_id" -------------

	err = runtime.BindQueryParameter("form", true, false, "last_id", r.URL.Query(), &params.LastId)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "last_id", Err: err})
		return
	}

	// ------------- Optional query parameter "sort_by" -------------

	err = runtime.BindQueryParameter("form", true, false, "sort_by", r.URL.Query(), &params.SortBy)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "sort_by", Err: err})
		return
	}

	// ------------- Optional query parameter "sort_dir" -------------

	err = runtime.BindQueryParameter("form", true, false, "sort_dir", r.URL.Query(), &params.SortDir)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "sort_dir", Err: err})
		return
	}

	// ------------- Optional query parameter "name" -------------

	err = runtime.BindQueryParameter("form", true, false, "name", r.URL.Query(), &params.Name)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "name", Err: err})
		return
	}

	// ------------- Optional query parameter "tag" -------------

	err = runtime.BindQueryParameter("form", true, false, "tag", r.URL.Query(), &params.Tag)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "tag", Err: err})
		return
	}

	// ------------- Optional query parameter "team" -------------

	err = runtime.BindQueryParameter("form", true, false, "team", r.URL.Query(), &params.Team)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "team", Err: err})
		return
	}

	// ------------- Optional query parameter "nationality" -------------

	err = runtime.BindQueryParameter("form", true, false, "nationality", r.URL.Query(), &params.Nationality)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "nationality", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.handler.GetDrivers(cw, r, params)
	}))

	handler.ServeHTTP(cw, r.WithContext(ctx))
}

// GetRaceResults operation middleware
func (siw *ServerInterfaceWrapper) GetRaceResults(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cw := uhttp.NewResponseWriter(w,
		uhttp.WithDefaultStatusCode(http.StatusOK),
		uhttp.WithDefaultHeader("X-Request-ID", uhttp.RequestIDFromContext(ctx)),
		uhttp.WithDefaultHeader(uhttp.HeaderContentType, uhttp.ContentTypeJSON),
	)

	defer func() {
		if siw.metricsMiddleware != nil {
			siw.metricsMiddleware(cw, r)
		}
	}()

	var err error

	// ------------- Path parameter "race_id" -------------
	var raceId PathRaceId

	err = runtime.BindStyledParameterWithOptions("simple", "race_id", mux.Vars(r)["race_id"], &raceId, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "race_id", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetRaceResultsParams

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", r.URL.Query(), &params.Limit)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "limit", Err: err})
		return
	}

	// ------------- Optional query parameter "last_val" -------------

	err = runtime.BindQueryParameter("form", true, false, "last_val", r.URL.Query(), &params.LastVal)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "last_val", Err: err})
		return
	}

	// ------------- Optional query parameter "last_id" -------------

	err = runtime.BindQueryParameter("form", true, false, "last_id", r.URL.Query(), &params.LastId)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "last_id", Err: err})
		return
	}

	// ------------- Optional query parameter "sort_by" -------------

	err = runtime.BindQueryParameter("form", true, false, "sort_by", r.URL.Query(), &params.SortBy)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "sort_by", Err: err})
		return
	}

	// ------------- Optional query parameter "sort_dir" -------------

	err = runtime.BindQueryParameter("form", true, false, "sort_dir", r.URL.Query(), &params.SortDir)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "sort_dir", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.handler.GetRaceResults(cw, r, raceId, params)
	}))

	handler.ServeHTTP(cw, r.WithContext(ctx))
}

// GetSeasons operation middleware
func (siw *ServerInterfaceWrapper) GetSeasons(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cw := uhttp.NewResponseWriter(w,
		uhttp.WithDefaultStatusCode(http.StatusOK),
		uhttp.WithDefaultHeader("X-Request-ID", uhttp.RequestIDFromContext(ctx)),
		uhttp.WithDefaultHeader(uhttp.HeaderContentType, uhttp.ContentTypeJSON),
	)

	defer func() {
		if siw.metricsMiddleware != nil {
			siw.metricsMiddleware(cw, r)
		}
	}()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetSeasonsParams

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", r.URL.Query(), &params.Limit)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "limit", Err: err})
		return
	}

	// ------------- Optional query parameter "last_val" -------------

	err = runtime.BindQueryParameter("form", true, false, "last_val", r.URL.Query(), &params.LastVal)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "last_val", Err: err})
		return
	}

	// ------------- Optional query parameter "last_id" -------------

	err = runtime.BindQueryParameter("form", true, false, "last_id", r.URL.Query(), &params.LastId)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "last_id", Err: err})
		return
	}

	// ------------- Optional query parameter "sort_by" -------------

	err = runtime.BindQueryParameter("form", true, false, "sort_by", r.URL.Query(), &params.SortBy)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "sort_by", Err: err})
		return
	}

	// ------------- Optional query parameter "sort_dir" -------------

	err = runtime.BindQueryParameter("form", true, false, "sort_dir", r.URL.Query(), &params.SortDir)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "sort_dir", Err: err})
		return
	}

	// ------------- Optional query parameter "year" -------------

	err = runtime.BindQueryParameter("form", true, false, "year", r.URL.Query(), &params.Year)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "year", Err: err})
		return
	}

	// ------------- Optional query parameter "year_min" -------------

	err = runtime.BindQueryParameter("form", true, false, "year_min", r.URL.Query(), &params.YearMin)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "year_min", Err: err})
		return
	}

	// ------------- Optional query parameter "year_max" -------------

	err = runtime.BindQueryParameter("form", true, false, "year_max", r.URL.Query(), &params.YearMax)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "year_max", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.handler.GetSeasons(cw, r, params)
	}))

	handler.ServeHTTP(cw, r.WithContext(ctx))
}

// GetSeasonRaces operation middleware
func (siw *ServerInterfaceWrapper) GetSeasonRaces(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cw := uhttp.NewResponseWriter(w,
		uhttp.WithDefaultStatusCode(http.StatusOK),
		uhttp.WithDefaultHeader("X-Request-ID", uhttp.RequestIDFromContext(ctx)),
		uhttp.WithDefaultHeader(uhttp.HeaderContentType, uhttp.ContentTypeJSON),
	)

	defer func() {
		if siw.metricsMiddleware != nil {
			siw.metricsMiddleware(cw, r)
		}
	}()

	var err error

	// ------------- Path parameter "year" -------------
	var year PathYear

	err = runtime.BindStyledParameterWithOptions("simple", "year", mux.Vars(r)["year"], &year, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "year", Err: err})
		return
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetSeasonRacesParams

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", r.URL.Query(), &params.Limit)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "limit", Err: err})
		return
	}

	// ------------- Optional query parameter "last_val" -------------

	err = runtime.BindQueryParameter("form", true, false, "last_val", r.URL.Query(), &params.LastVal)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "last_val", Err: err})
		return
	}

	// ------------- Optional query parameter "last_id" -------------

	err = runtime.BindQueryParameter("form", true, false, "last_id", r.URL.Query(), &params.LastId)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "last_id", Err: err})
		return
	}

	// ------------- Optional query parameter "sort_by" -------------

	err = runtime.BindQueryParameter("form", true, false, "sort_by", r.URL.Query(), &params.SortBy)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "sort_by", Err: err})
		return
	}

	// ------------- Optional query parameter "sort_dir" -------------

	err = runtime.BindQueryParameter("form", true, false, "sort_dir", r.URL.Query(), &params.SortDir)
	if err != nil {
		siw.errorHandlerFunc(cw, r, &InvalidParamFormatError{ParamName: "sort_dir", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.handler.GetSeasonRaces(cw, r, year, params)
	}))

	handler.ServeHTTP(cw, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// wrapHandler will wrap the handler with middlewares in the other specified
// making the execution order the inverse of the parameter declaration
func wrapHandler(handler http.HandlerFunc, middlewares ...mux.MiddlewareFunc) http.Handler {
	var wrappedHandler http.Handler = handler
	for _, middleware := range middlewares {
		if middleware == nil {
			continue
		}
		wrappedHandler = middleware(wrappedHandler)
	}
	return wrappedHandler
}

// RegisterUnauthedHandlers registers any api handlers which do not have any authentication on them. Most services will not have any.
func RegisterUnauthedHandlers(router *mux.Router, si ServerInterface, opts ...ServerOption) {
	wrapper := ServerInterfaceWrapper{
		handler: si,
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(&wrapper)
	}

	router.Use(uhttp.AuthHeaderToContextMux())
	router.Use(uhttp.GenerateOrCopyRequestIDMux())

	// We do not have a gateway preparer here as no auth is sent.

	router.Methods(http.MethodGet).Path("/championships/{year}/constructors").Handler(wrapHandler(wrapper.GetConstructorsChampionship))

	router.Methods(http.MethodGet).Path("/championships/{year}/drivers").Handler(wrapHandler(wrapper.GetDriversChampionship))

	router.Methods(http.MethodGet).Path("/drivers").Handler(wrapHandler(wrapper.GetDrivers))

	router.Methods(http.MethodGet).Path("/races/{race_id}/results").Handler(wrapHandler(wrapper.GetRaceResults))

	router.Methods(http.MethodGet).Path("/seasons").Handler(wrapHandler(wrapper.GetSeasons))

	router.Methods(http.MethodGet).Path("/seasons/{year}/races").Handler(wrapHandler(wrapper.GetSeasonRaces))
}
