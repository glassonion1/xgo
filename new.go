package xgo

// Deprecated: Use ToPtr instead.
func New[T any](x T) *T { return &x }

// ToPtr returns a pointer copy of value.
func ToPtr[T any](x T) *T { return &x }

// FromPtr returns the pointer value or empty.
func FromPtr[T any](x *T) T {
	if x == nil {
		var zero T
		return zero
	}
	return *x
}
