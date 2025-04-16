package tricks

import (
	"strings"
	"testing"
)

func TestParseParamsToVars(t *testing.T) {
	testCases := []struct {
		name   string
		dst    []string
		params []string
		expect []string
	}{
		{
			name:   "empty dst and params",
			dst:    []string{},
			params: []string{},
			expect: []string{},
		},
		{
			name:   "empty dst",
			dst:    []string{},
			params: []string{"param1", "param2"},
			expect: []string{},
		},
		{
			name:   "empty params",
			dst:    []string{"dst1", "dst2"},
			params: []string{},
			expect: []string{"dst1", "dst2"},
		},
		{
			name:   "dst and params with same length",
			dst:    []string{"dst1", "dst2"},
			params: []string{"param1", "param2"},
			expect: []string{"param1", "param2"},
		},
		{
			name:   "dst longer than params",
			dst:    []string{"dst1", "dst2", "dst3"},
			params: []string{"param1", "param2"},
			expect: []string{"param1", "param2", "dst3"},
		},
		{
			name:   "params longer than dst",
			dst:    []string{"dst1", "dst2"},
			params: []string{"param1", "param2", "param3"},
			expect: []string{"param1", "param2"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var dst1, dst2, dst3 string
			ParseParamsToVars(tc.params, &dst1, &dst2, &dst3)
			result := [3]string{dst1, dst2, dst3}
			if result != [3]string(tc.expect) {
				t.Errorf("expected %v, got %v", tc.expect, result)
			}
		})
	}
}

func TestParseParams_Usecase(t *testing.T) {
	{
		input := "IDXXX_SLOTXXX"
		var id, slot, name string
		ParseParamsToVars(strings.Split(input, "_"), &id, &slot, &name)
		t.Logf("ID: %s, Slot: %s, Name: %s", id, slot, name)
	}
	{
		input := "IDXXX_SLOTXXX_NAMEXXX"
		var id, slot, name string
		ParseParamsToVars(strings.Split(input, "_"), &id, &slot, &name)
		t.Logf("ID: %s, Slot: %s, Name: %s", id, slot, name)
	}
}
