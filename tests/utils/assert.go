package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// AssertStructFields validates that a struct has all expected fields
func AssertStructFields(t *testing.T, structPtr interface{}, expectedFields ...string) {
	t.Helper()

	v := reflect.ValueOf(structPtr)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		t.Fatalf("Expected struct, got %T", structPtr)
	}

	structType := v.Type()
	fieldMap := make(map[string]bool)

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldMap[field.Name] = true
	}

	for _, expectedField := range expectedFields {
		if !fieldMap[expectedField] {
			t.Errorf("Expected field %q not found in struct %s", expectedField, structType.Name())
		}
	}
}

// AssertNonNilFields validates that specified fields are not nil
func AssertNonNilFields(t *testing.T, structPtr interface{}, fieldNames ...string) {
	t.Helper()

	v := reflect.ValueOf(structPtr)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		t.Fatalf("Expected struct, got %T", structPtr)
	}

	for _, fieldName := range fieldNames {
		field := v.FieldByName(fieldName)
		if !field.IsValid() {
			t.Errorf("Field %q not found in struct", fieldName)
			continue
		}

		if field.Kind() == reflect.Ptr && field.IsNil() {
			t.Errorf("Field %q is nil but expected to be non-nil", fieldName)
		}
	}
}

// AssertMethodExists validates that a struct has a specific method
func AssertMethodExists(t *testing.T, obj interface{}, methodName string) {
	t.Helper()

	v := reflect.ValueOf(obj)
	method := v.MethodByName(methodName)

	if !method.IsValid() {
		t.Errorf("Method %q not found on %T", methodName, obj)
	}
}

// AssertPanicRecovery tests that a function recovers from panic gracefully
func AssertPanicRecovery(t *testing.T, fn func(), expectedMessage string) {
	t.Helper()

	defer func() {
		if r := recover(); r != nil {
			if expectedMessage != "" {
				errorMsg := fmt.Sprintf("%v", r)
				if !strings.Contains(errorMsg, expectedMessage) {
					t.Errorf("Expected panic message to contain %q, got %q", expectedMessage, errorMsg)
				}
			}
		} else if expectedMessage != "" {
			t.Error("Expected function to panic, but it didn't")
		}
	}()

	fn()
}

// AssertNoPanic tests that a function doesn't panic
func AssertNoPanic(t *testing.T, fn func()) {
	t.Helper()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Function panicked unexpectedly: %v", r)
		}
	}()

	fn()
}

// AssertJSONEqual compares two JSON strings for equality
func AssertJSONEqual(t *testing.T, expected, actual string) {
	t.Helper()

	var expectedData, actualData interface{}

	if err := json.Unmarshal([]byte(expected), &expectedData); err != nil {
		t.Fatalf("Failed to unmarshal expected JSON: %v", err)
	}

	if err := json.Unmarshal([]byte(actual), &actualData); err != nil {
		t.Fatalf("Failed to unmarshal actual JSON: %v", err)
	}

	if !reflect.DeepEqual(expectedData, actualData) {
		t.Errorf("JSON not equal:\nExpected: %s\nActual: %s", expected, actual)
	}
}

// AssertSliceEqual compares two slices for equality
func AssertSliceEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Slices not equal:\nExpected: %+v\nActual: %+v", expected, actual)
	}
}

// AssertStringContains validates that a string contains expected substrings
func AssertStringContains(t *testing.T, str string, expected ...string) {
	t.Helper()

	for _, exp := range expected {
		if !strings.Contains(str, exp) {
			t.Errorf("String %q does not contain expected substring %q", str, exp)
		}
	}
}

// AssertStringNotContains validates that a string doesn't contain forbidden substrings
func AssertStringNotContains(t *testing.T, str string, forbidden ...string) {
	t.Helper()

	for _, forb := range forbidden {
		if strings.Contains(str, forb) {
			t.Errorf("String %q contains forbidden substring %q", str, forb)
		}
	}
}

// AssertErrorType validates that an error is of a specific type
func AssertErrorType(t *testing.T, err error, expectedType interface{}) {
	t.Helper()

	if err == nil {
		t.Fatal("Expected error but got nil")
	}

	expectedTypeValue := reflect.TypeOf(expectedType)
	actualTypeValue := reflect.TypeOf(err)

	if !actualTypeValue.AssignableTo(expectedTypeValue) {
		t.Errorf("Expected error type %s, got %s", expectedTypeValue, actualTypeValue)
	}
}

// AssertErrorContains validates that an error message contains expected text
func AssertErrorContains(t *testing.T, err error, expected string) {
	t.Helper()

	if err == nil {
		t.Fatal("Expected error but got nil")
	}

	if !strings.Contains(err.Error(), expected) {
		t.Errorf("Expected error to contain %q, got %q", expected, err.Error())
	}
}

// AssertNoError validates that no error occurred
func AssertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// AssertMapContains validates that a map contains expected key-value pairs
func AssertMapContains(t *testing.T, actual map[string]interface{}, expected map[string]interface{}) {
	t.Helper()

	for key, expectedValue := range expected {
		actualValue, exists := actual[key]
		if !exists {
			t.Errorf("Expected key %q not found in map", key)
			continue
		}

		if !reflect.DeepEqual(expectedValue, actualValue) {
			t.Errorf("Key %q: expected %+v, got %+v", key, expectedValue, actualValue)
		}
	}
}

// AssertTableHeaders validates table headers against expected values
func AssertTableHeaders(t *testing.T, output string, expectedHeaders []string) {
	t.Helper()

	lines := strings.Split(output, "\n")
	if len(lines) == 0 {
		t.Fatal("No output lines found")
	}

	// Find header line (usually contains │ separators)
	headerLine := ""
	for _, line := range lines {
		if strings.Contains(line, "│") && !strings.Contains(line, "├") && !strings.Contains(line, "└") && !strings.Contains(line, "┌") {
			headerLine = line
			break
		}
	}

	if headerLine == "" {
		t.Fatal("No table header line found")
	}

	// Check that all expected headers are present
	for _, header := range expectedHeaders {
		if !strings.Contains(headerLine, header) {
			t.Errorf("Expected header %q not found in table output", header)
		}
	}
}

// AssertValidJSON validates that a string is valid JSON
func AssertValidJSON(t *testing.T, jsonStr string) {
	t.Helper()

	var result interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Errorf("Invalid JSON: %v\nJSON: %s", err, jsonStr)
	}
}

// AssertNotEmpty validates that a string/slice/map is not empty
func AssertNotEmpty(t *testing.T, value interface{}, description string) {
	t.Helper()

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		if v.Len() == 0 {
			t.Errorf("%s should not be empty", description)
		}
	case reflect.Slice, reflect.Array, reflect.Map:
		if v.Len() == 0 {
			t.Errorf("%s should not be empty", description)
		}
	case reflect.Ptr:
		if v.IsNil() {
			t.Errorf("%s should not be nil", description)
		}
	default:
		t.Errorf("AssertNotEmpty: unsupported type %T", value)
	}
}

// AssertEmpty validates that a string/slice/map is empty
func AssertEmpty(t *testing.T, value interface{}, description string) {
	t.Helper()

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		if v.Len() != 0 {
			t.Errorf("%s should be empty", description)
		}
	case reflect.Slice, reflect.Array, reflect.Map:
		if v.Len() != 0 {
			t.Errorf("%s should be empty", description)
		}
	case reflect.Ptr:
		if !v.IsNil() {
			t.Errorf("%s should be nil", description)
		}
	default:
		t.Errorf("AssertEmpty: unsupported type %T", value)
	}
}
