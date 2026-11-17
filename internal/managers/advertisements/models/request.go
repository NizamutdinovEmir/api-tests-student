package models

// CreateRequest позволяет отправлять как полную форму, так и форму без обязательного поля.
type CreateRequest struct {
	Fields map[string]string
	Photo  []byte
}
