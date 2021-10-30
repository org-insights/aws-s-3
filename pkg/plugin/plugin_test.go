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

	resp, err := ds.QueryData(
		context.Background(),
		&backend.QueryDataRequest{
			Queries: []backend.DataQuery{
				{
					RefID: "A",
					TimeRange: backend.TimeRange{
						From: time.Date(2021, time.Month(2), 10, 1, 10, 0, 0, time.UTC),
						To:   time.Date(2021, time.Month(2), 19, 1, 10, 0, 0, time.UTC),
					},
					Interval: 60 * 60,
					JSON: []byte("{\"Endpoint\": \"localhost\"}"),
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
	println(resp.Responses["A"].Frames[0].Fields[1].Name)  // TODO: compare result
}