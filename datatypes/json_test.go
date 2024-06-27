package main

import (
	"encoding/json"
	"testing"
)

func TestJSON(t *testing.T) {
	// 假设 datatypes.JSON 是 []byte 类型
	// jsonData := []byte(`{"name":"John Doe","age":30,"email":"john@example.com"}`)
	jsonData := []byte(`["https://example.com", "https://example.org"]`)

	// 解析到 map 中
	// var result map[string]interface{}
	var result []interface{}
	err := json.Unmarshal(jsonData, &result)
	if err != nil {
		t.Logf("Error parsing JSON: %v", err)
		return
	}

	// 遍历 map
	for key, value := range result {
		t.Logf("%s: %v\n", key, value)
	}
}
