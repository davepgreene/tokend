package status

type Status uint8

const (
    Ready Status = iota + 1
    Pending
    Error
)
