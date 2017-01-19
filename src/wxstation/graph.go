package main

import (
	"github.com/ziutek/rrd"
	"log"
	"net/http"
	"time"
)

func init() {
	http.HandleFunc("/g/temp/day", getGraphHandler(setupTemperatureGraph(), 24*time.Hour, "Temperature - Last 24 Hours"))
	http.HandleFunc("/g/temp/week", getGraphHandler(setupTemperatureGraph(), 7*24*time.Hour, "Temperature - Last 7 Days"))
	http.HandleFunc("/g/temp/month", getGraphHandler(setupTemperatureGraph(), 31*24*time.Hour, "Temperature - Last Month"))
	http.HandleFunc("/g/temp/year", getGraphHandler(setupTemperatureGraph(), 365*24*time.Hour, "Temperature - Last Year"))

	http.HandleFunc("/g/humidity/day", getGraphHandler(setupHumidityGraph(), 24*time.Hour, "Humidity - Last 24 Hours"))
	http.HandleFunc("/g/humidity/week", getGraphHandler(setupHumidityGraph(), 7*24*time.Hour, "Humidity - Last 7 Days"))
	http.HandleFunc("/g/humidity/month", getGraphHandler(setupHumidityGraph(), 31*24*time.Hour, "Humidity - Last Month"))
	http.HandleFunc("/g/humidity/year", getGraphHandler(setupHumidityGraph(), 365*24*time.Hour, "Humidity - Last Year"))
}

func setupTemperatureGraph() (g *rrd.Grapher) {
	g = rrd.NewGrapher()
	g.SetTitle("Temperature")
	g.SetWatermark("wxstation")
	g.Def("temp", rrdFile, dsTemperature, "AVERAGE")
	g.VDef("avg", "temp,AVERAGE")
	//g.SetAltAutoscale()
	g.VDef("min", "temp,MINIMUM")
	g.VDef("max", "temp,MAXIMUM")
	g.Line(1.0, "temp", "ff0000", "Temperature (C)")
	g.GPrint("min", "Minimum = %3.1lfC")
	g.GPrint("avg", "Mean = %3.1lfC")
	g.GPrint("max", "Maximum = %3.1lfC")
	g.SetSize(800, 150)

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
	g.VDef("min", "humidity,MINIMUM")
	g.VDef("max", "humidity,MAXIMUM")
	g.Line(1.0, "humidity", "00ff00", "Humidity (%)")
	g.GPrint("min", "Minimum = %2.1lf%%")
	g.GPrint("avg", "Mean = %2.1lf%%")
	g.GPrint("max", "Maximum = %2.1lf%%")
	g.SetSize(800, 150)

	// lesigh.  no magic for this one - we have to do it ourselves.
	g.AddOptions("-y", "5.0:2")
	//g.AddOptions("--left-axis-format", "%2.0lf")

	return g
}

// dodgy shortcut!
func getGraphHandler(g *rrd.Grapher, interval time.Duration, title string) (handler func(http.ResponseWriter, *http.Request)) {
	return func(resp http.ResponseWriter, req *http.Request) {
		g.SetTitle(title)
		now := time.Now()
		_, buf, err := g.Graph(now.Add(-interval), now)
		if err != nil {
			resp.WriteHeader(500)
			log.Println(err)
			return
		}
		resp.Write(buf)
	}
}
