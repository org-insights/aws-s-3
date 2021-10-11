package plugin

import (
	"reflect"
	"testing"
	"time"
)

var parseTimeTests = []struct {
	t        time.Time // time input
	format   string    // format input
	expected string    // expected result
}{
	{time.Date(2021, time.Month(2), 19, 1, 10, 0, 0, time.UTC), "yyyy-MM-dd", "2021-02-19"},
	{time.Date(2021, time.Month(2), 19, 1, 10, 0, 0, time.UTC), "yyy-MM-dd", "021-02-19"},
	{time.Date(2021, time.Month(2), 19, 1, 10, 0, 0, time.UTC), "yy-MM-dd", "21-02-19"},
	{time.Date(2021, time.Month(2), 19, 1, 10, 0, 0, time.UTC), "yy-M-dd", "21-2-19"},
	{time.Date(2021, time.Month(2), 19, 1, 10, 0, 0, time.UTC), "yy-MM-d", "21-02-19"},
	{time.Date(2021, time.Month(2), 19, 1, 10, 0, 0, time.UTC), "yy-M-d", "21-2-19"},
	{time.Date(2021, time.Month(2), 5, 1, 10, 0, 0, time.UTC), "yy-M-d", "21-2-5"},
	{time.Date(2021, time.Month(2), 19, 1, 10, 0, 0, time.UTC), "yy-MM-dd/hh:mm", "21-02-19/01:10"},
	{time.Date(2021, time.Month(2), 19, 1, 10, 0, 0, time.UTC), "yy-MM-dd/hh:m", "21-02-19/01:10"},
	{time.Date(2021, time.Month(2), 19, 1, 10, 0, 0, time.UTC), "yy-MM-dd/h:m", "21-02-19/1:10"},
	{time.Date(2021, time.Month(2), 19, 1, 10, 0, 0, time.UTC), "yy-MM-dd/h:mm", "21-02-19/1:10"},
	{time.Date(2021, time.Month(2), 19, 1, 8, 0, 0, time.UTC), "yy-MM-dd/h:m", "21-02-19/1:8"},
}

func TestParseTime(t *testing.T) {
	for _, testCase := range parseTimeTests {
		actual := parseTime(testCase.t, testCase.format)
		if actual != testCase.expected {
			t.Errorf("parseTime(%s, %s): expected %s, actual %s", testCase.t, testCase.format, testCase.expected, actual)
		}
	}
}

var splitPrefixTests = []struct {
	prefix   string   	// format input
	expected []string	// expected result
}{
	{"client=1000/<yyyy-MM-dd>", []string{"client=1000/", "yyyy-MM-dd"}},
	{"client=1000/<yyyy-MM-dd>/hour=<hh>", []string{"client=1000/", "yyyy-MM-dd", "/hour=", "hh"}},
	{"<yyyy-MM-dd>/client=1000/hour=<hh>", []string{"yyyy-MM-dd", "/client=1000/hour=", "hh"}},
}


func TestSplitPrefix(t *testing.T) {
	for _, testCase := range splitPrefixTests {
		actual := splitPrefix(testCase.prefix)
		if !reflect.DeepEqual(testCase.expected, actual) {
			t.Errorf("splitPrefix(%s): expected %s, actual %s", testCase.prefix, testCase.expected, actual)
		}
	}
}
