package common

import (
	"fmt"
	"math/rand"

	"github.com/jackc/pgx"

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
		case pgx.PgError:
			pge, _ := err.(pgx.PgError)
			switch pge.Code {
			default:
			}
			grpclog.Infoln("error type:",
				err.(pgx.PgError).Code, err.Error())
		}
		return err
	}

	return err
}

func RandomStringWithLength(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
