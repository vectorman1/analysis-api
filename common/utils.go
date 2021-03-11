package common

import (
	"fmt"

	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"

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

func GetErrorStatus(err error) error {
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			return st.Err()
		}
		switch err.(type) {
		default:
			grpclog.Infoln("error type:", err, err.Error())
		}
		return err
	}

	return err
}
