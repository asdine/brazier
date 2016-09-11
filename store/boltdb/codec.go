package boltdb

type rawCodec int

func (c rawCodec) Encode(v interface{}) ([]byte, error) {
	switch t := v.(type) {
	case string:
		return []byte(t), nil
	case []byte:
		return t, nil
	default:
		panic("need a string or slice of bytes")
	}
}

func (c rawCodec) Decode(b []byte, v interface{}) error {
	ptr := v.(*[]byte)
	*ptr = b
	return nil
}
