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

// Stream
