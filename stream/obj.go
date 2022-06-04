package stream

// 组合对象

type StreamList struct {
	Stream
	streams []Stream
}

// 迭代器模式
func (s *StreamList) Iterator() func() Stream {
	i := 0
	return func() Stream {
		if i < len(s.streams) {
			i += 1
			return s.streams[i-1]
		}

		return nil
	}

}

func StreamLists(self Stream, streams ...Stream) Stream {
	return &StreamList{
		Stream:  self,
		streams: streams,
	}
}

func StreamListRangeInt(self Stream, only func(i int) (Stream, error), start, end int) (Stream, error) {
	lists := make([]Stream, 0, end-start+1)
	for i := start; i <= end; i++ {
		stream, err := only(i)
		if err != nil {
			return nil, err
		}

		lists = append(lists, stream)
	}

	return &StreamList{
		Stream:  self,
		streams: lists,
	}, nil
}

func StreamListRangeString(self Stream, only func(value string) (Stream, error), s []string) (Stream, error) {
	lists := make([]Stream, 0, len(s))

	for _, val := range s {
		stm, err := only(val)
		if err != nil {
			return nil, err
		}

		lists = append(lists, stm)
	}

	return &StreamList{
		Stream:  self,
		streams: lists,
	}, nil
}

func StreamListRangeFloat(self Stream, only func(value float64) (Stream, error), f []float64) (Stream, error) {
	lists := make([]Stream, 0, len(f))

	for _, val := range f {
		stm, err := only(val)
		if err != nil {
			return nil, err
		}

		lists = append(lists, stm)
	}

	return &StreamList{
		Stream:  self,
		streams: lists,
	}, nil
}
