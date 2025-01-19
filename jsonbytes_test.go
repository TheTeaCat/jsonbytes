package jsonbytes

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Tests

func TestIsJson(t *testing.T) {
	testCases := []string{
		// Strings
		"\"\"",
		"\"a\"",
		"\"foo\"",
		"\"\\\"\"",
		// Numbers
		"-3.14159E+123",
		"-3.14159e+123",
		"-3.14159",
		"-1",
		"-0.1",
		"-3.14159E-123",
		"-3.14159e-123",
		"-0",
		"0",
		"3.14159E-123",
		"3.14159e-123",
		"0.1",
		"1",
		"3.14159",
		"3.14159E+123",
		"3.14159e+123",
		// Booleans & null
		"true",
		"false",
		"null",
		// Various arrays
		"[]",
		"[0]",
		"[0,1,2,3,4,5,6,7,8,9]",
		// Various objects
		"{}",
		"{\"\":\"\"}",
		"{\"foo\":\"bar\"}",
		"{\"foo\":0}",
		"{\"foo\":1}",
		"{\"foo\":0.1}",
		"{\"foo\":3.14159}",
		"{\"foo\":{}}",
		"{\"foo\":{\"bar\":\"baz\"}}",
		"{\"foo\":[]}",
		"{\"foo\":[\"bar\"]}",
		"{\"foo\": [\"bar\"]}",
		"{\"foo\":true}",
		"{\"foo\":false}",
		"{\"foo\":null}",
		// Whitespace in various positions
		" 0 ",
		"  0  ",
		"   0   ",
		" -0 ",
		" 0.1 ",
		" 1 ",
		" 3.14159 ",
		" 3.14159E+123 ",
		" true ",
		" false ",
		" null ",
		" [ ] ",
		" { } ",
		" { \"foo\" : \"bar\" } ",
		" { \"foo\" : 0 } ",
		" { \"foo\" : 1 } ",
		" { \"foo\" : 0.1 } ",
		" { \"foo\" : 3.14159 } ",
		" { \"foo\" : { \"bar\" : \"baz\" } } ",
		" { \"foo\" : [ \"bar\" ] } ",
		" { \"foo\" : true } ",
		" { \"foo\" : false } ",
		" { \"foo\" : null } ",
	}
	for _, testCase := range testCases {
		t.Run(
			testCase,
			func(t *testing.T) {
				err := IsJson([]byte(testCase))
				require.Nil(t, err)
			},
		)
	}
}

