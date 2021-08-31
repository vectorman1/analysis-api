package model

import (
	"github.com/jackc/pgtype"
	"github.com/vectorman1/analysis/analysis-api/generated/instrument_service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Instrument struct {
	ID           uint        `json:"-"`
	Uuid         pgtype.UUID `json:"uuid"`
	CurrencyCode string      `json:"currency_code"`

	Isin                 string        `json:"isin"`
	Identifier           string        `json:"identifier"`
	Name                 string        `json:"name"`
	MinimumOrderQuantity pgtype.Float4 `json:"minimum_order_quantity"`
	MarketName           string        `json:"market_name"`
	MarketHoursGmt       string        `json:"market_hours_gmt"`

	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
	DeletedAt pgtype.Timestamptz `json:"deleted_at"`
}

func (Instrument) FromProtoObject(sym *instrument_service.Instrument) *Instrument {
	moq := pgtype.Float4{}
	moq.Set(sym.MinimumOrderQuantity)

	u := pgtype.UUID{}
	u.Set(sym.Uuid)

	res := &Instrument{
		Uuid:                 u,
		CurrencyCode:         sym.CurrencyCode,
		Isin:                 sym.Isin,
		Identifier:           sym.Identifier,
		Name:                 sym.Name,
		MinimumOrderQuantity: moq,
		MarketName:           sym.MarketName,
		MarketHoursGmt:       sym.MarketHoursGmt,
	}

	return res
}

func (s *Instrument) ToProto() *instrument_service.Instrument {
	// db constraint
	var u string
	s.Uuid.AssignTo(&u)

	res := &instrument_service.Instrument{
		CurrencyCode:         s.CurrencyCode,
		Isin:                 s.Isin,
		Uuid:                 u,
		Identifier:           s.Identifier,
		Name:                 s.Name,
		MinimumOrderQuantity: s.MinimumOrderQuantity.Float,
		MarketName:           s.MarketName,
		MarketHoursGmt:       s.MarketHoursGmt,
		CreatedAt:            timestamppb.New(s.CreatedAt.Time),
		UpdatedAt:            timestamppb.New(s.UpdatedAt.Time),
		DeletedAt:            timestamppb.New(s.DeletedAt.Time),
	}

	return res
}
