package util

func PickA[A, B any](t A, _ B) A {
	return t
}
