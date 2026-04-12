// Example application demonstrating the use of the bambusdk-go library.
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/asfrm/bambusdk-go/printer"
)

func main() {
	// Get connection details from environment or use defaults
	ipAddress := os.Getenv("BAMBU_IP")
	if ipAddress == "" {
		ipAddress = "192.168.1.200"
	}

	serial := os.Getenv("BAMBU_SERIAL")
	if serial == "" {
		serial = "AC12309BH109"
	}

	accessCode := os.Getenv("BAMBU_ACCESS_CODE")
	if accessCode == "" {
		accessCode = "12347890"
	}

	fmt.Println("Starting bambusdk-go example")
	fmt.Println("Connecting to BambuLab 3D printer")
	fmt.Printf("IP: %s\n", ipAddress)
	fmt.Printf("Serial: %s\n", serial)
	fmt.Printf("Access Code: %s\n", accessCode)

	// Create a new printer instance
	p := printer.NewPrinter(ipAddress, accessCode, serial)

	// Connect to the printer
	err := p.Connect()
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}
	defer p.Disconnect()

	// Wait for connection to establish
	time.Sleep(2 * time.Second)

	// Check connection status
	fmt.Printf("MQTT Connected: %v\n", p.MQTTClientConnected())
	fmt.Printf("Camera Alive: %v\n", p.CameraClientAlive())

	// Get printer status
	state := p.GetState()
	fmt.Printf("Printer State: %s\n", state)

	// Get print percentage
	percentage := p.GetPercentage()
	fmt.Printf("Print Percentage: %v\n", percentage)

	// Get remaining time
	remainingTime := p.GetTime()
	fmt.Printf("Remaining Time: %v\n", remainingTime)

	// Get temperatures
	bedTemp := p.GetBedTemperature()
	nozzleTemp := p.GetNozzleTemperature()
	chamberTemp := p.GetChamberTemperature()
	fmt.Printf("Bed Temperature: %.1f°C\n", bedTemp)
	fmt.Printf("Nozzle Temperature: %.1f°C\n", nozzleTemp)
	fmt.Printf("Chamber Temperature: %.1f°C\n", chamberTemp)

	// Get print speed
	printSpeed := p.GetPrintSpeed()
	fmt.Printf("Print Speed: %d%%\n", printSpeed)

	// Get file name
	fileName := p.GetFileName()
	fmt.Printf("Current File: %s\n", fileName)

	// Get light state
	lightState := p.GetLightState()
	fmt.Printf("Light State: %s\n", lightState)

	// Turn light off
	fmt.Println("Turning light off...")
	if err := p.TurnLightOff(); err != nil {
		fmt.Printf("Failed to turn light off: %v\n", err)
	}
	time.Sleep(2 * time.Second)

	// Turn light on
	fmt.Println("Turning light on...")
	if err := p.TurnLightOn(); err != nil {
		fmt.Printf("Failed to turn light on: %v\n", err)
	}

	// Get nozzle info
	nozzleType := p.NozzleType()
	nozzleDiameter := p.NozzleDiameter()
	fmt.Printf("Nozzle Type: %s, Diameter: %.1fmm\n", nozzleType, nozzleDiameter)

	// Get current layer info
	currentLayer := p.CurrentLayerNum()
	totalLayers := p.TotalLayerNum()
	fmt.Printf("Layer: %d/%d\n", currentLayer, totalLayers)

	// Get current state (detailed)
	currentState := p.GetCurrentState()
	fmt.Printf("Current State: %s\n", currentState)

	// Get AMS info (if available)
	amsHub := p.AMSHub()
	if amsHub != nil {
		fmt.Println("AMS Hub available")
		for i := range 4 {
			ams := amsHub.Get(i)
			if ams != nil {
				fmt.Printf("AMS %d: Humidity=%d%%, Temperature=%.1f°C\n", i, ams.Humidity, ams.Temperature)
			}
		}
	}

	// Get external spool info
	vtTray := p.VTTray()
	if vtTray.TrayInfoIdx != "" {
		fmt.Printf("External Spool: %s (%s)\n", vtTray.TrayIDName, vtTray.TrayType)
	}

	// Get WiFi signal
	wifiSignal := p.WifiSignal()
	fmt.Printf("WiFi Signal: %s dBm\n", wifiSignal)

	// Get print type
	printType := p.PrintType()
	fmt.Printf("Print Type: %s\n", printType)

	// Get subtask name
	subtaskName := p.SubtaskName()
	fmt.Printf("Subtask: %s\n", subtaskName)

	// Example: Send G-code (home printer)
	// fmt.Println("Homing printer...")
	// p.HomePrinter()

	// Example: Set fan speed
	// p.SetPartFanSpeedInt(128) // 50% speed

	// Example: Get camera frame
	if p.CameraClientAlive() {
		frame, err := p.GetCameraFrame()
		if err != nil {
			fmt.Printf("Camera frame error: %v\n", err)
		} else {
			fmt.Printf("Camera frame received (%d bytes base64)\n", len(frame))
		}
	}

	fmt.Println("\nExample completed successfully!")
	fmt.Println("Disconnecting...")
}
