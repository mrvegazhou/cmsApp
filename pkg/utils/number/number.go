package number

func RemoveRepeatedInArr[T uint64 | int | float32 | float64](s []T) []T {
	result := make([]T, 0)
	m := make(map[T]bool) //map的值不重要
	for _, v := range s {
		if _, ok := m[v]; !ok {
			result = append(result, v)
			m[v] = true
		}
	}
	return result
}
