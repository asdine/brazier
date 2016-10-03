package boltdb

type rawCodec int

func (c rawCodec) Marshal(v interface{}) ([]byte, error) {
	switch t := v.(type) {
	case string:
		return []byte(t), nil
	case []byte:
		return t, nil
	default:
		panic("need a string or slice of bytes")
	}
}

func (c rawCodec) Unmarshal(b []byte, v interface{}) error {
	ptr, ok := v.(*[]byte)
	if !ok {
		return nil
	}
	*ptr = b
	return nil
}

func (c rawCodec) Name() string {
	return "raw"
}
