package errors

type APIError interface {
	error
	StatusCode() int
}
