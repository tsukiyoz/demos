package tricks

// ParseParamsToVars parses the input params into a fixed number of variables.
// It dynamically assigns values to the provided pointers, ensuring no panic occurs.
func ParseParamsToVars(params []string, vars ...*string) {
	for i := 0; i < len(vars); i++ {
		if i < len(params) {
			*vars[i] = params[i]
		} else {
			*vars[i] = "" // Default to empty string if params are insufficient
		}
	}
}
