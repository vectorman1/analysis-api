package jobs

import (
	"context"
	"time"

	"github.com/vectorman1/analysis/analysis-api/generated/symbol_service"
	"github.com/vectorman1/analysis/analysis-api/service"
	"google.golang.org/grpc/grpclog"
)

type HistoryUpdateJob struct {
	symbolService  *service.SymbolService
	historyService *service.HistoryService
}

func (j HistoryUpdateJob) Run() {
	grpclog.Infoln("[HISTORY JOB] Starting update job...")

	timeNow := time.Now().Weekday()
	if timeNow == time.Saturday ||
		timeNow == time.Sunday {
		grpclog.Infoln("[HISTORY JOB] Skipping job, not a working day")
	}

	res, _, err := j.symbolService.GetPaged(
		context.Background(),
		&symbol_service.ReadPagedSymbolRequest{
			Filter: &symbol_service.SymbolFilter{
				PageSize:   100000,
				PageNumber: 1,
				Order:      "identifier",
				Ascending:  false,
			},
		})
	if err != nil {
		grpclog.Errorf("[HISTORY JOB] Failed update job: %v", err)
	}

	grpclog.Infoln("[HISTORY JOB] Length of symbols to update: ", len(*res))

	hoursApprox := float32(len(*res) / 2000)
	grpclog.Infof("[HISTORY JOB] Job will take at least: %.2f hours", hoursApprox)

	for _, sym := range *res {
		if sym.MarketName == "NASDAQ" ||
			sym.MarketName == "NYSE" {

			ctx, c := context.WithTimeout(context.Background(), 2*time.Second)
			entries, err := j.historyService.UpdateSymbolHistory(ctx, sym.Uuid, sym.Identifier)
			c()
			if err != nil {
				grpclog.Errorf("[HISTORY JOB] Failed to update histories at: %s %s %s %s err: %v", sym.Isin, sym.Identifier, sym.Name, sym.MarketName, err)
				continue
			} else if entries == 0 {
				grpclog.Infof("[HISTORY JOB] No need to update: %s %s %s %s ", sym.Isin, sym.Identifier, sym.Name, sym.MarketName)
				continue
			}

			grpclog.Infof("[HISTORY JOB] Updated: %s %s %s %s Added entries: %d", sym.Isin, sym.Identifier, sym.Name, sym.MarketName, entries)

			// timeout to avoid throttle
			time.Sleep(2 * time.Second)
		} else {
			grpclog.Infof("[HISTORY JOB] Skipping: ", sym.Isin, sym.Identifier, sym.Name, sym.MarketName)
		}
	}

	grpclog.Infoln("[SYMBOL JOB] Finished update job: %v", res)
}
