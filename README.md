# Introduction

This is a first cut at a driver for the [Xpider](http://xpider.me) robot for golang.

![Picture of a Xpider Robot](http://res.cloudinary.com/hrscywv4p/image/upload/c_limit,fl_lossy,h_9000,w_1200,f_auto,q_auto/v1/666328/620_415_2-01_cwuoh9.jpg)


There are a few assumptions:

	* You've already connected to the Xpider WIFI
		* Xpider_XXX, Pass: 12345678
	* The Xpider is listening for commands on the IP address 
		* 192.168.100.1

This driver does not capture the video feed from the robot, to do this you'll needa
to fire up a connection to rtsp://admin:admin@192.168.100.1:554/cam1/h264

# Getting started

```
	go get github.com/celer/xpider
```

## Connect to the robot

To get to this point you'll need to connect to the 
wifi network that the Xpider robot uses, mine appears
as Xpider_242, yours will probably look similar. The
robot is simply listening on an IP address on the network
it's setup. 

```go
	// Connect to the xpider
	x := &xpider.Controller{}
	err := x.Connect("192.168.100.1:80")
	if err != nil {
		panic(err)
	}
```

## Send a command to the robot

```go

	// Set our front LEDs to be green and red
	x.FrontLED(0, 0xFF, 0, 0xFF, 0, 0)
```

## Get the state of the robot

```go
	state:=x.GetState()
	fmt.Printf("Observed Distance %d\n", state.ObsticalDistance)

```
# See examples/ for examples


