package helper

func NilIfEmptyStr(val string) *string {
	if val == "" {
		return nil
	}
	return &val
}

func NilIfEmptyInt(val int) *int {
	if val == 0 {
		return nil
	}
	return &val
}