func TestRedactAllValues(t *testing.T) {
	testCases := []struct {
		testJson     string
		expectedJson string
	}{
		// Strings
		{"\"\"", "\"\""},
		{"\"a\"", "\"\""},
		{"\"foo\"", "\"\""},
		{"\"\\\"\"", "\"\""},
		// Numbers
		{"-3.14159E+123", "0"},
		{"-3.14159e+123", "0"},
		{"-3.14159", "0"},
		{"-1", "0"},
		{"-0.1", "0"},
		{"-3.14159E-123", "0"},
		{"-3.14159e-123", "0"},
		{"-0", "0"},
		{"0", "0"},
		{"3.14159E-123", "0"},
		{"3.14159e-123", "0"},
		{"0.1", "0"},
		{"1", "0"},
		{"3.14159", "0"},
		{"3.14159E+123", "0"},
		{"3.14159e+123", "0"},
		// Booleans & null
		{"true", "true"},
		{"false", "true"},
		{"null", "null"},
		// Various arrays
		{"[]", "[]"},
		{"[0]", "[0]"},
		{"[0,1,2,3,4,5,6,7,8,9]", "[0,0,0,0,0,0,0,0,0,0]"},
		// Various objects
		{"{}", "{}"},
		{"{\"\":\"\"}", "{\"\":\"\"}"},
		{"{\"foo\":\"bar\"}", "{\"foo\":\"\"}"},
		{"{\"foo\":0}", "{\"foo\":0}"},
		{"{\"foo\":1}", "{\"foo\":0}"},
		{"{\"foo\":0.1}", "{\"foo\":0}"},
		{"{\"foo\":3.14159}", "{\"foo\":0}"},
		{"{\"foo\":{}}", "{\"foo\":{}}"},
		{"{\"foo\":{\"bar\":\"baz\"}}", "{\"foo\":{\"bar\":\"\"}}"},
		{"{\"foo\":[]}", "{\"foo\":[]}"},
		{"{\"foo\":[\"bar\"]}", "{\"foo\":[\"\"]}"},
		{"{\"foo\": [\"bar\"]}", "{\"foo\":[\"\"]}"},
		{"{\"foo\":true}", "{\"foo\":true}"},
		{"{\"foo\":false}", "{\"foo\":true}"},
		{"{\"foo\":null}", "{\"foo\":null}"},
		// Whitespace in various positions
		{" 0 ", "0"},
		{"  0  ", "0"},
		{"   0   ", "0"},
		{" -0 ", "0"},
		{" 0.1 ", "0"},
		{" 1 ", "0"},
		{" 3.14159 ", "0"},
		{" 3.14159E+123 ", "0"},
		{" true ", "true"},
		{" false ", "true"},
		{" null ", "null"},
		{" [ ] ", "[]"},
		{" { } ", "{}"},
		{" { \"foo\" : \"bar\" } ", "{\"foo\":\"\"}"},
		{" { \"foo\" : 0 } ", "{\"foo\":0}"},
		{" { \"foo\" : 1 } ", "{\"foo\":0}"},
		{" { \"foo\" : 0.1 } ", "{\"foo\":0}"},
		{" { \"foo\" : 3.14159 } ", "{\"foo\":0}"},
		{" { \"foo\" : { \"bar\" : \"baz\" } } ", "{\"foo\":{\"bar\":\"\"}}"},
		{" { \"foo\" : [ \"bar\" ] } ", "{\"foo\":[\"\"]}"},
		{" { \"foo\" : true } ", "{\"foo\":true}"},
		{" { \"foo\" : false } ", "{\"foo\":true}"},
		{" { \"foo\" : null } ", "{\"foo\":null}"},
	}
	for _, testCase := range testCases {
		t.Run(
			testCase.testJson,
			func(t *testing.T) {
				redactedJson, err := RedactAllValues([]byte(testCase.testJson))
				require.Nil(t, err)
				require.Equal(t, testCase.expectedJson, string(redactedJson))
			},
		)
	}
}

type invalidJsonTestCase struct {
	testJson      string
	expectedError string
}

var invalidJsonTestCases []invalidJsonTestCase = []invalidJsonTestCase{
	{"", "jsonvalidator needs more than zero bytes"},
	{" ", "read head ran out of json"},
	{"j", "expected any of \"10123456789{[tfn at index 0 but read 'j'"},
	{"\"", "read head ran out of json"},
	{"\"\x1f", "expected any codepoint except \" or \\ or control characters at index 1 but read '\x1f'"},
	{"\"\\z", "expected any of \"/\\bfnrtu at index 2 but read 'z'"},
	{"\"\\u", "expected 4 hex digits but reached end of json"},
	{"\"\\u0", "expected 4 hex digits but reached end of json"},
	{"\"\\u00", "expected 4 hex digits but reached end of json"},
	{"\"\\u000", "expected 4 hex digits but reached end of json"},
	{"\"\\u0000", "expected \" but reached end of json"},
	{"\"\\u000z", "expected 4 hex digits at index 6 but read 'z'"},
	{"\"\\u00000", "expected \" but reached end of json"},
	{"\"\\u0000z", "expected \" but reached end of json"},
	{"{", "read head ran out of json"},
	{"{\"foo\":\"bar\"", "read head ran out of json"},
	{"{\"foo\":\"bar\",", "read head ran out of json"},
	{"{\"foo", "expected \" but reached end of json"},
	{"{\"foo\" ", "read head ran out of json"},
	{"{\"foo\";", "expected : at index 6 but read ';'"},
	{"{\"foo\":", "read head ran out of json"},
	{"[", "read head ran out of json"},
	{"[,", "read head ran out of json"},
	{"[\"", "read head ran out of json"},
	{"\"f", "expected \" but reached end of json"},
	{"-", "expected any of 0123456789 but reached end of json"},
	{"0.", "expected any of 0123456789 but reached end of json"},
	{"0E", "expected + or - but reached end of json"},
	{"0E-", "read head ran out of json"},
	{"foo", "expected a at index 1 but read 'o'"},
	{"f", "expected a but reached end of json"},
	{"t", "expected r but reached end of json"},
	{"n", "expected u but reached end of json"},
	{"{\"foo\":\"bar\"} {\"foo\":\"bar\"}", "failed to consume entire json string"},
}

