package storage

type Stop struct {
	ID       int    `json:"ID,omitempty"`
	StopList string `json:"stopList,omitempty"`
}

type Interface interface {
	AllList() ([]Stop, error)
	AddList(c Stop) error
	CreateStopTable() error
	DropStopTable() error
}
