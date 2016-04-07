package models

const (
	Stop  = 0
	Start = 1
	Pause = 2
)

type WorckProfile struct {
	Chan   chan []byte
	Status uint8
}