func TestIsJsonInvalidJsons(t *testing.T) {
	for _, testCase := range invalidJsonTestCases {
		t.Run(
			testCase.testJson,
			func(t *testing.T) {
				err := IsJson([]byte(testCase.testJson))
				require.NotNil(t, err)
				require.Equal(t, testCase.expectedError, err.Error())
			},
		)
	}
}

func TestRedactAllValuesInvalidJsons(t *testing.T) {
	for _, testCase := range invalidJsonTestCases {
		t.Run(
			testCase.testJson,
			func(t *testing.T) {
				_, err := RedactAllValues([]byte(testCase.testJson))
				require.NotNil(t, err)
				require.Equal(t, testCase.expectedError, err.Error())
			},
		)
	}
}

// Sample values for benchmarks

var longString []byte = []byte("\"" + strings.Repeat("a", 10240-2) + "\"")         // Precisely 10KiB
var longNumber []byte = []byte(strings.Repeat("1", 10240))                         // Precisely 10KiB
var longName []byte = []byte("{\"f" + strings.Repeat("o", 10240-3-5) + "\":\"\"}") // Precisely 10KiB
var longWhitespace []byte = []byte(strings.Repeat(" ", 10240-1) + "1")             // Precisely 1KiB
var longArray string
var longObject string
var manyArrays string
var manyObjects string
var manyTrues string
var manyFalses string
var manyNulls string
var packageLockAxios []byte

func init() {
	longArray = "["
	//  5119 repeats of 1 with , separators and enclosing [] gives us a testJson of 10241 bytes (~1.0001KiB)
	for i := 1; i <= 5119; i++ {
		longArray += "1"
		if i < 5119 {
			longArray += ","
		}
	}
	longArray += "]"
	longObject = "{"
	// 1261 copies of "i":0 with i 0-1261, ',' separators and enclosing [] gives us a testJson of exactly 10KiB
	for i := 1; i <= 1261; i++ {
		longObject += "\"\":0"
		if i < 1261 {
			longObject += ","
		}
	}
	longObject += "}"
	manyArrays = "["
	//  3413 repeats of 1 with , separators and enclosing [] gives us a testJson of exactly 10KiB
	for i := 1; i <= 3413; i++ {
		manyArrays += "[]"
		if i < 3413 {
			manyArrays += ","
		}
	}
	manyArrays += "]"
	manyObjects = "["
	// 1463 copies of {"":0} with , separators and enclosing [] gives us a testJson of 10242 bytes (~1.0002KiB)
	for i := 1; i <= 1463; i++ {
		manyObjects += "{\"\":0}"
		if i < 1463 {
			manyObjects += ","
		}
	}
	manyObjects += "]"
	manyTrues = "["
	// 2048 repeats of true with , separators and enclosing [] gives us a testJson of 10241 bytes (~1.0001KiB)
	for i := 1; i <= 2048; i++ {
		manyTrues += "true"
		if i < 2048 {
			manyTrues += ","
		}
	}
	manyTrues += "]"
	manyFalses = "["
	// 1707 repeats of false with , separators and enclosing [] gives us a testJson of 10243 bytes (~10.003KiB)
	for i := 1; i <= 1707; i++ {
		manyFalses += "false"
		if i < 1707 {
			manyFalses += ","
		}
	}
	manyFalses += "]"
	manyNulls = "["
	// 2048 repeats of null with , separators and enclosing [] gives us a testJson of 10241 bytes (~1.0001KiB)
	for i := 1; i <= 2048; i++ {
		manyNulls += "null"
		if i < 2048 {
			manyNulls += ","
		}
	}
	manyNulls += "]"
	var err error
	packageLockAxios, err = os.ReadFile("testdata/package-lock-axios.json")
	if err != nil {
		panic(err)
	}
}

