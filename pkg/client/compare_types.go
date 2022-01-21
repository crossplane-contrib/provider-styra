package client

// IsEqualString determines if StringValue's of two pointers are equal.
func IsEqualString(s1 *string, s2 *string) bool {
	return StringValue(s1) == StringValue(s2)
}

// IsEqualStringArrayContent determines if two string arrays contain the
// same elements, regardless of order.
func IsEqualStringArrayContent(a1 []string, a2 []string) bool {
	if len(a1) != len(a2) {
		return false
	}

	map2 := make(map[string]struct{}, len(a2))
	for _, a2v := range a2 {
		map2[a2v] = struct{}{}
	}

	for _, a1v := range a1 {
		if _, exists := map2[a1v]; !exists {
			return false
		}
	}

	return true
}

// IsEqualBool determines if BoolValue's of two pointers are equal.
func IsEqualBool(b1 *bool, b2 *bool) bool {
	return BoolValue(b1) == BoolValue(b2)
}

// IsEqualInt64 determines if Int64Value's of two pointers are equal.
func IsEqualInt64(i1 *int64, i2 *int64) bool {
	return Int64Value(i1) == Int64Value(i2)
}
