// Package states provides printer state types and constants.
package states

import "fmt"

// PrintStatus represents the current state of the printer.
type PrintStatus int

const (
	PrintStatusPrinting                             PrintStatus = 0
	PrintStatusAutoBedLeveling                      PrintStatus = 1
	PrintStatusHeatbedPreheating                    PrintStatus = 2
	PrintStatusSweepingXYMechMode                   PrintStatus = 3
	PrintStatusChangingFilament                     PrintStatus = 4
	PrintStatusM400Pause                            PrintStatus = 5
	PrintStatusPausedFilamentRunout                 PrintStatus = 6
	PrintStatusHeatingHotend                        PrintStatus = 7
	PrintStatusCalibratingExtrusion                 PrintStatus = 8
	PrintStatusScanningBedSurface                   PrintStatus = 9
	PrintStatusInspectingFirstLayer                 PrintStatus = 10
	PrintStatusIdentifyingBuildPlateType            PrintStatus = 11
	PrintStatusCalibratingMicroLidar                PrintStatus = 12
	PrintStatusHomingToolhead                       PrintStatus = 13
	PrintStatusCleaningNozzleTip                    PrintStatus = 14
	PrintStatusCheckingExtruderTemperature          PrintStatus = 15
	PrintStatusPausedUser                           PrintStatus = 16
	PrintStatusPausedFrontCoverFalling              PrintStatus = 17
	PrintStatusCalibratingLidar                     PrintStatus = 18
	PrintStatusCalibratingExtrusionFlow             PrintStatus = 19
	PrintStatusPausedNozzleTemperatureMalfunction   PrintStatus = 20
	PrintStatusPausedHeatBedTemperatureMalfunction  PrintStatus = 21
	PrintStatusFilamentUnloading                    PrintStatus = 22
	PrintStatusPausedSkippedStep                    PrintStatus = 23
	PrintStatusFilamentLoading                      PrintStatus = 24
	PrintStatusCalibratingMotorNoise                PrintStatus = 25
	PrintStatusPausedAMSLost                        PrintStatus = 26
	PrintStatusPausedLowFanSpeedHeatBreak           PrintStatus = 27
	PrintStatusPausedChamberTemperatureControlError PrintStatus = 28
	PrintStatusCoolingChamber                       PrintStatus = 29
	PrintStatusPausedUserGcode                      PrintStatus = 30
	PrintStatusMotorNoiseShowoff                    PrintStatus = 31
	PrintStatusPausedNozzleFilamentCoveredDetected  PrintStatus = 32
	PrintStatusPausedCutterError                    PrintStatus = 33
	PrintStatusPausedFirstLayerError                PrintStatus = 34
	PrintStatusPausedNozzleClog                     PrintStatus = 35
	PrintStatusUnknown                              PrintStatus = -1
	PrintStatusIdle                                 PrintStatus = 255
)

