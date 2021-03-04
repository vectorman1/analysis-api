package service

import "github.com/vectorman1/analysis/analysis-api/db"

type historicalService interface {
}

type HistoricalService struct {
	historicalService
	historicalRepository *db.HistoricalRepository
}

func NewHistoricalService(historicalRepository *db.HistoricalRepository) *HistoricalService {
	return &HistoricalService{
		historicalRepository: historicalRepository,
	}
}
