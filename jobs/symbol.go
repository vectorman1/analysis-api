package jobs

import (
	"context"
	"time"

	"google.golang.org/grpc/grpclog"

	"github.com/vectorman1/analysis/analysis-api/service"
)

type SymbolRecalculationJob struct {
	symbolService *service.SymbolsService
}

func (j SymbolRecalculationJob) Run() {
	grpclog.Infoln("[SYMBOL JOB] Starting recalculation job")
	ctx, c := context.WithTimeout(context.Background(), 10*time.Second)
	defer c()

	res, err := j.symbolService.Recalculate(ctx)
	if err != nil {
		grpclog.Errorf("[SYMBOL JOB] Failed recalculation job: %v", err)
	}

	grpclog.Infoln("[SYMBOL JOB] Finished recalculation job: %v", res)
}