// Benchmarks

func BenchmarkIsJsonLongString10KiB(b *testing.B) {
	benchmarkImplementationIsJson(b, IsJson, longString)
}

func BenchmarkIsJsonLongNumber10KiB(b *testing.B) {
	benchmarkImplementationIsJson(b, IsJson, longNumber)
}

func BenchmarkIsJsonLongName10KiB(b *testing.B) {
	benchmarkImplementationIsJson(b, IsJson, longName)
}

func BenchmarkIsJsonLongWhitespace10KiB(b *testing.B) {
	benchmarkImplementationIsJson(b, IsJson, longWhitespace)
}

func BenchmarkIsJsonLongArray10KiB(b *testing.B) {
	benchmarkImplementationIsJson(b, IsJson, []byte(longArray))
}

func BenchmarkIsJsonLongObject10KiB(b *testing.B) {
	benchmarkImplementationIsJson(b, IsJson, []byte(longObject))
}

func BenchmarkIsJsonManyArrays10KiB(b *testing.B) {
	benchmarkImplementationIsJson(b, IsJson, []byte(manyArrays))
}

func BenchmarkIsJsonManyObjects10KiB(b *testing.B) {
	benchmarkImplementationIsJson(b, IsJson, []byte(manyObjects))
}

func BenchmarkIsJsonManyTrues10KiB(b *testing.B) {
	benchmarkImplementationIsJson(b, IsJson, []byte(manyTrues))
}

func BenchmarkIsJsonManyFalses10KiB(b *testing.B) {

	benchmarkImplementationIsJson(b, IsJson, []byte(manyFalses))
}

func BenchmarkIsJsonManyNulls10KiB(b *testing.B) {
	benchmarkImplementationIsJson(b, IsJson, []byte(manyNulls))
}

func benchmarkImplementationIsJson(b *testing.B, implementation func(json []byte) error, testJson []byte) {
	testJsonLen := len(testJson)
	b.ResetTimer()
	b.StopTimer()
	for n := 0; n < b.N; n++ {
		testJsonCopy := make([]byte, testJsonLen)
		copy(testJsonCopy, testJson)
		b.StartTimer()
		err := implementation(testJsonCopy)
		b.StopTimer()
		if err != nil {
			log.Println(err.Error())
			b.FailNow()
		}
	}
}

func BenchmarkRedactAllValuesLongString10KiB(b *testing.B) {
	benchmarkImplementationRedactAllValues(b, RedactAllValues, longString)
}

func BenchmarkRedactAllValuesLongNumber10KiB(b *testing.B) {
	benchmarkImplementationRedactAllValues(b, RedactAllValues, longNumber)
}

func BenchmarkRedactAllValuesLongName10KiB(b *testing.B) {
	benchmarkImplementationRedactAllValues(b, RedactAllValues, longName)
}

func BenchmarkRedactAllValuesLongWhitespace10KiB(b *testing.B) {
	benchmarkImplementationRedactAllValues(b, RedactAllValues, longWhitespace)
}

func BenchmarkRedactAllValuesLongArray10KiB(b *testing.B) {
	benchmarkImplementationRedactAllValues(b, RedactAllValues, []byte(longArray))
}

func BenchmarkRedactAllValuesLongObject10KiB(b *testing.B) {
	benchmarkImplementationRedactAllValues(b, RedactAllValues, []byte(longObject))
}

func BenchmarkRedactAllValuesManyArrays10KiB(b *testing.B) {
	benchmarkImplementationRedactAllValues(b, RedactAllValues, []byte(manyArrays))
}

