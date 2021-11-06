package plugin_test

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-starter-datasource-backend/pkg/plugin"
)


func TestNewSampleDatasourceWithoutJSON(t *testing.T) {
	var settings backend.DataSourceInstanceSettings
	_, err := plugin.NewSampleDatasource(settings)
	if err == nil {
		t.Error("expecting error due to missing JSON")
	}
}


func TestNewSampleDatasource(t *testing.T) {
	var settings backend.DataSourceInstanceSettings
	settings.JSONData = []byte("{}")
	_, err := plugin.NewSampleDatasource(settings)
	if err != nil {
		t.Error(err)
	}
}


func TestNewSampleDatasourceWithCreds(t *testing.T) {
	var settings backend.DataSourceInstanceSettings
	settings.JSONData = []byte("{\"authenticationProvider\": 1, \"accessKeyId\": \"test_key\"}")
	settings.DecryptedSecureJSONData = map[string]string{"secretAccessKey": "test_secret"}
	_, err := plugin.NewSampleDatasource(settings)
	if err != nil {
		t.Error(err)
	}
}


func TestNewSampleDatasourceWithEndpoint(t *testing.T) {
	var settings backend.DataSourceInstanceSettings
	settings.JSONData = []byte("{\"authenticationProvider\": 1, \"accessKeyId\": \"test_key\", \"endpoint\": \"http://localhost:9000\"}")
	_, err := plugin.NewSampleDatasource(settings)
	if err != nil {
		t.Error(err)
	}
}


func TestQueryDataWithError(t *testing.T) {
	ds := plugin.SampleDatasource{}

	resp, err := ds.QueryData(
		context.Background(),
		&backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{RefID: "A"},
			},
		},
	)
	if err != nil {
		t.Error(err)
	}

	if len(resp.Responses) != 1 {
		t.Fatal("QueryData must return a response")
	}
}


type MockS3Client struct {
}

func (client* MockS3Client) ListObjectsV2(context.Context, *s3.ListObjectsV2Input, ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	key := "some_key"
	return &s3.ListObjectsV2Output{
		Contents:              []types.Object{{
				Key: &key,
				Size: 1024,
			},
		},
	}, nil
}


func TestQueryData(t *testing.T) {
	var client s3.ListObjectsV2APIClient = &MockS3Client{}
	ds := plugin.SampleDatasource{
		Client: &client,
	}
	from := 10
	to := 19

	resp, err := ds.QueryData(
		context.Background(),
		&backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{
					RefID: "A",
					TimeRange: backend.TimeRange{
						From: time.Date(2021, time.Month(2), from, 1, 10, 0, 0, time.UTC),
						To:   time.Date(2021, time.Month(2), to, 1, 10, 0, 0, time.UTC),
					},
					Interval: 60 * 60,
					JSON: []byte("{\"Endpoint\": \"localhost\", \"Metric\": 1}"),
				},
			},
		},
	)
	if err != nil {
		t.Error(err)
	}

	if len(resp.Responses) != 1 {
		t.Fatal("QueryData must return a response")
	}

	if (to - from - 1) != resp.Responses["A"].Frames[0].Fields[1].Len() {
		t.Fatal("wrong number of values")
	}
}