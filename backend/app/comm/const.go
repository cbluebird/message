package comm

type MessageType int

const (
	Text MessageType = iota
	Image
	Video
)
