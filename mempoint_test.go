package mempoint

import (
	"fmt"
	"log"
	"testing"
	"time"

	influxClient "github.com/influxdata/influxdb/client/v2"
)

func TestMain(m *testing.M) {
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
	stats := NewSysStat("sys_stats")

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
