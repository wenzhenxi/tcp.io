package tcpio

const (
	_CONNECT = "connection"
	_DISCONNECT = "disconnection"
	_EMIT = "emit"
)

type packet struct {
	Type string
	Data string
}