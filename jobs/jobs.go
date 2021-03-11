package jobs

import (
	"github.com/bamzi/jobrunner"
	"github.com/vectorman1/analysis/analysis-api/service"
)

func ScheduleJobs(symbolService *service.SymbolsService) error {
	jobrunner.Start()
	err := jobrunner.Schedule("@every 8h", SymbolRecalculationJob{symbolService: symbolService})
	if err != nil {
		return err
	}

	return nil
}
