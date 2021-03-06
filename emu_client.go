package main

import (
	"log"
	"time"
	"math/rand"
	"flag"

	"github.com/influxdata/influxdb/client/v2"
)

func main() {
	// Influxdb configuration descriptor
	var conf client.HTTPConfig
	// Intermediate variables for scanning command line
	// * Configuration properties for db connection
	var addr, username, password *string
	// * Database properties
	var dbname *string
	// Number of frames to be sent
	var nFrames int
	// Delay between frames
	var delaySec int

	addr = flag.String("a", "http://localhost:8086", "Influxdb server address")
	username = flag.String("u", "", "Influxdb server authorized user name")
	password = flag.String("p", "", "Password for Influxdb server authorization")
	dbname = flag.String("d", "test", "Influxdb database name")
	flag.IntVar(&nFrames, "n", 1, "Number of frames to be sent or 0 for infinity")
	flag.IntVar(&delaySec, "i", 15, "Delay between frames, sec")

	flag.Parse()

	if nFrames < 0 {
		log.Fatal("Number of frames should not be less than 0")
	}

	if delaySec < 0 {
		log.Fatal("Delay between frames can not be less than 0")
	}

	if *password != "" && *username == "" {
		log.Fatal("Password defined without username")
	}

	conf.Addr = *addr
	conf.Username = *username
	conf.Password = *password

	for i := 0; nFrames == 0 || i < nFrames; i++ {
		send_data(conf, dbname)
		if nFrames > 1 {
			time.Sleep(time.Duration(delaySec) * time.Second)
		}
	}
}

func send_data(conf client.HTTPConfig, dbname *string) {
	// Create a new HTTPClient
	c, err := client.NewHTTPClient(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  *dbname,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a point and add to batch
	tags := map[string]string{
		"client_id": "emu123",
		"temp_enabled": "yes",
		"hum_enabled" : "yes",
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

