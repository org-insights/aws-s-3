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

S3 data source consists of both frontend and backend components.

1. Install dependencies

   ```bash
   yarn install
   yarn build
   mage -v
   docker-compose up
   ```