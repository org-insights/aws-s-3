package plugin

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Make sure SampleDatasource implements required interfaces. This is important to do
// since otherwise we will only get a not implemented error response from plugin in
// runtime. In this example datasource instance implements backend.QueryDataHandler,
// backend.CheckHealthHandler, backend.StreamHandler interfaces. Plugin should not
// implement all these interfaces - only those which are required for a particular task.
// For example if plugin does not need streaming functionality then you are free to remove
// methods that implement backend.StreamHandler. Implementing instancemgmt.InstanceDisposer
// is useful to clean up resources used by previous datasource instance when a new datasource
// instance created upon datasource settings changed.
var (
	_ backend.QueryDataHandler      = (*SampleDatasource)(nil)
	_ backend.CheckHealthHandler    = (*SampleDatasource)(nil)
	_ backend.StreamHandler         = (*SampleDatasource)(nil)
	_ instancemgmt.InstanceDisposer = (*SampleDatasource)(nil)
)

type dataSourceConfig struct {
	AuthenticationProvider int    `json:"authenticationProvider"`
	AccessKeyId            string `json:"accessKeyId"`
	Endpoint               string `json:"endpoint"`
}

// NewSampleDatasource creates a new datasource instance.
func NewSampleDatasource(settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	var awsConfig aws.Config
	var credentialsProviderFunc config.LoadOptionsFunc
	var endpointResolverFunc config.LoadOptionsFunc

	var dsConfig dataSourceConfig
	err := json.Unmarshal(settings.JSONData, &dsConfig)
	if err != nil {
		log.DefaultLogger.Warn("error marshalling", "err", err)
		return nil, err
	}
	log.DefaultLogger.Info("Configurations", "authenticationProvider", dsConfig.AuthenticationProvider, "endpoint", dsConfig.Endpoint)

	credentialsProviderFunc = getCredentialsProviderFunc(dsConfig, settings.DecryptedSecureJSONData)

	if len(dsConfig.Endpoint) > 0 {
		customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:               dsConfig.Endpoint,
				HostnameImmutable: true,
			}, nil
		})

		endpointResolverFunc = config.WithEndpointResolver(customResolver)
	} else {
		endpointResolverFunc = DummyLoadOptionsFunc()
	}

	// Load the Shared AWS Configuration (~/.aws/config)
	awsConfig, err = config.LoadDefaultConfig(
		context.TODO(),
		endpointResolverFunc,
		credentialsProviderFunc,
	)
	if err != nil {
		log.DefaultLogger.Error("NewSampleDatasource called", "err", err)
	}

	log.DefaultLogger.Info("Create an Amazon S3 service client")
	// Create an Amazon S3 service client
	client := s3.NewFromConfig(awsConfig)
	log.DefaultLogger.Info("Amazon S3 service client created successfully")

	return &SampleDatasource{
		client: client,
	}, nil
}

func getCredentialsProviderFunc(dsConfig dataSourceConfig, secureData map[string]string) config.LoadOptionsFunc {
	if dsConfig.AuthenticationProvider == 1 {
		secretAccessKey, hasSecretAccessKey := secureData["secretAccessKey"]
		if hasSecretAccessKey {
			log.DefaultLogger.Info("Adding secretAccessKey for access key", "AccessKeyID", dsConfig.AccessKeyId)
		}
		return config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(dsConfig.AccessKeyId, secretAccessKey, ""))
	}
	return DummyLoadOptionsFunc()
}

func DummyLoadOptionsFunc() config.LoadOptionsFunc {
	return func(o *config.LoadOptions) error {
		return nil
	}
}

// SampleDatasource is an example datasource which can respond to data queries, reports
// its health and has streaming skills.
type SampleDatasource struct {
	client *s3.Client
}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created. As soon as datasource settings change detected by SDK old datasource instance will
// be disposed and a new one will be created using NewSampleDatasource factory function.
func (d *SampleDatasource) Dispose() {
	// Clean up datasource instance resources.
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (d *SampleDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	log.DefaultLogger.Info("QueryData called", "request", req)

	// create response struct
	response := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := d.query(ctx, req.PluginContext, q)

		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}

	return response, nil
}

