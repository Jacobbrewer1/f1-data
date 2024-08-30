package filters

// RaceResultRaceID is a filter that filters race results by race ID.
type RaceResultRaceID struct {
	raceID int
}

// NewRaceResultRaceID creates a new RaceResultRaceID.
func NewRaceResultRaceID(raceID int) *RaceResultRaceID {
	return &RaceResultRaceID{
		raceID: raceID,
	}
}

// Where returns the where clause for the filter.
func (r *RaceResultRaceID) Where() (string, []interface{}) {
	return `AND t.race_id = ?`, []interface{}{r.raceID}
}
