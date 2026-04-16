package utils

import (
	"encoding/json"
	"testing"
)

var testFatalf = func(t testing.TB, format string, args ...any) {
	t.Fatalf(format, args...)
}

func DecodeBodyMap(t testing.TB, body []byte) map[string]any {
	t.Helper()

	var out map[string]any
	if err := json.Unmarshal(body, &out); err != nil {
		testFatalf(t, "unmarshal error: %v", err)
		return nil
	}

	return out
}

func AssertMessagePresent(t testing.TB, body map[string]any) {
	t.Helper()

	if _, ok := body["mensagem"]; !ok {
		testFatalf(t, "expected mensagem field, got %+v", body)
		return
	}
}

func AssertMessageEquals(t testing.TB, body map[string]any, expected string) {
	t.Helper()

	if body["mensagem"] != expected {
		testFatalf(t, "mensagem = %v", body["mensagem"])
		return
	}
}

func AssertDataNil(t testing.TB, body map[string]any) {
	t.Helper()

	if body["data"] != nil {
		testFatalf(t, "expected data=nil, got %v", body["data"])
		return
	}
}

func AssertDataArray(t testing.TB, body map[string]any) []interface{} {
	t.Helper()

	data, ok := body["data"].([]interface{})
	if !ok {
		testFatalf(t, "expected data array, got %T", body["data"])
		return nil
	}

	return data
}