type partitionInfo struct {
	Size         int64
	NumberOfKeys int64
}

func getPartitionInfo(client *s3.Client, bucket string, prefix string) (*partitionInfo, error) {
	// TODO: pagination support, currently limited to 1,000 keys per call
	var info partitionInfo
	output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		log.DefaultLogger.Error("getPartitionSize called", "err", err)
		return nil, err
	}

	for _, object := range output.Contents {
		info.Size += object.Size
		info.NumberOfKeys += 1
	}

	return &info, nil
}

type queryModel struct {
	Endpoint      string `json:"endpoint"`
	Bucket        string `json:"bucket"`
	Prefix        string `json:"prefix"`
	Metric        int    `json:"metric"`
	WithStreaming bool   `json:"withStreaming"`
}

type aggrData struct {
	Timestamp    time.Time
	Day          int
	Month        time.Month
	Year         int
	Size         int64
	NumberOfKeys int64
}

func (d *SampleDatasource) query(_ context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	response := backend.DataResponse{}

	// Unmarshal the JSON into our queryModel.
	var qm queryModel

	response.Error = json.Unmarshal(query.JSON, &qm)
	if response.Error != nil {
		return response
	}

	// create data frame response.
	frame := data.NewFrame("response")

	current := query.TimeRange.From
	// numOfFields := int(query.TimeRange.To.Sub(query.TimeRange.From).Hours() / 24)
	times := []time.Time{}
	values := []int64{}

	granularity := parseGranularityInMinutes(qm.Prefix)

	var currentDate aggrData
	currentDate.Timestamp = current
	currentDate.Day = current.Day()
	currentDate.Month = current.Month()
	currentDate.Year = current.Year()
	currentDate.Size = 0
	currentDate.NumberOfKeys = 0

	for query.TimeRange.To.After(current) {
		parsedPrefix := parsePrefix(qm.Prefix, current)
		info, err := getPartitionInfo(d.client, qm.Bucket, parsedPrefix)
		if err != nil {
			log.DefaultLogger.Error("query called", "err", err)
			response.Error = err
			return response
		}

		if current.Day() == currentDate.Day && current.Month() == currentDate.Month && current.Year() == currentDate.Year {
			currentDate.Size += info.Size
			currentDate.NumberOfKeys += info.NumberOfKeys
		} else {
			times = append(times, currentDate.Timestamp)
			if qm.Metric == 0 {
				values = append(values, currentDate.Size)
			} else {
				values = append(values, currentDate.NumberOfKeys)
			}
			currentDate.Timestamp = current
			currentDate.Day = current.Day()
			currentDate.Month = current.Month()
			currentDate.Year = current.Year()
			currentDate.Size = info.Size
			currentDate.NumberOfKeys = info.NumberOfKeys
		}

		current = current.Add(time.Duration(granularity) * time.Minute)
	}
	// TODO: add last currentDate

	// add fields.
	frame.Fields = append(frame.Fields,
		data.NewField("time", nil, times),
		data.NewField("values", nil, values),
	)

	// add the frames to the response.
	response.Frames = append(response.Frames, frame)

	return response
}

func parseGranularityInMinutes(input string) int {
	minGranularity := 60 * 24 // Day in minutes
	var oddIndex int = 1
	if strings.HasPrefix(input, "<") {
		oddIndex = 0
	}
	splited := splitPrefix(input)
	for i := 0; i < len(splited); i++ {
		if i%2 == oddIndex {
			if strings.Contains(splited[i], "mm") && minGranularity > 1 {
				minGranularity = 1
			} else if (strings.Contains(splited[i], "h") || strings.Contains(splited[i], "H")) && minGranularity > 60 {
				minGranularity = 60
			}
		}
	}
	return minGranularity
}

func parsePrefix(input string, current time.Time) string {
	var oddIndex int = 1
	if strings.HasPrefix(input, "<") {
		oddIndex = 0
	}
	splited := splitPrefix(input)
	for i := 0; i < len(splited); i++ {
		if i%2 == oddIndex {
			splited[i] = parseTime(current, splited[i])
		}
	}
	return strings.Join(splited, "")
}

