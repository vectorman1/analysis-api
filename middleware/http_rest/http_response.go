package http_rest

type HttpResponse struct {
	Code    int
	Message string
	Details []interface{}
}
