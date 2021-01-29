package webutility

// IntOrDefault ...
func IntOrDefault(v *int) int {
	if v != nil {
		return *v
	}
	return 0
}

// Int32OrDefault ...
func Int32OrDefault(v *int32) int32 {
	if v != nil {
		return *v
	}
	return 0
}

// Int64OrDefault ...
func Int64OrDefault(v *int64) int64 {
	if v != nil {
		return *v
	}
	return 0
}

// Uint32OrDefault ...
func Uint32OrDefault(v *uint32) uint32 {
	if v != nil {
		return *v
	}
	return 0
}

// Uint64OrDefault ...
func Uint64OrDefault(v *uint64) uint64 {
	if v != nil {
		return *v
	}
	return 0
}

// StringOrDefault ...
func StringOrDefault(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

// Float32OrDefault ...
func Float32OrDefault(v *float32) float32 {
	if v != nil {
		return *v
	}
	return 0.0
}

// Float64OrDefault ...
func Float64OrDefault(v *float64) float64 {
	if v != nil {
		return *v
	}
	return 0.0
}

// NilSafeInt32OrDefault ...
func NilSafeInt32OrDefault(val *interface{}) int32 {
	if *val != nil {
		return (*val).(int32)
	}
	return 0
}

// NilSafeInt64OrDefault ...
func NilSafeInt64OrDefault(val *interface{}) int64 {
	if *val != nil {
		return (*val).(int64)
	}
	return 0
}

// NilSafeUint32OrDefault ...
func NilSafeUint32OrDefault(val *interface{}) uint32 {
	if *val != nil {
		return (*val).(uint32)
	}
	return 0
}

// NilSafeUint64OrDefault ...
func NilSafeUint64OrDefault(val *interface{}) uint64 {
	if *val != nil {
		return (*val).(uint64)
	}
	return 0
}

// NilSafeFloat32OrDefault ...
func NilSafeFloat32OrDefault(val *interface{}) float32 {
	if *val != nil {
		return (*val).(float32)
	}
	return 0.0
}

// NilSafeFloat64OrDefault ...
func NilSafeFloat64OrDefault(val *interface{}) float64 {
	if *val != nil {
		return (*val).(float64)
	}
	return 0.0
}

// NilSafeStringOrDefault ...
func NilSafeStringOrDefault(val *interface{}) string {
	if *val != nil {
		return (*val).(string)
	}
	return ""
}

