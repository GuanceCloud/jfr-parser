package parser

import "fmt"

func ParseResolvableAsString(p ParseResolvable) (string, error) {
	switch v := p.(type) {
	case *RawValue:
		switch raw := v.Value.(type) {
		case string:
			return raw, nil
		case fmt.Stringer:
			return raw.String(), nil
		}
	case *String:
		return string(*v), nil
	case *Symbol:
		return v.String, nil
	case *ThreadState:
		return v.Name, nil
	case *Class:
		if v.Name != nil {
			return v.Name.String, nil
		}
	case *Method:
		if v.Name != nil {
			return v.Name.String, nil
		}
	case *Package:
		if v.Name != nil {
			return v.Name.String, nil
		}
	}
	return "", fmt.Errorf("unsupported string conversion for %T", p)
}

func ParseResolvableAsInt64(p ParseResolvable) (int64, error) {
	switch v := p.(type) {
	case *RawValue:
		switch raw := v.Value.(type) {
		case int:
			return int64(raw), nil
		case int8:
			return int64(raw), nil
		case int16:
			return int64(raw), nil
		case int32:
			return int64(raw), nil
		case int64:
			return raw, nil
		case uint8:
			return int64(raw), nil
		case uint16:
			return int64(raw), nil
		case uint32:
			return int64(raw), nil
		case uint64:
			return int64(raw), nil
		case bool:
			if raw {
				return 1, nil
			}
			return 0, nil
		}
	case *Long:
		return int64(*v), nil
	case *Int:
		return int64(*v), nil
	case *Short:
		return int64(*v), nil
	case *Byte:
		return int64(*v), nil
	case *Boolean:
		if bool(*v) {
			return 1, nil
		}
		return 0, nil
	}
	return 0, fmt.Errorf("unsupported integer conversion for %T", p)
}
