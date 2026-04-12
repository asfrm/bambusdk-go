// Package fleet provides printer pool management for handling multiple printers concurrently.
package fleet

import (
	"fmt"
	"maps"
	"sync"

	"github.com/asfrm/bambusdk-go/printer"
)

// PrinterConfig holds the configuration for a printer.
type PrinterConfig struct {
	IP         string
	AccessCode string
	Serial     string
	Name       string // Optional friendly name
}

// PrinterInfo contains information about a printer in the pool.
type PrinterInfo struct {
	Serial      string
	Name        string
	IP          string
	Connected   bool
	State       string
	Temperature struct {
		Nozzle  float64
		Bed     float64
		Chamber float64
	}
	Progress int
}

// PrinterPool manages a collection of printers for fleet operations.
// It is safe for concurrent use by multiple goroutines.
type PrinterPool struct {
	mu       sync.RWMutex
	printers map[string]*printer.Printer
	configs  map[string]*PrinterConfig
}

// NewPrinterPool creates a new empty printer pool.
func NewPrinterPool() *PrinterPool {
	return &PrinterPool{
		printers: make(map[string]*printer.Printer),
		configs:  make(map[string]*PrinterConfig),
	}
}

// AddPrinter adds a printer configuration to the pool.
// The printer is not connected until ConnectPrinter or ConnectAll is called.
func (p *PrinterPool) AddPrinter(cfg *PrinterConfig) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if cfg.Serial == "" {
		return fmt.Errorf("serial number is required")
	}
	if cfg.IP == "" {
		return fmt.Errorf("IP address is required")
	}
	if cfg.AccessCode == "" {
		return fmt.Errorf("access code is required")
	}

	// Check if already exists
	if _, exists := p.configs[cfg.Serial]; exists {
		return fmt.Errorf("printer with serial %s already exists", cfg.Serial)
	}

	// Store config
	p.configs[cfg.Serial] = cfg

	return nil
}

// AddPrinters adds multiple printer configurations to the pool.
func (p *PrinterPool) AddPrinters(configs []*PrinterConfig) error {
	for _, cfg := range configs {
		if err := p.AddPrinter(cfg); err != nil {
			return fmt.Errorf("failed to add printer %s: %w", cfg.Serial, err)
		}
	}
	return nil
}

// RemovePrinter removes a printer from the pool.
// If the printer is connected, it will be disconnected first.
func (p *PrinterPool) RemovePrinter(serial string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Disconnect if connected
	if prtr, exists := p.printers[serial]; exists {
		prtr.Disconnect()
		delete(p.printers, serial)
	}

	// Remove config
	delete(p.configs, serial)

	return nil
}

// GetPrinter returns a printer instance by serial number.
// Returns nil if the printer is not in the pool or not connected.
func (p *PrinterPool) GetPrinter(serial string) *printer.Printer {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.printers[serial]
}

// GetPrinterConfig returns the configuration for a printer.
func (p *PrinterPool) GetPrinterConfig(serial string) *PrinterConfig {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.configs[serial]
}

// ConnectPrinter connects a specific printer in the pool.
// The printer must have been added via AddPrinter first.
func (p *PrinterPool) ConnectPrinter(serial string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	cfg, exists := p.configs[serial]
	if !exists {
		return fmt.Errorf("printer %s not found in pool", serial)
	}

	// Check if already connected
	if _, exists := p.printers[serial]; exists {
		return nil // Already connected
	}

	// Create and connect printer
	prtr := printer.NewPrinter(cfg.IP, cfg.AccessCode, cfg.Serial)
	if err := prtr.Connect(); err != nil {
		return fmt.Errorf("failed to connect to printer %s: %w", serial, err)
	}

	p.printers[serial] = prtr
	return nil
}

// ConnectAll connects all printers in the pool concurrently.
// Returns a map of serial numbers to connection errors (nil = success).
func (p *PrinterPool) ConnectAll() map[string]error {
	p.mu.Lock()
	configs := make([]*PrinterConfig, 0, len(p.configs))
	for _, cfg := range p.configs {
		configs = append(configs, cfg)
	}
	p.mu.Unlock()

	results := make(map[string]error)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, cfg := range configs {
		wg.Add(1)
		go func(cfg *PrinterConfig) {
			defer wg.Done()

			p.mu.Lock()
			_, exists := p.printers[cfg.Serial]
			p.mu.Unlock()

			if exists {
				mu.Lock()
				results[cfg.Serial] = nil // Already connected
				mu.Unlock()
				return
			}

			prtr := printer.NewPrinter(cfg.IP, cfg.AccessCode, cfg.Serial)
			if err := prtr.Connect(); err != nil {
				mu.Lock()
				results[cfg.Serial] = err
				mu.Unlock()
				return
			}

			p.mu.Lock()
			p.printers[cfg.Serial] = prtr
			p.mu.Unlock()

			mu.Lock()
			results[cfg.Serial] = nil
			mu.Unlock()
		}(cfg)
	}

	wg.Wait()
	return results
}

// DisconnectPrinter disconnects a specific printer.
func (p *PrinterPool) DisconnectPrinter(serial string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if prtr, exists := p.printers[serial]; exists {
		prtr.Disconnect()
		delete(p.printers, serial)
	}

	return nil
}

// DisconnectAll disconnects all printers in the pool.
func (p *PrinterPool) DisconnectAll() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for serial, prtr := range p.printers {
		prtr.Disconnect()
		delete(p.printers, serial)
	}
}

// ListPrinters returns a list of all printer serial numbers in the pool.
func (p *PrinterPool) ListPrinters() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	serials := make([]string, 0, len(p.configs))
	for serial := range p.configs {
		serials = append(serials, serial)
	}
	return serials
}

