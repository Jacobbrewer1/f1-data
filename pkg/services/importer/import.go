package importer

import (
	"fmt"
	"log/slog"
)

func (s *service) Import(from, to int) error {
	for i := from; i <= to; i++ {
		slog.Debug(fmt.Sprintf("Importing season %d", i))

		err := s.ImportSeason(i)
		if err != nil {
			return fmt.Errorf("error importing season %d: %w", i, err)
		}

		slog.Debug(fmt.Sprintf("Imported season %d", i))
	}

	return nil
}
