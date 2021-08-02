package client

// LateInitializeStringPtr returns in if it's non-nil, otherwise returns from
// which is the backup for the cases in is nil.
func LateInitializeStringPtr(in *string, from *string) *string {
	if in != nil {
		return in
	}
	return from
}

// LateInitializeString returns `from` if `in` is empty and `from` is non-nil,
// in other cases it returns `in`.
func LateInitializeString(in string, from *string) string {
	if in == "" && from != nil {
		return *from
	}
	return in
}

// LateInitializeIntPtr returns in if it's non-nil, otherwise returns from
// which is the backup for the cases in is nil.
func LateInitializeIntPtr(in *int, from *int64) *int {
	if in != nil {
		return in
	}
	if from != nil {
		i := int(*from)
		return &i
	}
	return nil
}

// LateInitializeInt64Ptr returns in if it's non-nil, otherwise returns from
// which is the backup for the cases in is nil.
func LateInitializeInt64Ptr(in *int64, from *int64) *int64 {
	if in != nil {
		return in
	}
	return from
}

// LateInitializeBoolPtr returns in if it's non-nil, otherwise returns from
// which is the backup for the cases in is nil.
func LateInitializeBoolPtr(in *bool, from *bool) *bool {
	if in != nil {
		return in
	}
	return from
}
