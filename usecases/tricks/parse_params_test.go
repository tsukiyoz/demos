package tricks

import (
	"strconv"
	"strings"
	"testing"

	"github.com/samber/lo"
)

func TestNewParamParser(t *testing.T) {
	testCases := []struct {
		name   string
		dstLen int
		params []string
		expect []string
	}{
		{
			name:   "empty dst and params",
			dstLen: 0,
			params: []string{},
			expect: []string{},
		},
		{
			name:   "empty dst",
			dstLen: 0,
			params: []string{"param1", "param2"},
			expect: []string{},
		},
		{
			name:   "empty params",
			dstLen: 2,
			params: []string{},
			expect: []string{"", ""},
		},
		{
			name:   "dst and params with same length",
			dstLen: 2,
			params: []string{"param1", "param2"},
			expect: []string{"param1", "param2"},
		},
		{
			name:   "dst longer than params",
			dstLen: 3,
			params: []string{"param1", "param2"},
			expect: []string{"param1", "param2", ""},
		},
		{
			name:   "params longer than dst",
			dstLen: 2,
			params: []string{"param1", "param2", "param3"},
			expect: []string{"param1", "param2"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dst := make([]*string, tc.dstLen)
			for i := range dst {
				dst[i] = new(string) // Initialize each pointer to a new string
			}
			NewParamParser[string]().ParseParamsToVars(tc.params, dst...)
			// ParseParamsToVars(tc.params, dst...)
			for i := range tc.dstLen {
				if i < len(tc.expect) {
					if *dst[i] != tc.expect[i] {
						t.Errorf("expected %s, got %s", tc.expect[i], *dst[i])
					}
				} else {
					if *dst[i] != "" {
						t.Errorf("expected empty string, got %s", *dst[i])
					}
				}
			}
		})
	}
}

func TestParseParams_Usecase(t *testing.T) {
	{
		input := "IDXXX_SLOTXXX"
		var id, slot, name string
		NewParamParser[string]().ParseParamsToVars(strings.Split(input, "_"), &id, &slot, &name)
		// ParseParamsToVars(strings.Split(input, "_"), &id, &slot, &name)
		t.Logf("ID: %s, Slot: %s, Name: %s", id, slot, name)
	}
	{
		input := "IDXXX_SLOTXXX_NAMEXXX"
		var id, slot, name string
		id = "IDDEFAULT"
		NewParamParser[string]().ParseParamsToVars(strings.Split(input, "_"), &id, &slot, &name)
		// ParseParamsToVars(strings.Split(input, "_"), &id, &slot, &name)
		t.Logf("ID: %s, Slot: %s, Name: %s", id, slot, name)
	}
	{
		input := "123,456,789"
		var val1, val2, val3 int
		NewParamParser[int]().ParseParamsToVars(lo.Map(strings.Split(input, ","), func(s string, i int) int {
			v, _ := strconv.Atoi(s)
			return v
		}), &val1, &val2, &val3)
		t.Logf("Val1: %d, Val2: %d, Val3: %d", val1, val2, val3)
	}
}
