package jobs

import (
	"context"
	"time"

	"github.com/vectorman1/analysis/analysis-api/domain/instrument/service"

	"google.golang.org/grpc/grpclog"
)

type SymbolRecalculationJob struct {
	symbolService *service.InstrumentsService
}

func NewSymbolUpdateJob(symbolService *service.InstrumentsService) *SymbolRecalculationJob {
	return &SymbolRecalculationJob{symbolService: symbolService}
}

func (j SymbolRecalculationJob) Run() {
	grpclog.Infoln("[SYMBOL JOB] Starting recalculation job")
	ctx, c := context.WithTimeout(context.Background(), 30*time.Second)
	defer c()

	res, err := j.symbolService.UpdateAll(ctx)
	if err != nil {
		grpclog.Errorf("[SYMBOL JOB] Failed recalculation job: %v", err)
	}

	grpclog.Infoln("[SYMBOL JOB] Finished recalculation job: %v", res)
}
