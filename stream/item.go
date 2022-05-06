package stream

type Item interface {
	Stream
	Filed()
}

type item struct {
	Stream
}

func (i *item) Filed() {}

func BasicItem(self Stream) Item {
	return &item{
		Stream: self,
	}
}
