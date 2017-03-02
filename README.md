# Weatherstation

WeatherStation (`wxstation`) is a small daemon to measure and record the 
temperature from a hardware probe.

This project originally was written to run using a USB-enabled Arduino
connected to a small computer (we actually deployed using a Raspberry Pi
in the end as we had one spare, but any Linux PC with a USB port would work).

# Parts Required

* DHT11 (or compatible) sensor.
* Arduino with USB Serial (we used a Leonardo clone)
* Host that can run Go binaries.
* (optional) switch to monitor door (or whatever) state.

# Assembly

See the notes in `hw/` on how to connect up the Arduino.

# Building the software

Ensure you have the `librrd` headers installed before build.

```
export GOPATH=$(pwd)
go install wxstation
```

# Running the software

Ensure you start `wxstation` with the top level directory of this tree as
the CWD otherwise `wxstation` will not find the static assets to serve!

The binary accepts the following arguments:

`-rrd-file=<filename>` specifies the location to write out the RRD datafile.

`-device=<path>` specifies the Serial port to use for communicating with the
temperature probe.

`-bind=<addrspec>` specifies the address and TCP port to bind the web server
to.

See the `systemd/wxstation.service` file for suggestions on how to run
wxstation as a service.

# Using the Software

Point a browser at http port.

There is a page available at `/status` which contains the last observed data 
and it's age encoded in JSON.

For example:

> {"temp":28,"temp_age":1,"humidity":28,"humidity_age":1,"door_status":1}

Which indicates a temperature of 28C, relative humidity of 28%, that both the
temperature and humidity datums were updated 1 second ago, and the door
sense wire is currently high.