func BenchmarkRedactAllValuesManyObjects10KiB(b *testing.B) {
	benchmarkImplementationRedactAllValues(b, RedactAllValues, []byte(manyObjects))
}

func BenchmarkRedactAllValuesManyTrues10KiB(b *testing.B) {
	benchmarkImplementationRedactAllValues(b, RedactAllValues, []byte(manyTrues))
}

func BenchmarkRedactAllValuesManyFalses10KiB(b *testing.B) {
	benchmarkImplementationRedactAllValues(b, RedactAllValues, []byte(manyFalses))
}

func BenchmarkRedactAllValuesManyNulls10KiB(b *testing.B) {
	benchmarkImplementationRedactAllValues(b, RedactAllValues, []byte(manyNulls))
}

func BenchmarkIsJsonSamplePackageLockAxios(b *testing.B) {
	benchmarkImplementationIsJson(b, IsJson, packageLockAxios)
}

func BenchmarkIsJsonSamplePackageLockAxiosReferenceImplementation(b *testing.B) {
	benchmarkImplementationIsJson(b, referenceImplementationIsJson, packageLockAxios)
}

func BenchmarkRedactAllValuesSamplePackageLockAxios(b *testing.B) {
	benchmarkImplementationRedactAllValues(b, RedactAllValues, packageLockAxios)
}

func BenchmarkRedactAllValuesSamplePackageLockAxiosReferenceImplementation(b *testing.B) {
	benchmarkImplementationRedactAllValues(b, referenceImplementationRedactAllValues, packageLockAxios)
}

func benchmarkImplementationRedactAllValues(b *testing.B, implementation func(json []byte) ([]byte, error), testJson []byte) {
	testJsonLen := len(testJson)
	b.ResetTimer()
	b.StopTimer()
	for n := 0; n < b.N; n++ {
		testJsonCopy := make([]byte, testJsonLen)
		copy(testJsonCopy, testJson)
		b.StartTimer()
		_, err := implementation(testJsonCopy)
		b.StopTimer()
		if err != nil {
			log.Println(err.Error())
			b.FailNow()
		}
	}
}

func BenchmarkRedactAllValuesSamplePackageLockAxiosReferenceImplementationPremarshalled(b *testing.B) {
	b.ResetTimer()
	b.StopTimer()
	for n := 0; n < b.N; n++ {
		var testJsonUnmarshalled interface{}
		err := json.Unmarshal(packageLockAxios, &testJsonUnmarshalled)
		if err != nil {
			b.Error(err)
			b.FailNow()
		}
		b.StartTimer()
		referenceImplementationRedactAllValuesPremarshalled(testJsonUnmarshalled)
		b.StopTimer()
		if err != nil {
			b.Error(err)
			b.FailNow()
		}
	}
}

// Reference implementations

func referenceImplementationIsJson(v []byte) error {
	var unmashalled interface{}
	return json.Unmarshal(v, &unmashalled)
}

func referenceImplementationRedactAllValues(v []byte) ([]byte, error) {
	var unmashalled interface{}
	err := json.Unmarshal(v, &unmashalled)
	if err != nil {
		return nil, err
	}
	return json.Marshal(referenceImplementationRedactAllValuesPremarshalled(unmashalled))
}

func referenceImplementationRedactAllValuesPremarshalled(to_redact interface{}) interface{} {
	switch to_react_type_asserted := to_redact.(type) {
	case map[string]interface{}:
		for k, v := range to_react_type_asserted {
			to_react_type_asserted[k] = referenceImplementationRedactAllValuesPremarshalled(v)
		}
		return to_react_type_asserted
	case []interface{}:
		for i, v := range to_react_type_asserted {
			to_react_type_asserted[i] = referenceImplementationRedactAllValuesPremarshalled(v)
		}
		return to_react_type_asserted
	case string:
		return ""
	case float64:
		return 0
	case bool:
		return true
	default:
		return nil
	}
}
