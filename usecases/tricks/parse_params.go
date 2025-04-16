package tricks

// ParseParamsToVars parses the input params into a fixed number of variables.
// It dynamically assigns values to the provided pointers, ensuring no panic occurs.
func ParseParamsToVars(params []string, vars ...*string) {
	for i := range vars {
		if i >= len(params) { // No more params to assign
			break
		}
		if i < len(params) {
			*vars[i] = params[i]
		}
	}
}

type ParamParser[T any] struct{}

func NewParamParser[T any]() *ParamParser[T] {
	return &ParamParser[T]{}
}

func (p *ParamParser[T]) ParseParamsToVars(params []T, vars ...*T) {
	for i := range vars {
		if i >= len(params) {
			break
		}
		if i < len(params) {
			*vars[i] = params[i]
		}
	}
}