// ListConnectedPrinters returns a list of connected printer serial numbers.
func (p *PrinterPool) ListConnectedPrinters() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	serials := make([]string, 0, len(p.printers))
	for serial := range p.printers {
		serials = append(serials, serial)
	}
	return serials
}

// GetPrinterInfo returns detailed information about a printer.
func (p *PrinterPool) GetPrinterInfo(serial string) (*PrinterInfo, error) {
	p.mu.RLock()
	prtr, exists := p.printers[serial]
	cfg, cfgExists := p.configs[serial]
	p.mu.RUnlock()

	if !exists || !cfgExists {
		return nil, fmt.Errorf("printer %s not found", serial)
	}

	info := &PrinterInfo{
		Serial:    serial,
		Name:      cfg.Name,
		IP:        cfg.IP,
		Connected: true,
		State:     prtr.GetState().String(),
		Progress:  prtr.GetPercentage(),
	}
	info.Temperature.Nozzle = prtr.GetNozzleTemperature()
	info.Temperature.Bed = prtr.GetBedTemperature()
	info.Temperature.Chamber = prtr.GetChamberTemperature()

	return info, nil
}

// GetAllPrinterInfo returns information about all printers in the pool.
// For disconnected printers, only configuration info is returned.
func (p *PrinterPool) GetAllPrinterInfo() []*PrinterInfo {
	p.mu.RLock()
	defer p.mu.RUnlock()

	infos := make([]*PrinterInfo, 0, len(p.configs))

	for serial, cfg := range p.configs {
		info := &PrinterInfo{
			Serial:    serial,
			Name:      cfg.Name,
			IP:        cfg.IP,
			Connected: false,
		}

		if prtr, exists := p.printers[serial]; exists {
			info.Connected = true
			info.State = prtr.GetState().String()
			info.Progress = prtr.GetPercentage()
			info.Temperature.Nozzle = prtr.GetNozzleTemperature()
			info.Temperature.Bed = prtr.GetBedTemperature()
			info.Temperature.Chamber = prtr.GetChamberTemperature()
		}

		infos = append(infos, info)
	}

	return infos
}

// GetPoolStatus returns a summary of the pool status.
type PoolStatus struct {
	TotalPrinters     int
	ConnectedCount    int
	DisconnectedCount int
}

func (p *PrinterPool) GetStatus() PoolStatus {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return PoolStatus{
		TotalPrinters:     len(p.configs),
		ConnectedCount:    len(p.printers),
		DisconnectedCount: len(p.configs) - len(p.printers),
	}
}

// ForEachPrinter iterates over all connected printers concurrently.
// The function fn is called for each printer with its serial number.
func (p *PrinterPool) ForEachPrinter(fn func(serial string, prtr *printer.Printer)) {
	p.mu.RLock()
	printers := make(map[string]*printer.Printer)
	maps.Copy(printers, p.printers)
	p.mu.RUnlock()

	var wg sync.WaitGroup
	for serial, prtr := range printers {
		wg.Add(1)
		go func(serial string, prtr *printer.Printer) {
			defer wg.Done()
			fn(serial, prtr)
		}(serial, prtr)
	}
	wg.Wait()
}

// ForEachPrinterSerial iterates over all printer serials (connected or not).
func (p *PrinterPool) ForEachPrinterSerial(fn func(serial string, cfg *PrinterConfig)) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for serial, cfg := range p.configs {
		fn(serial, cfg)
	}
}

// SetStateUpdateCallback sets a state update callback for all connected printers.
func (p *PrinterPool) SetStateUpdateCallback(callback func(serial string)) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for serial, prtr := range p.printers {
		s := serial
		prtr.SetStateUpdateCallback(func() {
			callback(s)
		})
	}
}

// PrinterStatus contains the status of a single printer.
type PrinterStatus struct {
	Serial      string
	State       string
	Progress    int
	NozzleTemp  float64
	BedTemp     float64
	ChamberTemp float64
}

// BroadcastGcode sends a G-code command to all connected printers.
// Returns a map of serial numbers to success results.
func (p *PrinterPool) BroadcastGcode(gcode string, gcodeCheck bool) map[string]bool {
	p.mu.RLock()
	printers := make(map[string]*printer.Printer)
	maps.Copy(printers, p.printers)
	p.mu.RUnlock()

	results := make(map[string]bool)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for serial, prtr := range printers {
		wg.Add(1)
		go func(serial string, prtr *printer.Printer) {
			defer wg.Done()
			success, err := prtr.Gcode(gcode, gcodeCheck)
			mu.Lock()
			if err != nil {
				results[serial] = false
			} else {
				results[serial] = success
			}
			mu.Unlock()
		}(serial, prtr)
	}
	wg.Wait()
	return results
}

// BroadcastStatus gets the status from all connected printers.
// Returns a map of serial numbers to PrinterStatus.
func (p *PrinterPool) BroadcastStatus() map[string]*PrinterStatus {
	p.mu.RLock()
	printers := make(map[string]*printer.Printer)
	maps.Copy(printers, p.printers)
	p.mu.RUnlock()

	results := make(map[string]*PrinterStatus)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for serial, prtr := range printers {
		wg.Add(1)
		go func(serial string, prtr *printer.Printer) {
			defer wg.Done()
			status := &PrinterStatus{
				Serial:      serial,
				State:       prtr.GetState().String(),
				Progress:    prtr.GetPercentage(),
				NozzleTemp:  prtr.GetNozzleTemperature(),
				BedTemp:     prtr.GetBedTemperature(),
				ChamberTemp: prtr.GetChamberTemperature(),
			}
			mu.Lock()
			results[serial] = status
			mu.Unlock()
		}(serial, prtr)
	}
	wg.Wait()
	return results
}
