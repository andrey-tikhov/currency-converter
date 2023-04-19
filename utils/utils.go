package utils

func ToPointer[T any](a T) *T {
	return &a
}
