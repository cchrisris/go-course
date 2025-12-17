package library

func SequentialIDGenerator(start uint64) func() uint64 {
	next := start
	return func() uint64 {
		current := next
		next++
		return current
	}
}
