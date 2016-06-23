package mempoint

import (
	"log"
	"os"
	"runtime"
	"time"

	influxClient "github.com/influxdata/influxdb/client/v2"
)

// SysStat is a holder for the influx database table to which we send system stats to
type SysStat struct {
	Name string
}

// NewSysStat will create a MemStat reference so that we can can call the GetMemPoint function
// and create a new point to send to InfluxDB for monitoring
func NewSysStat(name string) *SysStat {
	return &SysStat{Name: name}
}

// GetSysPoint will create a new InfluxDB data point with system stats which can be called whenever
// you want to send system stats to InfluxDB.
func (stat *SysStat) GetSysPoint() *influxClient.Point {
	// Create system stat point
	host, _ := os.Hostname()

	tags := map[string]string{"host": host}
	fields := map[string]interface{}{}
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// gather the system stats that we want to gather
	fields["num_go_routines"] = runtime.NumGoroutine()
	fields["memory_total_alloc"] = int64(m.TotalAlloc)
	fields["memory_total_heap_alloc"] = int64(m.HeapAlloc)
	fields["memory_total_heap_in_use"] = int64(m.HeapAlloc)
	fields["memory_total_sys"] = int64(m.Sys)

	// create a new influx point and return that point
	pt, err := influxClient.NewPoint(stat.Name, tags, fields, time.Now())
	if err != nil {
		log.Printf("Unable to create new memory point: %s\n", err.Error())
		return nil
	}

	return pt
}
