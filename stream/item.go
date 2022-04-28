package stream

type Item interface {
	Stream
}

func BasicItem(self Stream) Item {
	return self
}
