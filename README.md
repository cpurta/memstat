# mempoint

mempoint is a simple package that allows to get system statistics for your go applications and
send those statistics to InfluxDB which can be used when monitoring your applications.

## Pre-requisites

In order to use this package you will need the InfluxDB Go client, which can be retrieved using
`go get`

Example:
```
$ go get github.com/influxdata/influxdb/client/v2
```

## Use

When using this package all that is needed is you create a new SysStat reference with the table name
you want to send system statistics to and then call the GetSysPoint function and send that to InfluxDB.

Example:

```go
package main

import (
    "fmt"
	"log"
	"testing"
	"time"

    "github.com/cpurta/mempoint"
    influxClient "github.com/influxdata/influxdb/client/v2"
)

func main() {
    client, err := influxClient.NewHTTPClient(influxClient.HTTPConfig{
		Addr:     "http://192.168.99.100:8086",
		Username: "root",
		Password: "temppwd",
		Timeout:  time.Second * 10,
	})
	if err != nil {
		log.Fatalf("Error connecting to InfluxDB: %s", err.Error())
	} else {
		fmt.Println("Succesfully connected to InfluxDB on 192.168.99.100")
	}

	// Create a new point batch
	bp, err := influxClient.NewBatchPoints(influxClient.BatchPointsConfig{
		Database:  "system_resources",
		Precision: "s",
	})
	if err != nil {
		log.Printf("Unable to create batch points for influx: %s\n", err.Error())
	}

	// send system statistics with name sys_stats
	stats := mempoint.NewSysStat("sys_stats")

	done := make(chan bool)

	limit := 50

	fmt.Println("Making some go routines...")
	for i := 0; i < limit; i++ {
		// lets make some go routines so our metrics look cool...
		go func(i int) {
			time.Sleep(time.Duration(i) * time.Second)
			done <- true
		}(i)

		time.Sleep(time.Duration(1) * time.Second)

		point := stats.GetSysPoint()

		bp.AddPoint(point)

		if i%5 == 0 {
			client.Write(bp)
		}
	}

	fmt.Println("Now let's wait for our go routines to finish...")
	// wait for all routines to finish
	for i := 0; i < limit; i++ {
		<-done
	}

	fmt.Println("Done!")
}
```

## Monitoring

With the above example we can then use an interface like (Grafana)[http://http://grafana.org/] or (Chronograf)[https://influxdata.com/time-series-platform/chronograf/]
to visualize our system resources for our go application in real-time.

In the above example I was using a Docker container using the (tutum/influxdb)[https://hub.docker.com/r/tutum/influxdb/] image for my InfluxDB which allows
you to set the user, password and any database that you want with environment variables. Then you can use Grafana, Chronograf or any other visualization project
to view the system usage metrics. I highly recommend that you do the same and play around with the main provided to make some cool looking graphs.
