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

	addr = flag.String("a", "http://localhost:8086", "Influxdb server address")
	username = flag.String("u", "", "Influxdb server authorized user name")
	password = flag.String("p", "", "Password for Influxdb server authorization")
	dbname = flag.String("d", "test", "Influxdb database name")

	flag.Parse()

	if *password != "" && *username == "" {
		log.Fatal("Password defined without username")
	}

	conf.Addr = *addr
	conf.Username = *username
	conf.Password = *password

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

