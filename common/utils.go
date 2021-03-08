package common

import (
	"fmt"

	"github.com/vectorman1/analysis/analysis-api/generated/proto_models"
)

func FormatOrderQuery(attr string, asc bool) string {
	var d string
	if asc {
		d = "asc"
	} else {
		d = "desc"
	}

	return fmt.Sprintf("%s %s", attr, d)
}

func ContainsSymbol(uuid string, arr []*proto_models.Symbol) (bool, *proto_models.Symbol) {
	for _, v := range arr {
		if v.Uuid == uuid {
			return true, v
		}
	}
	return false, nil
}
