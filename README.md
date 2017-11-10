# Introduction

This is a first cut at a driver for the Xpider robot (xpider.me) for golang.

![alt text](http://res.cloudinary.com/hrscywv4p/image/upload/c_limit,fl_lossy,h_9000,w_1200,f_auto,q_auto/v1/666328/620_415_2-01_cwuoh9.jpg)


There are a few assumptions:

	* You've already connected to the Xpider WIFI
		* Xpider_XXX, Pass: 12345678
	* The Xpider is listening for commands on the IP address 
		* 192.168.100.1

This driver does not capture the video feed from the robot, to do this you'll needa
to fire up a connection to rtsp://admin:admin@192.168.100.1:554/cam1/h264

# Getting started

'''
	go get github.com/celer/xpider
'''

# See examples/ for examples


