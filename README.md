# Grafana S3 data source plugin

You wrote a lot of ETLs using great fancy tool like [Spark](https://spark.apache.org/), [Flink](https://flink.apache.org/) etc. A common pattern usually partitions the massive data by customer, date, hour etc. The S3 bucket's tree might looks like:
```
.
├── main.go
├── plugin.go
├── plugin_internal_test.go
├── plugin_test.go
└── client=2000
    ├── date=2021-09-25
    │   ├── hour=00
    │   │   ├── partition-00.data
    │   │   ├── partition-01.data
    │   │   ├── partition-02.data
    │   │   └── partition-03.data
    │   ├── hour=01
    │   │   ├── partition-00.data
    │   │   ├── partition-01.data
    │   │   ├── partition-02.data
    │   │   └── partition-03.data
    │   ├── hour=02
    │   │   ├── ...
    │   └── ... 
    ├── date=2021-09-26
    └── date=2021-09-27
```

How can you ensure that the production data exist?
How can you detect anomalies?
If you ever used AWS Console WEB you know how hard is to get such insights. You can struggle between sub-directiories, trying to remember sizes in your memory.

Here is where S3 AWS data-source comes in: it will give you a better observations about your partitioned data.
It will help you to ensure that there is nothing missing and detect anomalies.  


## Getting started

You can run end-to-end example with Grafana and [minio](https://min.io/) (as S3 compatible object storage): 

1. Run the following command in your terminal:

   ```bash
   docker-compose up
   ```

2. Open your browser and go to: http://localhost:3000/

## Templating
S3 Data source supports Date/Time formats such as:
- Year is represented by 2-4 y digits: `yyyy`, `yyy` or `yy`.
- Month is represented by 1-2 M digits: `MM` or `M`.
- Day is represented by 1-2 d digits: `dd` or `d`.
- Hour is represented by 1-2 h digits: `hh` or `h`.
- Minute is represented by 1-2 m digits: `mm` or `m`.

Templates should be wrapped with triangular brackets. Those will be rendered to prefixes within the selected Grafana's time range.

### Examples
**Time range:** now-30d (assuming we are on 08/31/2021)
**Prefix:** client=1000/date=<yyyy-MM-dd>
S3 Data source will list objects of the following prefixes:
```
client=1000/date=2021-08-31
client=1000/date=2021-09-01
...
client=1000/date=2021-09-29
```

**Time range:** now-10d (assuming we are on 08/31/2021 00:00)
**Prefix:** client=2000/date=<yyyy-MM-dd>/hour=<hh>
S3 Data source will list objects of the following prefixes:
```
client=2000/date=2021-08-31/hour=00
client=2000/date=2021-08-31/hour=01
..
client=2000/date=2021-08-31/hour=23
client=2000/date=2021-09-01/hour=00
client=2000/date=2021-09-01/hour=01
...
client=2000/date=2021-09-09/hour=23
```

## Screenshots

- **Data source**: Overview of data source configurations.

  ![Data source](https://raw.githubusercontent.com/org-insights/aws-s-3/master/src/screenshots/datasource-configurations.png)

- **Dashboard example**: Overview of graph panel configurations.

  ![Data source](https://raw.githubusercontent.com/org-insights/aws-s-3/master/src/screenshots/dashboard-example-1.png)
