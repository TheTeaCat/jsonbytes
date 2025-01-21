package jsonbytes

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

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

var longString []byte = []byte("\"" + strings.Repeat("a", 10240-2) + "\"")         // Precisely 10KiB
var longNumber []byte = []byte(strings.Repeat("1", 10240))                         // Precisely 10KiB
var longName []byte = []byte("{\"f" + strings.Repeat("o", 10240-3-5) + "\":\"\"}") // Precisely 10KiB
var longWhitespace []byte = []byte(strings.Repeat(" ", 10240-1) + "1")             // Precisely 1KiB
var longArray []byte
var longObject []byte
var manyArrays []byte
var manyObjects []byte
var manyTrues []byte
var manyFalses []byte
var manyNulls []byte
var deeplyNestedArray []byte
var deeplyNestedObject []byte
var packageLockAxios []byte

type testJsonCase struct {
	name     string
	testJson *[]byte
}

var testJsonCases []testJsonCase = []testJsonCase{
	{"LongString", &longString},
	{"LongNumber", &longNumber},
	{"LongName", &longName},
	{"LongWhitespace", &longWhitespace},
	{"LongArray", &longArray},
	{"LongObject", &longObject},
	{"ManyArrays", &manyArrays},
	{"ManyObjects", &manyObjects},
	{"ManyTrues", &manyTrues},
	{"ManyFalses", &manyFalses},
	{"ManyNulls", &manyNulls},
	{"DeeplyNestedArray", &deeplyNestedArray},
	{"DeeplyNestedObject", &deeplyNestedObject},
	{"PackageLockAxios", &packageLockAxios},
}

func init() {
	var err error

	//  5120 repeats of 1 with , separators and enclosing [] gives us a testJson of 10241 bytes (~1.0001KiB)
	longArraySlice := make([]int, 5120)
	for i := 0; i < len(longArraySlice); i++ {
		longArraySlice[i] = 1
	}
	longArray, err = json.Marshal(longArraySlice)
	if err != nil {
		panic(err)
	}

	// 1261 copies of "i":0 with i 0-1261, ',' separators and enclosing [] gives us a testJson of exactly 10KiB
	longObjectMap := map[string]int{}
	for i := 0; i < 1261; i++ {
		longObjectMap[strconv.Itoa(i)] = 0
	}
	longObject, err = json.Marshal(longObjectMap)
	if err != nil {
		panic(err)
	}

	//  3413 repeats of [] with , separators and enclosing [] gives us a testJson of exactly 10KiB
	manyArraysSlice := make([][]interface{}, 3413)
	for i := 0; i < len(manyArraysSlice); i++ {
		manyArraysSlice[i] = []interface{}{}
	}
	manyArrays, err = json.Marshal(manyArraysSlice)
	if err != nil {
		panic(err)
	}

	// 1463 copies of {"":0} with , separators and enclosing [] gives us a testJson of 10242 bytes (~1.0002KiB)
	manyObjectsSlice := make([]map[string]int, 1463)
	for i := 0; i < len(manyObjectsSlice); i++ {
		manyObjectsSlice[i] = map[string]int{"": 1}
	}
	manyObjects, err = json.Marshal(manyObjectsSlice)
	if err != nil {
		panic(err)
	}

	// 2048 repeats of true with , separators and enclosing [] gives us a testJson of 10241 bytes (~1.0001KiB)
	manyTruesSlice := make([]bool, 2048)
	for i := 0; i < len(manyTruesSlice); i++ {
		manyTruesSlice[i] = true
	}
	manyTrues, err = json.Marshal(manyTruesSlice)
	if err != nil {
		panic(err)
	}

	// 1707 repeats of false with , separators and enclosing [] gives us a testJson of 10243 bytes (~10.003KiB)
	manyFalsesSlice := make([]bool, 1707)
	manyFalses, err = json.Marshal(manyFalsesSlice)
	if err != nil {
		panic(err)
	}

	// 2048 repeats of null with , separators and enclosing [] gives us a testJson of 10241 bytes (~1.0001KiB)
	manyNullsSlice := make([]*interface{}, 2048)
	manyNulls, err = json.Marshal(manyNullsSlice)
	if err != nil {
		panic(err)
	}

	deeplyNestedArray = make([]byte, 10240)
	for i := 0; i < len(deeplyNestedArray)/2; i++ {
		deeplyNestedArray[i] = '['
		deeplyNestedArray[i+len(deeplyNestedArray)/2] = ']'
	}

	deeplyNestedObject = make([]byte, 10242)
	deeplyNestedObject[0] = '{'
	for i := 0; i < len(deeplyNestedObject)/5; i++ {
		copy(deeplyNestedObject[i*4+1:i*4+5], []byte("\"\":{"))
		deeplyNestedObject[len(deeplyNestedObject)-2-i] = '}'
	}
	deeplyNestedObject[len(deeplyNestedObject)-1] = '}'

	packageLockAxios, err = os.ReadFile("testdata/package-lock-axios.json")
	if err != nil {
		panic(err)
	}
}

func BenchmarkIsJson(b *testing.B) {
	implementations := []struct {
		name           string
		implementation func(json []byte) error
	}{
		{"JsonBytes", IsJson},
		{"EncodingJson", func(maybeJson []byte) error {
			var unmashalled interface{}
			return json.Unmarshal(maybeJson, &unmashalled)
		}},
	}
	for _, testCase := range testJsonCases {
		for _, implementation := range implementations {
			// TODO: setup encoding/json reference implementation which can handle really long numbers.
			if implementation.name == "EncodingJson" && testCase.name == "LongNumber" {
				continue
			}
			b.Run(
				testCase.name+"/"+implementation.name,
				func(b *testing.B) {
					testJsonLen := len(*testCase.testJson)
					b.ResetTimer()
					b.StopTimer()
					for n := 0; n < b.N; n++ {
						testJsonCopy := make([]byte, testJsonLen)
						copy(testJsonCopy, *testCase.testJson)
						b.StartTimer()
						err := implementation.implementation(testJsonCopy)
						b.StopTimer()
						if err != nil {
							log.Println(err.Error())
							b.FailNow()
						}
					}

				},
			)

		}
	}
}

func BenchmarkRedactAllValues(b *testing.B) {
	implementations := []struct {
		name           string
		implementation func(json []byte) ([]byte, error)
	}{
		{"JsonBytes", RedactAllValues},
		{"EncodingJson", func(v []byte) ([]byte, error) {
			var unmashalled interface{}
			err := json.Unmarshal(v, &unmashalled)
			if err != nil {
				return nil, err
			}
			return json.Marshal(referenceImplementationRedactAllValuesPremarshalled(unmashalled))
		}},
	}
	for _, testCase := range testJsonCases {
		for _, implementation := range implementations {
			// TODO: setup encoding/json reference implementation which can handle really long numbers.
			if implementation.name == "EncodingJson" && testCase.name == "LongNumber" {
				continue
			}
			b.Run(
				testCase.name+"/"+implementation.name,
				func(b *testing.B) {
					testJsonLen := len(*testCase.testJson)
					b.ResetTimer()
					b.StopTimer()
					for n := 0; n < b.N; n++ {
						testJsonCopy := make([]byte, testJsonLen)
						copy(testJsonCopy, *testCase.testJson)
						b.StartTimer()
						_, err := implementation.implementation(testJsonCopy)
						b.StopTimer()
						if err != nil {
							log.Println(err.Error())
							b.FailNow()
						}
					}
				},
			)
		}
	}
}

func BenchmarkRedactAllValuesPackageLockAxiosEncodingJsonPremarshalled(b *testing.B) {
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
