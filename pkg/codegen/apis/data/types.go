// Package data provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package data

import (
	externalRef0 "github.com/Jacobbrewer1/f1-data/pkg/codegen/apis/common"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Race defines the model for race.
type Race struct {
	Date *openapi_types.Date `json:"date,omitempty"`
	Id   *int64              `json:"id,omitempty"`
	Name *string             `json:"name,omitempty"`
}

// RaceResult defines the model for race_result.
type RaceResult struct {
	DriverName    *string `json:"driver_name,omitempty"`
	DriverNumber  *int64  `json:"driver_number,omitempty"`
	Id            *int64  `json:"id,omitempty"`
	LapsCompleted *int64  `json:"laps_completed,omitempty"`
	Points        *int64  `json:"points,omitempty"`
	Position      *string `json:"position,omitempty"`
	TeamName      *string `json:"team_name,omitempty"`
	TimeOrRetired *string `json:"time_or_retired,omitempty"`
}

// Season defines the model for season.
type Season struct {
	Id   *int64 `json:"id,omitempty"`
	Year *int64 `json:"year,omitempty"`
}

// PathRaceId defines the model for path_race_id.
type PathRaceId = int64

// PathYear defines the model for path_year.
type PathYear = int64

// QueryYear defines the model for query_year.
type QueryYear = int64

// QueryYearMax defines the model for query_year_max.
type QueryYearMax = int64

// QueryYearMin defines the model for query_year_min.
type QueryYearMin = int64

// GetRaceResultsParams defines parameters for GetRaceResults.
type GetRaceResultsParams struct {
	// Limit Report type
	Limit *externalRef0.LimitParam `form:"limit,omitempty" json:"limit,omitempty"`

	// LastVal Pagination details, last value of the sort column on the previous page.
	LastVal *externalRef0.LastValue `form:"last_val,omitempty" json:"last_val,omitempty"`

	// LastId Pagination details, last value of the id column on the previous page.
	LastId *externalRef0.LastId `form:"last_id,omitempty" json:"last_id,omitempty"`

	// SortBy Pagination details, sort column, if empty uses the id column.
	SortBy *externalRef0.SortBy `form:"sort_by,omitempty" json:"sort_by,omitempty"`

	// SortDir Pagination details, sorting order.
	SortDir *GetRaceResultsParamsSortDir `form:"sort_dir,omitempty" json:"sort_dir,omitempty"`
}

// GetRaceResultsParamsSortDir defines parameters for GetRaceResults.
type GetRaceResultsParamsSortDir string

// GetSeasonsParams defines parameters for GetSeasons.
type GetSeasonsParams struct {
	// Limit Report type
	Limit *externalRef0.LimitParam `form:"limit,omitempty" json:"limit,omitempty"`

	// LastVal Pagination details, last value of the sort column on the previous page.
	LastVal *externalRef0.LastValue `form:"last_val,omitempty" json:"last_val,omitempty"`

	// LastId Pagination details, last value of the id column on the previous page.
	LastId *externalRef0.LastId `form:"last_id,omitempty" json:"last_id,omitempty"`

	// SortBy Pagination details, sort column, if empty uses the id column.
	SortBy *externalRef0.SortBy `form:"sort_by,omitempty" json:"sort_by,omitempty"`

	// SortDir Pagination details, sorting order.
	SortDir *GetSeasonsParamsSortDir `form:"sort_dir,omitempty" json:"sort_dir,omitempty"`

	// Year The year of the season
	Year *QueryYear `form:"year,omitempty" json:"year,omitempty"`

	// YearMin The minimum year of the season
	YearMin *QueryYearMin `form:"year_min,omitempty" json:"year_min,omitempty"`

	// YearMax The maximum year of the season
	YearMax *QueryYearMax `form:"year_max,omitempty" json:"year_max,omitempty"`
}

// GetSeasonsParamsSortDir defines parameters for GetSeasons.
type GetSeasonsParamsSortDir string

// GetSeasonRacesParams defines parameters for GetSeasonRaces.
type GetSeasonRacesParams struct {
	// Limit Report type
	Limit *externalRef0.LimitParam `form:"limit,omitempty" json:"limit,omitempty"`

	// LastVal Pagination details, last value of the sort column on the previous page.
	LastVal *externalRef0.LastValue `form:"last_val,omitempty" json:"last_val,omitempty"`

	// LastId Pagination details, last value of the id column on the previous page.
	LastId *externalRef0.LastId `form:"last_id,omitempty" json:"last_id,omitempty"`

	// SortBy Pagination details, sort column, if empty uses the id column.
	SortBy *externalRef0.SortBy `form:"sort_by,omitempty" json:"sort_by,omitempty"`

	// SortDir Pagination details, sorting order.
	SortDir *GetSeasonRacesParamsSortDir `form:"sort_dir,omitempty" json:"sort_dir,omitempty"`
}

// GetSeasonRacesParamsSortDir defines parameters for GetSeasonRaces.
type GetSeasonRacesParamsSortDir string