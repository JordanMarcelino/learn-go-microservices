package pageutils

func GetOffset(page, limit int64) int64 {
	return limit * (page - 1)
}
