package jobs

import (
	"context"
	"time"

	"github.com/vectorman1/analysis/analysis-api/domain/instrument/service"

	"google.golang.org/grpc/grpclog"
)

type HistoryUpdateJob struct {
	historyService *service.HistoryService
}

func NewHistoryUpdateJob(historyService *service.HistoryService) *HistoryUpdateJob {
	return &HistoryUpdateJob{historyService: historyService}
}

func (j HistoryUpdateJob) Run() {
	grpclog.Infoln("[HISTORY JOB] Starting update job...")

	timeNow := time.Now()
	weekday := timeNow.Weekday()
	if weekday == time.Saturday ||
		weekday == time.Sunday {
		grpclog.Infoln("[HISTORY JOB] Skipping job, not a working day")
	}

	ctx := context.Background()
	err := j.historyService.UpdateAll(ctx)
	if err != nil {
		grpclog.Errorf("[HISTORY JOB] Failed update job: %v", err)
	}

	grpclog.Infof("[HISTORY JOB] Finished update job:\n - elapsed: %v", time.Since(timeNow))
}
