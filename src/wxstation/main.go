// main.go
//
package main

import (
	"encoding/json"
	"fmt"
	"github.com/ziutek/rrd"
	"log"
	"net/http"
	"os"
	"probe"
	"time"
)

const (
	SensorTickRate = 5 // anticipated update rate in seconds.
	dsTemperature  = "temperature"
	dsHumidity     = "humidity"
)

var (
	tempUpdated  time.Time
	temperature  float64 = 0.0
	humidUpdated time.Time
	humidity     float64 = 0.0
	doorUpdated  time.Time
	door         bool

	rrdFile    string = "data.rrd"
	rrdUpdater *rrd.Updater
)

func setupRRD() {
	c := rrd.NewCreator(rrdFile, time.Now(), SensorTickRate)

	ticksPerMinute := 60 / SensorTickRate
	// average per minute...  keep for a day
	c.RRA("AVERAGE", 0.3, ticksPerMinute, 24*60)
	// average per 5 minutes... keep for a week
	c.RRA("AVERAGE", 0.3, 5*ticksPerMinute, 7*24*60/5)
	// average per 30 minutes... keep for a month
	c.RRA("AVERAGE", 0.3, 15*ticksPerMinute, 31*24*4)
	// average per hour.... keep for a year
	c.RRA("AVERAGE", 0.3, 60*ticksPerMinute, 365*24)
	//FIXME: these limits are based on the DHT22 capabilities - may need to be changed for better sensor modules
	c.DS(dsTemperature, "GAUGE", SensorTickRate*3, -40, 125)
	c.DS(dsHumidity, "GAUGE", SensorTickRate*3, 0, 100)
	err := c.Create(false)
	if err != nil {
		if !os.IsExist(err) {
			log.Fatal(err)
		}
	}
	rrdUpdater = rrd.NewUpdater(rrdFile)
}

func main() {
	p, err := probe.Open("/dev/ttyACM0")
	if err != nil {
		fmt.Printf("Failed: %s\n", err)
		os.Exit(1)
	}
	setupRRD()
	p.Start()
	go maintainStatus(p)
	defer p.Stop()
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/g/temp", tempDayGraphHandler)
	http.HandleFunc("/g/humidity", humidityDayGraphHandler)
	log.Fatal(http.ListenAndServe(":9998", nil))
}

func maintainStatus(p *probe.Probe) {
	for {
		update, ok := <-p.Updates
		if !ok {
			break
		}
		var updParams [3]interface{}
		updParams[0] = time.Now()
		if update.Temperature != nil {
			tempUpdated = update.When
			temperature = *update.Temperature
			updParams[1] = temperature
		}
		if update.Humidity != nil {
			humidUpdated = update.When
			humidity = *update.Humidity
			updParams[2] = humidity
		}
		doorUpdated = update.When
		door = update.DoorStatus
		rrdUpdater.Update(updParams[:]...)
	}
}

type JSONResponse struct {
	Temperature    float64 `json:"temp"`
	TemperatureAge int     `json:"temp_age"`
	Humidity       float64 `json:"humidity"`
	HumidityAge    int     `json:"humidity_age"`
	DoorStatus     int     `json:"door_status"`
}

func statusHandler(resp http.ResponseWriter, req *http.Request) {
	var r JSONResponse
	jenc := json.NewEncoder(resp)

	r.Temperature = temperature
	r.TemperatureAge = int(time.Since(tempUpdated).Seconds())
	r.Humidity = humidity
	r.HumidityAge = int(time.Since(humidUpdated).Seconds())
	if door {
		r.DoorStatus = 1
	} else {
		r.DoorStatus = 0
	}
	jenc.Encode(r)
}

func setupTemperatureGraph() (g *rrd.Grapher) {
	g = rrd.NewGrapher()
	g.SetTitle("Temperature")
	g.SetWatermark("wxstation")
	g.Def("temp", rrdFile, dsTemperature, "AVERAGE")
	g.VDef("avg", "temp,AVERAGE")
	//g.SetAltAutoscale()
	//g.VDef("min", "temp,MINNAN")
	//g.VDef("max", "temp,MAX")
	g.Line(1.0, "temp", "ff0000", "Temperature (C)")
	//g.Print("min", "Minimum = %3.1lfC")
	g.GPrint("avg", "Mean = %3.1lfC")
	//g.Print("max", "Maximum = %3.1lfC")
	g.SetSize(800, 300)

	// lesigh.  no magic for this one - we have to do it ourselves.
	g.AddOptions("-y", "0.5:2")
	//g.AddOptions("--left-axis-format", "%3.1lf")

	return g
}

func setupHumidityGraph() (g *rrd.Grapher) {
	g = rrd.NewGrapher()
	g.SetTitle("Humidity")
	g.SetWatermark("wxstation")
	g.Def("humidity", rrdFile, dsHumidity, "AVERAGE")
	g.VDef("avg", "humidity,AVERAGE")
	//g.SetAltAutoscale()
	//g.VDef("min", "temp,MINNAN")
	//g.VDef("max", "temp,MAX")
	g.Line(1.0, "humidity", "00ff00", "Humidity (%)")
	//g.Print("min", "Minimum = %3.1lfC")
	g.GPrint("avg", "Mean = %2.1lf%%")
	//g.Print("max", "Maximum = %3.1lfC")
	g.SetSize(800, 300)

	// lesigh.  no magic for this one - we have to do it ourselves.
	g.AddOptions("-y", "5.0:2")
	//g.AddOptions("--left-axis-format", "%2.0lf")

	return g
}

func tempDayGraphHandler(resp http.ResponseWriter, req *http.Request) {
	g := setupTemperatureGraph()

	now := time.Now()
	_, buf, err := g.Graph(now.Add(-24*time.Hour), now)
	if err != nil {
		resp.WriteHeader(500)
		log.Println(err)
		return
	}
	resp.Write(buf)
}

func humidityDayGraphHandler(resp http.ResponseWriter, req *http.Request) {
	g := setupHumidityGraph()

	now := time.Now()
	_, buf, err := g.Graph(now.Add(-24*time.Hour), now)
	if err != nil {
		resp.WriteHeader(500)
		log.Println(err)
		return
	}
	resp.Write(buf)
}
