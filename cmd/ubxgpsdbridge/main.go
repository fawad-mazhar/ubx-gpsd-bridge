package main

import (
	"log"

	"example.com/ubx-gpsd-bridge/pkg/gpsd"
	"example.com/ubx-gpsd-bridge/pkg/ubx"
)

func main() {
	ubxDevice, err := ubx.NewI2CDevice(1, 0x42)
	if err != nil {
		log.Fatalf("Failed to create UBX I2C device: %v", err)
	}

	feeder, err := gpsd.NewFeeder(ubxDevice, "127.0.0.1:49000")
	if err != nil {
		log.Fatalf("Failed to create UBX GPSD Feeder: %v", err)
	}

	log.Println("Starting UBX GPSD Feeder...")
	if err := feeder.Serve(); err != nil {
		log.Fatalf("UBX GPSD Feeder failed: %v", err)
	}
}