func parseTime(date time.Time, format string) string {
	// example: MM/dd/yyyy HH:mm
	// golang works with
	if strings.Contains(format, "yyyy") {
		format = strings.Replace(format, "yyyy", "2006", -1)
	} else if strings.Contains(format, "yyy") {
		format = strings.Replace(format, "yyy", "006", -1)
	} else if strings.Contains(format, "yy") {
		format = strings.Replace(format, "yy", "06", -1)
	}

	if strings.Contains(format, "MM") {
		format = strings.Replace(format, "MM", "01", -1)
	} else if strings.Contains(format, "M") {
		format = strings.Replace(format, "M", "1", -1)
	}

	if strings.Contains(format, "dd") {
		format = strings.Replace(format, "dd", "02", -1)
	} else if strings.Contains(format, "d") {
		format = strings.Replace(format, "d", "2", -1)
	}

	if strings.Contains(format, "HH") {
		format = strings.Replace(format, "HH", "15", -1)
	} else if strings.Contains(format, "H") {
		format = strings.Replace(format, "H", "15", -1)
	}

	if strings.Contains(format, "hh") {
		format = strings.Replace(format, "hh", "03", -1)
	} else if strings.Contains(format, "h") {
		format = strings.Replace(format, "h", "3", -1)
	}

	if strings.Contains(format, "mm") {
		format = strings.Replace(format, "mm", "04", -1)
	} else if strings.Contains(format, "m") {
		format = strings.Replace(format, "m", "4", -1)
	}

	return date.Format(format)
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *SampleDatasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	log.DefaultLogger.Info("CheckHealth called", "request", req)

	var status = backend.HealthStatusOk
	var message = "Data source is working"

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}

func split(r rune) bool {
	return r == '<' || r == '>'
}

func splitPrefix(prefix string) []string {
	return strings.FieldsFunc(prefix, split)
}

// SubscribeStream is called when a client wants to connect to a stream. This callback
// allows sending the first message.
func (d *SampleDatasource) SubscribeStream(_ context.Context, req *backend.SubscribeStreamRequest) (*backend.SubscribeStreamResponse, error) {
	log.DefaultLogger.Info("SubscribeStream called", "request", req)

	status := backend.SubscribeStreamStatusPermissionDenied
	if req.Path == "stream" {
		// Allow subscribing only on expected path.
		status = backend.SubscribeStreamStatusOK
	}
	return &backend.SubscribeStreamResponse{
		Status: status,
	}, nil
}

// RunStream is called once for any open channel.  Results are shared with everyone
// subscribed to the same channel.
func (d *SampleDatasource) RunStream(ctx context.Context, req *backend.RunStreamRequest, sender *backend.StreamSender) error {
	log.DefaultLogger.Info("RunStream called", "request", req)

	// Create the same data frame as for query data.
	frame := data.NewFrame("response")

	// Add fields (matching the same schema used in QueryData).
	frame.Fields = append(frame.Fields,
		data.NewField("time", nil, make([]time.Time, 1)),
		data.NewField("values", nil, make([]int64, 1)),
	)

	counter := 0

	// Stream data frames periodically till stream closed by Grafana.
	for {
		select {
		case <-ctx.Done():
			log.DefaultLogger.Info("Context done, finish streaming", "path", req.Path)
			return nil
		case <-time.After(time.Second):
			// Send new data periodically.
			frame.Fields[0].Set(0, time.Now())
			frame.Fields[1].Set(0, int64(10*(counter%2+1)))

			counter++

			err := sender.SendFrame(frame, data.IncludeAll)
			if err != nil {
				log.DefaultLogger.Error("Error sending frame", "error", err)
				continue
			}
		}
	}
}

// PublishStream is called when a client sends a message to the stream.
func (d *SampleDatasource) PublishStream(_ context.Context, req *backend.PublishStreamRequest) (*backend.PublishStreamResponse, error) {
	log.DefaultLogger.Info("PublishStream called", "request", req)

	// Do not allow publishing at all.
	return &backend.PublishStreamResponse{
		Status: backend.PublishStreamStatusPermissionDenied,
	}, nil
}
