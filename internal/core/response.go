package core

type Response[T any] struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
	Data  T      `json:"data,omitempty"`
}
