package service

import (
	"context"

	"github.com/vectorman1/analysis/analysis-api/db"
)

type reportsService interface {
}

type ReportsService struct {
	reportsService
	historyRepository *db.HistoryRepository
}

func (s *ReportsService) GetDailyBySymbolUuid(ctx context.Context, symbolUuid string) {
}
