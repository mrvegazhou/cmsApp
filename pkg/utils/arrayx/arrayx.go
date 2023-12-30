package arrayx

type MyType interface {
	uint | uint64 | int | int64 | string
}

func IsContain[T MyType](items []T, item T) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}
