package events

type EventCapabilities struct {
	Raw []byte
}

func (e EventCapabilities) String() string {
	return "Capabilities"
}
