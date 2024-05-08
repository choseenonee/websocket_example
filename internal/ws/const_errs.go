package ws

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	RoomNotFound      = Error("error: room with this name does not exists")
	RoomAlreadyExists = Error("error: room with this name already exists")
)
