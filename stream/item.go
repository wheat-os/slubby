package stream

type Item interface {
	Stream
	Filed()
	IName() string
}

type item struct {
	Stream
	name string
}

func (i *item) Filed() {}

func (i *item) IName() string {
	return i.name
}

func BasicItem(self Stream) Item {
	return &item{
		Stream: self,
	}
}

func BasicItemBandName(self Stream, s string) Item {
	return &item{
		Stream: self,
		name:   s,
	}
}
