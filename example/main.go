package main

import (
	"fmt"
	"os"
	"time"

	bl "github.com/asfrm/bambuapi-go/bambulabs_api"
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

	fmt.Println("Starting bambulabs_api Go example")
	fmt.Println("Connecting to BambuLab 3D printer")
	fmt.Printf("IP: %s\n", ipAddress)
	fmt.Printf("Serial: %s\n", serial)
	fmt.Printf("Access Code: %s\n", accessCode)

	// Create a new printer instance
	printer := bl.NewPrinter(ipAddress, accessCode, serial)

	// Connect to the printer
	err := printer.Connect()
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}
	defer printer.Disconnect()

	// Wait for connection to establish
	time.Sleep(2 * time.Second)

	// Check connection status
	fmt.Printf("MQTT Connected: %v\n", printer.MQTTClientConnected())
	fmt.Printf("Camera Alive: %v\n", printer.CameraClientAlive())

	// Get printer status
	state := printer.GetState()
	fmt.Printf("Printer State: %s\n", state)

	// Get print percentage
	percentage := printer.GetPercentage()
	fmt.Printf("Print Percentage: %v\n", percentage)

	// Get remaining time
	remainingTime := printer.GetTime()
	fmt.Printf("Remaining Time: %v\n", remainingTime)

	// Get temperatures
	bedTemp := printer.GetBedTemperature()
	nozzleTemp := printer.GetNozzleTemperature()
	chamberTemp := printer.GetChamberTemperature()
	fmt.Printf("Bed Temperature: %.1f°C\n", bedTemp)
	fmt.Printf("Nozzle Temperature: %.1f°C\n", nozzleTemp)
	fmt.Printf("Chamber Temperature: %.1f°C\n", chamberTemp)

	// Get print speed
	printSpeed := printer.GetPrintSpeed()
	fmt.Printf("Print Speed: %d%%\n", printSpeed)

	// Get file name
	fileName := printer.GetFileName()
	fmt.Printf("Current File: %s\n", fileName)

	// Get light state
	lightState := printer.GetLightState()
	fmt.Printf("Light State: %s\n", lightState)

	// Turn light off
	fmt.Println("Turning light off...")
	printer.TurnLightOff()
	time.Sleep(2 * time.Second)

	// Turn light on
	fmt.Println("Turning light on...")
	printer.TurnLightOn()

	// Get nozzle info
	nozzleType := printer.NozzleType()
	nozzleDiameter := printer.NozzleDiameter()
	fmt.Printf("Nozzle Type: %s, Diameter: %.1fmm\n", nozzleType, nozzleDiameter)

	// Get current layer info
	currentLayer := printer.CurrentLayerNum()
	totalLayers := printer.TotalLayerNum()
	fmt.Printf("Layer: %d/%d\n", currentLayer, totalLayers)

	// Get current state (detailed)
	currentState := printer.GetCurrentState()
	fmt.Printf("Current State: %s\n", currentState)

	// Get AMS info (if available)
	amsHub := printer.AMSHub()
	if amsHub != nil {
		fmt.Println("AMS Hub available")
		for i := 0; i < 4; i++ {
			ams := amsHub.Get(i)
			if ams != nil {
				fmt.Printf("AMS %d: Humidity=%d%%, Temperature=%.1f°C\n", i, ams.Humidity, ams.Temperature)
			}
		}
	}

	// Get external spool info
	vtTray := printer.VTTray()
	if vtTray.TrayInfoIdx != "" {
		fmt.Printf("External Spool: %s (%s)\n", vtTray.TrayIDName, vtTray.TrayType)
	}

	// Get WiFi signal
	wifiSignal := printer.WifiSignal()
	fmt.Printf("WiFi Signal: %s dBm\n", wifiSignal)

	// Get print type
	printType := printer.PrintType()
	fmt.Printf("Print Type: %s\n", printType)

	// Get subtask name
	subtaskName := printer.SubtaskName()
	fmt.Printf("Subtask: %s\n", subtaskName)

	// Example: Send G-code (home printer)
	// fmt.Println("Homing printer...")
	// printer.HomePrinter()

	// Example: Set fan speed
	// printer.SetPartFanSpeedInt(128) // 50% speed

	// Example: Get camera frame
	if printer.CameraClientAlive() {
		frame, err := printer.GetCameraFrame()
		if err != nil {
			fmt.Printf("Camera frame error: %v\n", err)
		} else {
			fmt.Printf("Camera frame received (%d bytes base64)\n", len(frame))
		}
	}

	fmt.Println("\nExample completed successfully!")
	fmt.Println("Disconnecting...")
}
