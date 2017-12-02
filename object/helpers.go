package object

var (
	// NilObj is a predefined nil constant
	NilObj = &Nil{}

	// TrueObj is a predefined true constant
	TrueObj = &Boolean{Value: true}

	// FalseObj is a predefined false constant
	FalseObj = &Boolean{Value: false}

	// EmptyObj is a predefined empty tuple
	EmptyObj = &Tuple{Value: []Object{}}
)

// IsTruthy checks whether o is truthy.
func IsTruthy(o Object) bool {
	switch obj := o.(type) {
	case *Nil:
		return false

	case *Boolean:
		return obj.Value

	default:
		return true
	}
}

// MakeObj converts a native value to
// an Object type.
func MakeObj(v interface{}) Object {
	switch val := v.(type) {
	case nil:
		return NilObj

	case bool:
		return &Boolean{Value: val}

	case string:
		return &String{Value: val}

	case float64:
		return &Number{Value: val}

	case float32:
		return &Number{Value: float64(val)}

	case int:
		return &Number{Value: float64(val)}

	case int8:
		return &Number{Value: float64(val)}

	case int16:
		return &Number{Value: float64(val)}

	case int32:
		return &Number{Value: float64(val)}

	case int64:
		return &Number{Value: float64(val)}

	case uint:
		return &Number{Value: float64(val)}

	case uint8:
		return &Number{Value: float64(val)}

	case uint16:
		return &Number{Value: float64(val)}

	case uint32:
		return &Number{Value: float64(val)}

	case uint64:
		return &Number{Value: float64(val)}

	case []Object:
		return &List{Value: val}

	default:
		return NilObj
	}
}
