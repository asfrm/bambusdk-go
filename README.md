# Bambu Lab API - Go

Native Go implementation of the Bambu Lab 3D Printer API.

## Features

- **MQTT Control**: Full control over printer via MQTT protocol
- **Camera Stream**: Real-time camera frame retrieval
- **FTP File Transfer**: Upload/download files to/from printer
- **AMS Support**: Monitor and control AMS (Automated Material System)
- **Temperature Control**: Bed and nozzle temperature management
- **Fan Control**: Part, auxiliary, and chamber fan speed control
- **Print Management**: Start, stop, pause, resume prints
- **G-code Commands**: Send custom G-code commands

## Quick Start

```go
package main

import (
    "fmt"
    "time"
    
    bl "github.com/bambulabs_api/go/bambulabs_api"
)

func main() {
    // Create printer instance
    printer := bl.NewPrinter("192.168.1.200", "YOUR_ACCESS_CODE", "YOUR_SERIAL")
    
    // Connect to printer
    err := printer.Connect()
    if err != nil {
        panic(err)
    }
    defer printer.Disconnect()
    
    // Wait for connection
    time.Sleep(2 * time.Second)
    
    // Get printer status
    fmt.Printf("State: %s\n", printer.GetState())
    fmt.Printf("Bed Temp: %.1f°C\n", printer.GetBedTemperature())
    fmt.Printf("Nozzle Temp: %.1f°C\n", printer.GetNozzleTemperature())
    
    // Control printer
    printer.TurnLightOn()
    printer.HomePrinter()
}
```

## Usage Examples

### Monitor Print Progress

```go
printer := bl.NewPrinter(ip, accessCode, serial)
printer.Connect()
defer printer.Disconnect()

for {
    percentage := printer.GetPercentage()
    remaining := printer.GetTime()
    state := printer.GetState()
    
    fmt.Printf("Progress: %v%%, Remaining: %v, State: %s\n", 
        percentage, remaining, state)
    
    time.Sleep(5 * time.Second)
}
```

### Send G-code Commands

```go
// Home printer
printer.HomePrinter()

// Set bed temperature
printer.SetBedTemperature(60)

// Set nozzle temperature
printer.SetNozzleTemperature(200)

// Custom G-code
printer.Gcode("G1 X100 Y100 F3000", true)
```

### Control Fans

```go
// Set part fan speed (0-255)
printer.SetPartFanSpeedInt(128)

// Or use percentage (0.0-1.0)
printer.SetPartFanSpeed(0.5)

// Control auxiliary and chamber fans
printer.SetAuxFanSpeedInt(200)
printer.SetChamberFanSpeedInt(100)
```

### AMS Management

```go
// Get AMS hub info
amsHub := printer.AMSHub()

// Access individual AMS units
ams := amsHub.Get(0)
if ams != nil {
    fmt.Printf("Humidity: %d%%\n", ams.Humidity)
    fmt.Printf("Temperature: %.1f°C\n", ams.Temperature)
    
    // Get filament tray info
    tray := ams.GetFilamentTray(0)
    if tray != nil {
        fmt.Printf("Filament: %s\n", tray.TrayType)
    }
}
```

### Camera Access

```go
// Get camera frame as base64
frame, err := printer.GetCameraFrame()
if err != nil {
    log.Fatal(err)
}

// Get camera frame as image
img, err := printer.GetCameraImage()
if err != nil {
    log.Fatal(err)
}
```

### File Operations (FTP)

```go
// Upload file
file, _ := os.Open("print.gcode")
printer.UploadFile(file, "my_print.gcode")

// Start print
printer.StartPrint("my_print.gcode", 1, true, []int{0}, nil, true)

// Download file
data, err := printer.DownloadFile("image/last_print.jpg")

// Delete file
printer.DeleteFile("cache/temp.gcode")
```

### Print Control

```go
// Start print
printer.StartPrint("file.gcode", 1, true, []int{0}, nil, true)

// Pause print
printer.PausePrint()

// Resume print
printer.ResumePrint()

// Stop print
printer.StopPrint()

// Skip objects during print
printer.SkipObjects([]int{1, 3, 5})
```

### Temperature Control

```go
// Set bed temperature
printer.SetBedTemperature(60)

// Set nozzle temperature  
printer.SetNozzleTemperature(210)

// Low temperature warning (use override if needed)
printer.SetBedTemperatureOverride(30, true)
```

### Calibration

```go
// Full calibration
printer.CalibratePrinter(true, true, true)

// Bed leveling only
printer.CalibratePrinter(true, false, false)
```

### Filament Management

```go
// Load filament
printer.LoadFilamentSpool()

// Unload filament
printer.UnloadFilamentSpool()

// Set filament settings
printer.SetFilamentPrinter("FF0000", "BAMBU_PLA", 255, 254)

// Retry filament action
printer.RetryFilamentAction()
```

## API Reference

### Printer Methods

| Method | Description |
|--------|-------------|
| `Connect()` | Connect to printer (MQTT + Camera) |
| `Disconnect()` | Disconnect from printer |
| `GetState()` | Get printer G-code state |
| `GetCurrentState()` | Get detailed printer status |
| `GetPercentage()` | Get print completion percentage |
| `GetTime()` | Get remaining print time |
| `GetBedTemperature()` | Get bed temperature |
| `GetNozzleTemperature()` | Get nozzle temperature |
| `GetChamberTemperature()` | Get chamber temperature |
| `GetPrintSpeed()` | Get print speed percentage |
| `GetFileName()` | Get current print file name |
| `TurnLightOn()` | Turn on printer light |
| `TurnLightOff()` | Turn off printer light |
| `HomePrinter()` | Home all axes |
| `MoveZAxis(height)` | Move Z-axis to height |
| `PausePrint()` | Pause current print |
| `ResumePrint()` | Resume paused print |
| `StopPrint()` | Stop current print |
| `Gcode(cmd, check)` | Send G-code command |
| `UploadFile(data, name)` | Upload file via FTP |
| `StartPrint(...)` | Start printing a file |
| `GetCameraFrame()` | Get camera frame (base64) |
| `GetCameraImage()` | Get camera frame (image) |
| `AMSHub()` | Get AMS hub information |
| `Reboot()` | Reboot printer |

## Printer States

### GcodeState

- `IDLE` - Printer is idle
- `PREPARE` - Preparing to print
- `RUNNING` - Currently printing
- `PAUSE` - Print paused
- `FINISH` - Print finished
- `FAILED` - Print failed

### PrintStatus

Detailed status including:
- `PRINTING` - Actively printing
- `AUTO_BED_LEVELING` - Leveling bed
- `HEATBED_PREHEATING` - Preheating bed
- `CHANGING_FILAMENT` - Changing filament
- `CALIBRATING_EXTRUSION` - Calibrating extrusion
- And many more...

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `BAMBU_IP` | Printer IP address | - |
| `BAMBU_SERIAL` | Printer serial number | - |
| `BAMBU_ACCESS_CODE` | Printer access code | - |

## Troubleshooting

### Connection Issues

1. Ensure printer is on the same network
2. Verify access code is correct (8 digits)
3. Check firewall settings (ports 8883, 6000, 990)

### Camera Not Working

- X1 series: Camera support may be limited
- Ensure camera is enabled in printer settings
- Check access code permissions

### FTP Upload Fails

- Verify printer is idle or in safe state
- Check available storage on printer
- Ensure file name is valid

## License

MIT License - See LICENSE file for details

## Credits

Based on the Python implementation at https://github.com/BambuTools/bambulabs_api