func (s PrintStatus) String() string {
	switch s {
	case PrintStatusPrinting:
		return "PRINTING"
	case PrintStatusAutoBedLeveling:
		return "AUTO_BED_LEVELING"
	case PrintStatusHeatbedPreheating:
		return "HEATBED_PREHEATING"
	case PrintStatusSweepingXYMechMode:
		return "SWEEPING_XY_MECH_MODE"
	case PrintStatusChangingFilament:
		return "CHANGING_FILAMENT"
	case PrintStatusM400Pause:
		return "M400_PAUSE"
	case PrintStatusPausedFilamentRunout:
		return "PAUSED_FILAMENT_RUNOUT"
	case PrintStatusHeatingHotend:
		return "HEATING_HOTEND"
	case PrintStatusCalibratingExtrusion:
		return "CALIBRATING_EXTRUSION"
	case PrintStatusScanningBedSurface:
		return "SCANNING_BED_SURFACE"
	case PrintStatusInspectingFirstLayer:
		return "INSPECTING_FIRST_LAYER"
	case PrintStatusIdentifyingBuildPlateType:
		return "IDENTIFYING_BUILD_PLATE_TYPE"
	case PrintStatusCalibratingMicroLidar:
		return "CALIBRATING_MICRO_LIDAR"
	case PrintStatusHomingToolhead:
		return "HOMING_TOOLHEAD"
	case PrintStatusCleaningNozzleTip:
		return "CLEANING_NOZZLE_TIP"
	case PrintStatusCheckingExtruderTemperature:
		return "CHECKING_EXTRUDER_TEMPERATURE"
	case PrintStatusPausedUser:
		return "PAUSED_USER"
	case PrintStatusPausedFrontCoverFalling:
		return "PAUSED_FRONT_COVER_FALLING"
	case PrintStatusCalibratingLidar:
		return "CALIBRATING_LIDAR"
	case PrintStatusCalibratingExtrusionFlow:
		return "CALIBRATING_EXTRUSION_FLOW"
	case PrintStatusPausedNozzleTemperatureMalfunction:
		return "PAUSED_NOZZLE_TEMPERATURE_MALFUNCTION"
	case PrintStatusPausedHeatBedTemperatureMalfunction:
		return "PAUSED_HEAT_BED_TEMPERATURE_MALFUNCTION"
	case PrintStatusFilamentUnloading:
		return "FILAMENT_UNLOADING"
	case PrintStatusPausedSkippedStep:
		return "PAUSED_SKIPPED_STEP"
	case PrintStatusFilamentLoading:
		return "FILAMENT_LOADING"
	case PrintStatusCalibratingMotorNoise:
		return "CALIBRATING_MOTOR_NOISE"
	case PrintStatusPausedAMSLost:
		return "PAUSED_AMS_LOST"
	case PrintStatusPausedLowFanSpeedHeatBreak:
		return "PAUSED_LOW_FAN_SPEED_HEAT_BREAK"
	case PrintStatusPausedChamberTemperatureControlError:
		return "PAUSED_CHAMBER_TEMPERATURE_CONTROL_ERROR"
	case PrintStatusCoolingChamber:
		return "COOLING_CHAMBER"
	case PrintStatusPausedUserGcode:
		return "PAUSED_USER_GCODE"
	case PrintStatusMotorNoiseShowoff:
		return "MOTOR_NOISE_SHOWOFF"
	case PrintStatusPausedNozzleFilamentCoveredDetected:
		return "PAUSED_NOZZLE_FILAMENT_COVERED_DETECTED"
	case PrintStatusPausedCutterError:
		return "PAUSED_CUTTER_ERROR"
	case PrintStatusPausedFirstLayerError:
		return "PAUSED_FIRST_LAYER_ERROR"
	case PrintStatusPausedNozzleClog:
		return "PAUSED_NOZZLE_CLOG"
	case PrintStatusIdle:
		return "IDLE"
	default:
		return "UNKNOWN"
	}
}

// GcodeState represents the G-code state of the printer.
type GcodeState string

const (
	GcodeStateIdle    GcodeState = "IDLE"
	GcodeStatePrepare GcodeState = "PREPARE"
	GcodeStateRunning GcodeState = "RUNNING"
	GcodeStatePause   GcodeState = "PAUSE"
	GcodeStateFinish  GcodeState = "FINISH"
	GcodeStateFailed  GcodeState = "FAILED"
	GcodeStateUnknown GcodeState = "UNKNOWN"
)

func (s GcodeState) String() string {
	return string(s)
}

// ParseGcodeState parses a string into a GcodeState.
func ParseGcodeState(s string) GcodeState {
	switch s {
	case "IDLE":
		return GcodeStateIdle
	case "PREPARE":
		return GcodeStatePrepare
	case "RUNNING":
		return GcodeStateRunning
	case "PAUSE":
		return GcodeStatePause
	case "FINISH":
		return GcodeStateFinish
	case "FAILED":
		return GcodeStateFailed
	default:
		return GcodeStateUnknown
	}
}

// ParsePrintStatus parses an interface{} into a PrintStatus.
func ParsePrintStatus(v interface{}) PrintStatus {
	switch val := v.(type) {
	case int:
		return PrintStatus(val)
	case float64:
		return PrintStatus(int(val))
	case string:
		var status int
		fmt.Sscanf(val, "%d", &status)
		return PrintStatus(status)
	default:
		return PrintStatusUnknown
	}
}
