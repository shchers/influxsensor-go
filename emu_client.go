package main

import (
	"log"
	"time"
	"math/rand"

	"github.com/influxdata/influxdb/client/v2"
)

const (
	MyDB = "iwm"
)


func main() {
	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{ Addr: "http://localhost:8086" })
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  MyDB,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a point and add to batch
	tags := map[string]string{
		"client_id": "emu123",
		"temp": "yes",
		"hum" : "yes",
	}

	// w/a for SW randinm generator
	rand.Seed(time.Now().UnixNano())

	fields := map[string]interface{}{
		"temp": 21 + rand.Float64()*5 - rand.Float64()*5,
		"hum": 60 + rand.Float64()*10 - rand.Float64()*10,
	}

	pt, err := client.NewPoint("sensors", tags, fields)
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}

	// Close client resources
	if err := c.Close(); err != nil {
		log.Fatal(err)
	}
}

