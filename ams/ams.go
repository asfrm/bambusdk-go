// Package ams provides AMS (Automated Material System) support for Bambu Lab printers.
package ams

import "github.com/asfrm/bambuapi-go/filament"

// AMS represents the Bambu Lab AMS (Automated Material System).
type AMS struct {
	FilamentTrays map[int]*filament.FilamentTray
	Humidity      int
	Temperature   float64
}

// NewAMS creates a new AMS instance.
func NewAMS(humidity int, temperature float64) *AMS {
	return &AMS{
		FilamentTrays: make(map[int]*filament.FilamentTray),
		Humidity:      humidity,
		Temperature:   temperature,
	}
}

// SetFilamentTray sets a filament tray at the given index.
func (a *AMS) SetFilamentTray(trayIndex int, tray *filament.FilamentTray) {
	a.FilamentTrays[trayIndex] = tray
}

// GetFilamentTray gets the filament tray at the given index.
func (a *AMS) GetFilamentTray(trayIndex int) *filament.FilamentTray {
	if tray, ok := a.FilamentTrays[trayIndex]; ok {
		return tray
	}
	return nil
}

// ProcessTrays processes a list of tray data.
func (a *AMS) ProcessTrays(trays []map[string]interface{}) {
	for _, t := range trays {
		var id int
		if v, ok := t["id"]; ok {
			switch val := v.(type) {
			case float64:
				id = int(val)
			case int:
				id = val
			}
		}
		if _, ok := t["n"]; ok {
			tray := filament.FilamentTrayFromDict(t)
			a.SetFilamentTray(id, &tray)
		}
	}
}

// AMSHub holds all AMS units connected to the printer.
type AMSHub struct {
	AMSHub map[int]*AMS
}

// NewAMSHub creates a new AMSHub instance.
func NewAMSHub() *AMSHub {
	return &AMSHub{
		AMSHub: make(map[int]*AMS),
	}
}

// Get gets an AMS by index.
func (h *AMSHub) Get(index int) *AMS {
	if ams, ok := h.AMSHub[index]; ok {
		return ams
	}
	return nil
}

// Set sets an AMS at the given index.
func (h *AMSHub) Set(index int, ams *AMS) {
	h.AMSHub[index] = ams
}

// ParseList parses a list of AMS data.
func (h *AMSHub) ParseList(amsList []map[string]interface{}) {
	for _, a := range amsList {
		var id int
		if v, ok := a["id"]; ok {
			switch val := v.(type) {
			case float64:
				id = int(val)
			case int:
				id = val
			}
		}
		humidity := 0
		if v, ok := a["humidity"]; ok {
			switch val := v.(type) {
			case float64:
				humidity = int(val)
			case int:
				humidity = val
			}
		}
		temp := 0.0
		if v, ok := a["temp"]; ok {
			switch val := v.(type) {
			case float64:
				temp = val
			case int:
				temp = float64(val)
			}
		}

		ams := NewAMS(humidity, temp)
		if trays, ok := a["tray"].([]interface{}); ok {
			trayMaps := make([]map[string]interface{}, len(trays))
			for i, t := range trays {
				if tm, ok := t.(map[string]interface{}); ok {
					trayMaps[i] = tm
				}
			}
			ams.ProcessTrays(trayMaps)
		}
		h.Set(id, ams)
	}
}
