package transport

type Transport interface {
	Send() *DTO
}

type DTO interface {
}
