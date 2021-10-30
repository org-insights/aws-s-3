package plugin

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
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
	prefix   string   // format input
	expected []string // expected result
}{
	{"client=1000/<yyyy-MM-dd>", []string{"client=1000/", "yyyy-MM-dd"}},
	{"client=1000/<yyyy-MM-dd>/hour=<hh>", []string{"client=1000/", "yyyy-MM-dd", "/hour=", "hh"}},
	{"<yyyy-MM-dd>/client=1000/hour=<hh-mm>", []string{"yyyy-MM-dd", "/client=1000/hour=", "hh-mm"}},
}

func TestSplitPrefix(t *testing.T) {
	for _, testCase := range splitPrefixTests {
		actual := splitPrefix(testCase.prefix)
		if !reflect.DeepEqual(testCase.expected, actual) {
			t.Errorf("splitPrefix(%s): expected %s, actual %s", testCase.prefix, testCase.expected, actual)
		}
	}
}

var parseGranularityInMinutesTests = []struct {
	prefix   string // format input
	expected int    // expected result
}{
	{"client=1000/<yyyy-MM-dd>", 60 * 24}, // Day in minutes
	{"client=1000/<yyyy-MM-dd>/hour=<hh>", 60},
	{"client=1000/<yyyy-MM-dd>/hour=<HH>", 60},
	{"<yyyy-MM-dd>/client=1000/hour=<hh-mm>", 1},
	{"<yyyy-MM-dd>/client=1000/hour=<HH-mm>", 1},
}

func TestParseGranularityInMinutes(t *testing.T) {
	for _, testCase := range parseGranularityInMinutesTests {
		actual := parseGranularityInMinutes(testCase.prefix)
		if testCase.expected != actual {
			t.Errorf("parseGranularityInMinutes(%s): expected %d, actual %d", testCase.prefix, testCase.expected, actual)
		}
	}
}

var parsePrefixTests = []struct {
	prefix   string // format input
	expected string // expected result
}{
	{"client=1000/<yyyy-MM-dd>", "client=1000/2021-10-30"},
	{"client=1000/<yyyy-MM-dd>/hour=<HH>", "client=1000/2021-10-30/hour=17"},
	{"client=1000/<yyyy-MM-dd>/hour=<hh>", "client=1000/2021-10-30/hour=05"},
	{"<yyyy-MM-dd>/client=1000/hour=<HH:mm>", "2021-10-30/client=1000/hour=17:40"},
	{"<yyyy-MM-dd>/client=1000/hour=<H:mm>", "2021-10-30/client=1000/hour=17:40"},
	{"<yyyy-MM-dd>/client=1000/hour=<hh:mm>", "2021-10-30/client=1000/hour=05:40"},
}

func TestParsePrefix(t *testing.T) {
	currentTime := time.Date(2021, 10, 30, 17, 40, 0, 0, time.UTC)
	for _, testCase := range parsePrefixTests {
		actual := parsePrefix(testCase.prefix, currentTime)
		if testCase.expected != actual {
			t.Errorf("parsePrefix(%s): expected %s, actual %s", testCase.prefix, testCase.expected, actual)
		}
	}
}

type MockS3Client struct {
	error bool
}

func (client* MockS3Client) ListObjectsV2(context.Context, *s3.ListObjectsV2Input, ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	if client.error {
		return nil, errors.New("mocked failure")
	}
	key := "some_key"
	return &s3.ListObjectsV2Output{
		Contents:              []types.Object{{
				Key: &key,
				Size: 1024,
			},
		},
	}, nil
}

func TestGetPartitionInfoWithError(t *testing.T) {
	_, err := getPartitionInfo(&MockS3Client{true}, "", "")
	if err.Error() != "mocked failure" {
		t.Errorf("%s", err)
	}
}

func TestGetPartitionInfo(t *testing.T) {
	info, err := getPartitionInfo(&MockS3Client{false}, "", "")
	if err != nil {
		t.Errorf("nil error expected")
	}
	if info == nil {
		t.Errorf("expected info, got nil")
	}
}