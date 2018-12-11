package main

import (
	"log"
	"time"
	"math/rand"
	"flag"
	"io/ioutil"
	"strings"
	"regexp"
	"strconv"

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
	flag.IntVar(&nFrames, "n", 0, "Number of frames to be sent or 0 for infinity")
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
	// Read sensor data
	data := ReadDHTxx("/proc/am2301")
	if len(data) != 3 {
		return
	}

	if data[2] != "ok" {
		return
	}

	humidity, err := strconv.ParseFloat(data[0], 32)
	if err != nil {
		return
	}

	temperature, err := strconv.ParseFloat(data[1], 32)
	if err != nil {
		return
	}

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
		"client_id": getBoardSN(),
		"temp": "yes",
		"hum" : "yes",
	}

	// w/a for SW randinm generator
	rand.Seed(time.Now().UnixNano())

	fields := map[string]interface{}{
		"temp": temperature,
		"hum": humidity,
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

func ReadDHTxx(path string) []string {
	dat, _ := ioutil.ReadFile(path)
	value := strings.TrimSpace(string(dat))
	re := regexp.MustCompile(";")
	return re.Split(value, -1)
}

func getBoardSN() string {
	path := "/proc/device-tree/serial-number"
	dat, _ := ioutil.ReadFile(path)
	value := string(dat)
	value = strings.TrimSuffix(value, "\u0000")
	return value
}
