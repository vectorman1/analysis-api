package jobs

import (
	"github.com/bamzi/jobrunner"
	"github.com/vectorman1/analysis/analysis-api/service"
)

func ScheduleJobs(symbolService *service.SymbolService, historyService *service.HistoryService) error {
	jobrunner.Start()

	err := jobrunner.Schedule("@every 2h", SymbolRecalculationJob{symbolService: symbolService})
	if err != nil {
		return err
	}

	// run every day at 22:00
	err = jobrunner.Schedule("0 22 * * *", HistoryUpdateJob{historyService: historyService})
	return nil
}
