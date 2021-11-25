package docs

import "fmt"

type interData struct {
	val interface{}
}

func (d interData) ToInt64() (int64, error) {
	switch t := d.val.(type) {
	case int:
		return int64(t), nil
	case int8:
		return int64(t), nil
	case int16:
		return int64(t), nil
	case int32:
		return int64(t), nil
	case int64:
		return t, nil
	case uint8:
		return int64(t), nil
	case uint16:
		return int64(t), nil
	case uint32:
		return int64(t), nil
	case uint64:
		return int64(t), nil
	default:
		return 0, fmt.Errorf("can not convert to int64, val: %v", d.val)
	}
}

func (d interData) ToFloat() (float64, error) {
	i, err := d.ToInt64()
	if err == nil {
		return float64(i), nil
	}
	switch t := d.val.(type) {
	case float32:
		return float64(t), nil
	case float64:
		return t, nil
	default:
		return 0, fmt.Errorf("can not convert to float, val: %v", d.val)
	}
}

func (d interData) ToString() string {
	return fmt.Sprintf("%v", d.val)
}
