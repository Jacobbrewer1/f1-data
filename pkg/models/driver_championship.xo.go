// Package models contains the database interaction model code
//
// GENERATED BY GOSCHEMA. DO NOT EDIT.
package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/jacobbrewer1/patcher"
	"github.com/jacobbrewer1/patcher/inserter"
	"github.com/prometheus/client_golang/prometheus"
)

// DriverChampionship represents a row from 'driver_championship'.
type DriverChampionship struct {
	Id          int       `db:"id,pk,autoinc"`
	SeasonId    int       `db:"season_id"`
	Position    int       `db:"position"`
	Driver      string    `db:"driver"`
	DriverTag   string    `db:"driver_tag"`
	Nationality string    `db:"nationality"`
	Team        string    `db:"team"`
	Points      float64   `db:"points"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// Insert inserts the DriverChampionship to the database.
func (m *DriverChampionship) Insert(db DB) error {
	t := prometheus.NewTimer(DatabaseLatency.WithLabelValues("insert_DriverChampionship"))
	defer t.ObserveDuration()

	const sqlstr = "INSERT INTO driver_championship (" +
		"`season_id`, `position`, `driver`, `driver_tag`, `nationality`, `team`, `points`, `updated_at`" +
		") VALUES (" +
		"?, ?, ?, ?, ?, ?, ?, ?" +
		")"

	DBLog(sqlstr, m.SeasonId, m.Position, m.Driver, m.DriverTag, m.Nationality, m.Team, m.Points, m.UpdatedAt)
	res, err := db.Exec(sqlstr, m.SeasonId, m.Position, m.Driver, m.DriverTag, m.Nationality, m.Team, m.Points, m.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	m.Id = int(id)
	return nil
}

func InsertManyDriverChampionships(db DB, ms ...*DriverChampionship) error {
	if len(ms) == 0 {
		return nil
	}

	t := prometheus.NewTimer(DatabaseLatency.WithLabelValues("insert_many_DriverChampionship"))
	defer t.ObserveDuration()

	vals := make([]any, 0, len(ms))
	for _, m := range ms {
		// Dereference the pointer to get the struct value.
		vals = append(vals, []any{*m})
	}

	sqlstr, args, err := inserter.NewBatch(vals, inserter.WithTable("driver_championship")).GenerateSQL()
	if err != nil {
		return fmt.Errorf("failed to create batch insert: %w", err)
	}

	DBLog(sqlstr, args...)
	res, err := db.Exec(sqlstr, args...)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	for i, m := range ms {
		m.Id = int(id + int64(i))
	}

	return nil
}

// IsPrimaryKeySet returns true if all primary key fields are set to none zero values
func (m *DriverChampionship) IsPrimaryKeySet() bool {
	return IsKeySet(m.Id)
}

// Update updates the DriverChampionship in the database.
func (m *DriverChampionship) Update(db DB) error {
	t := prometheus.NewTimer(DatabaseLatency.WithLabelValues("update_DriverChampionship"))
	defer t.ObserveDuration()

	const sqlstr = "UPDATE driver_championship " +
		"SET `season_id` = ?, `position` = ?, `driver` = ?, `driver_tag` = ?, `nationality` = ?, `team` = ?, `points` = ?, `updated_at` = ? " +
		"WHERE `id` = ?"

	DBLog(sqlstr, m.SeasonId, m.Position, m.Driver, m.DriverTag, m.Nationality, m.Team, m.Points, m.UpdatedAt, m.Id)
	res, err := db.Exec(sqlstr, m.SeasonId, m.Position, m.Driver, m.DriverTag, m.Nationality, m.Team, m.Points, m.UpdatedAt, m.Id)
	if err != nil {
		return err
	}

	// Requires clientFoundRows=true
	if i, err := res.RowsAffected(); err != nil {
		return err
	} else if i <= 0 {
		return ErrNoAffectedRows
	}

	return nil
}

func (m *DriverChampionship) Patch(db DB, newT *DriverChampionship) error {
	if newT == nil {
		return errors.New("new driver_championship is nil")
	}

	res, err := patcher.NewDiffSQLPatch(m, newT, patcher.WithTable("driver_championship"))
	if err != nil {
		return fmt.Errorf("new diff sql patch: %w", err)
	}

	sqlstr, args, err := res.GenerateSQL()
	if err != nil {
		switch {
		case errors.Is(err, patcher.ErrNoChanges):
			return nil
		default:
			return fmt.Errorf("failed to create patch: %w", err)
		}
	}

	DBLog(sqlstr, args...)
	_, err = db.Exec(sqlstr, args...)
	if err != nil {
		return fmt.Errorf("failed to execute patch: %w", err)
	}

	return nil
}

// InsertWithUpdate inserts the DriverChampionship to the database, and tries to update
// on unique constraint violations.
func (m *DriverChampionship) InsertWithUpdate(db DB) error {
	t := prometheus.NewTimer(DatabaseLatency.WithLabelValues("insert_update_DriverChampionship"))
	defer t.ObserveDuration()

	const sqlstr = "INSERT INTO driver_championship (" +
		"`season_id`, `position`, `driver`, `driver_tag`, `nationality`, `team`, `points`, `updated_at`" +
		") VALUES (" +
		"?, ?, ?, ?, ?, ?, ?, ?" +
		") ON DUPLICATE KEY UPDATE " +
		"`season_id` = VALUES(`season_id`), `position` = VALUES(`position`), `driver` = VALUES(`driver`), `driver_tag` = VALUES(`driver_tag`), `nationality` = VALUES(`nationality`), `team` = VALUES(`team`), `points` = VALUES(`points`), `updated_at` = VALUES(`updated_at`)"

	DBLog(sqlstr, m.SeasonId, m.Position, m.Driver, m.DriverTag, m.Nationality, m.Team, m.Points, m.UpdatedAt)
	res, err := db.Exec(sqlstr, m.SeasonId, m.Position, m.Driver, m.DriverTag, m.Nationality, m.Team, m.Points, m.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	m.Id = int(id)
	return nil
}

// Save saves the DriverChampionship to the database.
func (m *DriverChampionship) Save(db DB) error {
	if m.IsPrimaryKeySet() {
		return m.Update(db)
	}
	return m.Insert(db)
}

// SaveOrUpdate saves the DriverChampionship to the database, but tries to update
// on unique constraint violations.
func (m *DriverChampionship) SaveOrUpdate(db DB) error {
	if m.IsPrimaryKeySet() {
		return m.Update(db)
	}
	return m.InsertWithUpdate(db)
}

// Delete deletes the DriverChampionship from the database.
func (m *DriverChampionship) Delete(db DB) error {
	t := prometheus.NewTimer(DatabaseLatency.WithLabelValues("delete_DriverChampionship"))
	defer t.ObserveDuration()

	const sqlstr = "DELETE FROM driver_championship WHERE `id` = ?"

	DBLog(sqlstr, m.Id)
	_, err := db.Exec(sqlstr, m.Id)

	return err
}

// DriverChampionshipById retrieves a row from 'driver_championship' as a DriverChampionship.
//
// Generated from primary key.
func DriverChampionshipById(db DB, id int) (*DriverChampionship, error) {
	t := prometheus.NewTimer(DatabaseLatency.WithLabelValues("insert_DriverChampionship"))
	defer t.ObserveDuration()

	const sqlstr = "SELECT `id`, `season_id`, `position`, `driver`, `driver_tag`, `nationality`, `team`, `points`, `updated_at` " +
		"FROM driver_championship " +
		"WHERE `id` = ?"

	DBLog(sqlstr, id)
	var m DriverChampionship
	if err := db.Get(&m, sqlstr, id); err != nil {
		return nil, err
	}

	return &m, nil
}

// GetSeasonIdSeason Gets an instance of Season
//
// Generated from constraint driver_championship_season_id_fk
func (m *DriverChampionship) GetSeasonIdSeason(db DB) (*Season, error) {
	return SeasonById(db, m.SeasonId)
}
