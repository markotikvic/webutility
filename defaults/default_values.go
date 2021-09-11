package defaults

// Int ...
func Int(v *int) int {
	if v != nil {
		return *v
	}
	return 0
}

// Int32 ...
func Int32(v *int32) int32 {
	if v != nil {
		return *v
	}
	return 0
}

// Int64 ...
func Int64(v *int64) int64 {
	if v != nil {
		return *v
	}
	return 0
}

// Uint32 ...
func Uint32(v *uint32) uint32 {
	if v != nil {
		return *v
	}
	return 0
}

// Uint64 ...
func Uint64(v *uint64) uint64 {
	if v != nil {
		return *v
	}
	return 0
}

// String ...
func String(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

// Float32 ...
func Float32(v *float32) float32 {
	if v != nil {
		return *v
	}
	return 0.0
}

// Float64 ...
func Float64(v *float64) float64 {
	if v != nil {
		return *v
	}
	return 0.0
}

// NilSafeInt32 ...
func NilSafeInt32(val *interface{}) int32 {
	if *val != nil {
		return (*val).(int32)
	}
	return 0
}

// NilSafeInt64 ...
func NilSafeInt64(val *interface{}) int64 {
	if *val != nil {
		return (*val).(int64)
	}
	return 0
}

// NilSafeUint32 ...
func NilSafeUint32(val *interface{}) uint32 {
	if *val != nil {
		return (*val).(uint32)
	}
	return 0
}

// NilSafeUint64 ...
func NilSafeUint64(val *interface{}) uint64 {
	if *val != nil {
		return (*val).(uint64)
	}
	return 0
}

// NilSafeFloat32 ...
func NilSafeFloat32(val *interface{}) float32 {
	if *val != nil {
		return (*val).(float32)
	}
	return 0.0
}

// NilSafeFloat64 ...
func NilSafeFloat64(val *interface{}) float64 {
	if *val != nil {
		return (*val).(float64)
	}
	return 0.0
}

// NilSafeString ...
func NilSafeString(val *interface{}) string {
	if *val != nil {
		return (*val).(string)
	}
	return ""
}
