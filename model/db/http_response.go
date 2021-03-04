package db

type HttpResponse struct {
	Code    int
	Message string
	Details []interface{}
}
