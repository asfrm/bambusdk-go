package printer

import (
	"testing"

	"github.com/asfrm/bambusdk-go/filament"
	"github.com/asfrm/bambusdk-go/printerinfo"
	"github.com/asfrm/bambusdk-go/states"
)

func TestNewPrinter(t *testing.T) {
	p := NewPrinter("192.168.1.100", "12345678", "TEST123")

	if p.IPAddress != "192.168.1.100" {
		t.Errorf("Expected IPAddress to be 192.168.1.100, got %s", p.IPAddress)
	}

	if p.AccessCode != "12345678" {
		t.Errorf("Expected AccessCode to be 12345678, got %s", p.AccessCode)
	}

	if p.Serial != "TEST123" {
		t.Errorf("Expected Serial to be TEST123, got %s", p.Serial)
	}
}

func TestPrintStatus(t *testing.T) {
	tests := []struct {
		value    int
		expected string
	}{
		{0, "PRINTING"},
		{1, "AUTO_BED_LEVELING"},
		{255, "IDLE"},
		{-1, "UNKNOWN"},
		{999, "UNKNOWN"},
	}

	for _, tt := range tests {
		status := states.PrintStatus(tt.value)
		if status.String() != tt.expected {
			t.Errorf("PrintStatus(%d).String() = %s, expected %s", tt.value, status.String(), tt.expected)
		}
	}
}

func TestGcodeState(t *testing.T) {
	tests := []struct {
		input    string
		expected states.GcodeState
	}{
		{"IDLE", states.GcodeStateIdle},
		{"RUNNING", states.GcodeStateRunning},
		{"PAUSE", states.GcodeStatePause},
		{"INVALID", states.GcodeStateUnknown},
	}

	for _, tt := range tests {
		result := states.ParseGcodeState(tt.input)
		if result != tt.expected {
			t.Errorf("ParseGcodeState(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

func TestNozzleType(t *testing.T) {
	tests := []struct {
		input    string
		expected printerinfo.NozzleType
	}{
		{"stainless_steel", printerinfo.NozzleTypeStainlessSteel},
		{"hardened_steel", printerinfo.NozzleTypeHardenedSteel},
		{"invalid", printerinfo.NozzleTypeStainlessSteel},
	}

	for _, tt := range tests {
		result := printerinfo.ParseNozzleType(tt.input)
		if result != tt.expected {
			t.Errorf("ParseNozzleType(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

func TestPrinterType(t *testing.T) {
	tests := []struct {
		input    string
		expected printerinfo.PrinterType
	}{
		{"P1S", printerinfo.PrinterTypeP1S},
		{"P1P", printerinfo.PrinterTypeP1P},
		{"X1C", printerinfo.PrinterTypeX1C},
		{"INVALID", printerinfo.PrinterTypeUnknown},
	}

	for _, tt := range tests {
		result := printerinfo.ParsePrinterType(tt.input)
		if result != tt.expected {
			t.Errorf("ParsePrinterType(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

func TestFilamentByName(t *testing.T) {
	tests := []struct {
		name      string
		wantError bool
	}{
		{"BAMBU_PLA_Basic", false},
		{"BAMBU_ABS", false},
		{"INVALID_FILAMENT", true},
	}

	for _, tt := range tests {
		_, err := filament.FilamentByName(tt.name)
		if (err != nil) != tt.wantError {
			t.Errorf("FilamentByName(%s) error = %v, wantError %v", tt.name, err, tt.wantError)
		}
	}
}

func TestFilamentSettings(t *testing.T) {
	filament := filament.FilamentBambuPLABasic
	settings := filament.GetSettings()

	if settings.TrayInfoIdx != "GFA00" {
		t.Errorf("Expected TrayInfoIdx to be GFA00, got %s", settings.TrayInfoIdx)
	}

	if settings.NozzleTempMin != 190 {
		t.Errorf("Expected NozzleTempMin to be 190, got %d", settings.NozzleTempMin)
	}

	if settings.NozzleTempMax != 250 {
		t.Errorf("Expected NozzleTempMax to be 250, got %d", settings.NozzleTempMax)
	}

	if settings.TrayType != "PLA" {
		t.Errorf("Expected TrayType to be PLA, got %s", settings.TrayType)
	}
}

func TestFilamentTrayFromDict(t *testing.T) {
	dict := map[string]any{
		"k":               0.5,
		"n":               1,
		"tag_uid":         "TEST123",
		"tray_id_name":    "Bambu PLA",
		"tray_info_idx":   "GFA00",
		"tray_type":       "PLA",
		"tray_sub_brands": "",
		"tray_color":      "FF0000FF",
		"tray_weight":     "1000",
		"tray_diameter":   "200",
		"tray_temp":       "25",
		"tray_time":       "2024-01-01",
		"bed_temp_type":   "1",
		"bed_temp":        "60",
		"nozzle_temp_max": 250,
		"nozzle_temp_min": 190,
		"xcam_info":       "",
		"tray_uuid":       "UUID123",
	}

	tray := filament.FilamentTrayFromDict(dict)

	if tray.K != 0.5 {
		t.Errorf("Expected K to be 0.5, got %f", tray.K)
	}

	if tray.N != 1 {
		t.Errorf("Expected N to be 1, got %d", tray.N)
	}

	if tray.TrayInfoIdx != "GFA00" {
		t.Errorf("Expected TrayInfoIdx to be GFA00, got %s", tray.TrayInfoIdx)
	}

	if tray.NozzleTempMin != 190 {
		t.Errorf("Expected NozzleTempMin to be 190, got %d", tray.NozzleTempMin)
	}

	if tray.NozzleTempMax != 250 {
		t.Errorf("Expected NozzleTempMax to be 250, got %d", tray.NozzleTempMax)
	}
}

func TestIsValidNozzleDiameter(t *testing.T) {
	tests := []struct {
		diameter float64
		expected bool
	}{
		{0.2, true},
		{0.4, true},
		{0.6, true},
		{0.8, true},
		{0.5, false},
		{1.0, false},
	}

	for _, tt := range tests {
		result := printerinfo.IsValidNozzleDiameter(tt.diameter)
		if result != tt.expected {
			t.Errorf("IsValidNozzleDiameter(%f) = %v, expected %v", tt.diameter, result, tt.expected)
		}
	}
}
