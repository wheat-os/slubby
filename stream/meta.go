package stream

type Context interface {
	Err() error
	SetErr(err error)
	SetMeta(key, value any)
	GetMeta(key any) any
}

type Meta struct {
	err   error
	value map[any]any
}

func (m *Meta) Err() error {
	return m.err
}

func (m *Meta) SetErr(err error) {
	m.err = err
}

func (m *Meta) SetMeta(key, value any) {
	if m.value == nil {
		m.value = make(map[any]any)
	}
	m.value[key] = value
}

func (m *Meta) GetMeta(key any) any {
	if m.value == nil {
		return nil
	}

	return m.value[key]
}
