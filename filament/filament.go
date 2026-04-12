// Package filament provides filament type definitions and settings for Bambu Lab printers.
package filament

import (
	"encoding/json"
	"fmt"
)

// AMSFilamentSettings holds filament settings.
type AMSFilamentSettings struct {
	TrayInfoIdx   string `json:"tray_info_idx"`
	NozzleTempMin int    `json:"nozzle_temp_min"`
	NozzleTempMax int    `json:"nozzle_temp_max"`
	TrayType      string `json:"tray_type"`
}

// Filament represents different filament types.
type Filament int

const (
	FilamentPolyLitePLA Filament = iota
	FilamentPolyTerraPLA
	FilamentBambuABS
	FilamentBambuPACF
	FilamentBambuPC
	FilamentBambuPLABasic
	FilamentBambuPLAMatte
	FilamentSupportG
	FilamentSupportW
	FilamentBambuTPU95A
	FilamentBambuASAAero
	FilamentBambuPLAMetal
	FilamentBambuPETGTranslucent
	FilamentBambuPLAMarble
	FilamentBambuPLAWood
	FilamentBambuPLASilkPlus
	FilamentBambuPETGHF
	FilamentBambuTPUForAMS
	FilamentBambuSupportForABS
	FilamentBambuPCFR
	FilamentBambuPLAGalaxy
	FilamentBambuPA6GF
	FilamentBambuPLAAero
	FilamentBambuASACF
	FilamentBambuPETGCF
	FilamentBambuSupportForPAPET
	FilamentBambuPLASparkle
	FilamentBambuABSGF
	FilamentBambuPAHTCF
	FilamentBambuPA6CF
	FilamentBambuPLASilk
	FilamentBambuPVA
	FilamentBambuPLACF
	FilamentBambuSupportForPLAPETG
	FilamentBambuTPU95AHF
	FilamentBambuPPACF
	FilamentBambuASA
	FilamentBambuPLAGlow
	FilamentABS
	FilamentGenericASA
	FilamentPA
	FilamentPACF
	FilamentPC
	FilamentPETG
	FilamentPLA
	FilamentPLACF
	FilamentGenericPVA
	FilamentTPU
)

// FilamentByName returns a Filament by its name.
func FilamentByName(name string) (Filament, error) {
	filaments := map[string]Filament{
		"POLYLITE_PLA":               FilamentPolyLitePLA,
		"POLYTERRA_PLA":              FilamentPolyTerraPLA,
		"BAMBU_ABS":                  FilamentBambuABS,
		"BAMBU_PA_CF":                FilamentBambuPACF,
		"BAMBU_PC":                   FilamentBambuPC,
		"BAMBU_PLA_Basic":            FilamentBambuPLABasic,
		"BAMBU_PLA_Matte":            FilamentBambuPLAMatte,
		"SUPPORT_G":                  FilamentSupportG,
		"SUPPORT_W":                  FilamentSupportW,
		"BAMBU_TPU_95A":              FilamentBambuTPU95A,
		"BAMBU_ASA_AERO":             FilamentBambuASAAero,
		"BAMBU_PLA_METAL":            FilamentBambuPLAMetal,
		"BAMBU_PETG_TRANSLUCENT":     FilamentBambuPETGTranslucent,
		"BAMBU_PLA_MARBLE":           FilamentBambuPLAMarble,
		"BAMBU_PLA_WOOD":             FilamentBambuPLAWood,
		"BAMBU_PLA_SILK_PLUS":        FilamentBambuPLASilkPlus,
		"BAMBU_PETG_HF":              FilamentBambuPETGHF,
		"BAMBU_TPU_FOR_AMS":          FilamentBambuTPUForAMS,
		"BAMBU_SUPPORT_FOR_ABS":      FilamentBambuSupportForABS,
		"BAMBU_PC_FR":                FilamentBambuPCFR,
		"BAMBU_PLA_GALAXY":           FilamentBambuPLAGalaxy,
		"BAMBU_PA6_GF":               FilamentBambuPA6GF,
		"BAMBU_PLA_AERO":             FilamentBambuPLAAero,
		"BAMBU_ASA_CF":               FilamentBambuASACF,
		"BAMBU_PETG_CF":              FilamentBambuPETGCF,
		"BAMBU_SUPPORT_FOR_PA_PET":   FilamentBambuSupportForPAPET,
		"BAMBU_PLA_SPARKLE":          FilamentBambuPLASparkle,
		"BAMBU_ABS_GF":               FilamentBambuABSGF,
		"BAMBU_PAHT_CF":              FilamentBambuPAHTCF,
		"BAMBU_PLA_BASIC":            FilamentBambuPLABasic,
		"BAMBU_PLA_MATTE":            FilamentBambuPLAMatte,
		"BAMBU_PA6_CF":               FilamentBambuPA6CF,
		"BAMBU_PLA_SILK":             FilamentBambuPLASilk,
		"BAMBU_PVA":                  FilamentBambuPVA,
		"BAMBU_PLA_CF":               FilamentBambuPLACF,
		"BAMBU_SUPPORT_FOR_PLA_PETG": FilamentBambuSupportForPLAPETG,
		"BAMBU_TPU_95A_HF":           FilamentBambuTPU95AHF,
		"BAMBU_PPA_CF":               FilamentBambuPPACF,
		"BAMBU_ASA":                  FilamentBambuASA,
		"BAMBU_PLA_GLOW":             FilamentBambuPLAGlow,
		"ABS":                        FilamentABS,
		"ASA":                        FilamentGenericASA,
		"PA":                         FilamentPA,
		"PA_CF":                      FilamentPACF,
		"PC":                         FilamentPC,
		"PETG":                       FilamentPETG,
		"PLA":                        FilamentPLA,
		"PLA_CF":                     FilamentPLACF,
		"PVA":                        FilamentGenericPVA,
		"TPU":                        FilamentTPU,
	}

	if f, ok := filaments[name]; ok {
		return f, nil
	}
	return 0, fmt.Errorf("filament %s not found", name)
}

// FilamentByIndex returns a Filament by its tray_info_idx code (e.g., "GFG99").
func FilamentByIndex(idx string) (Filament, error) {
	idxToFilament := map[string]Filament{
		"GFL00": FilamentPolyLitePLA,
		"GFL01": FilamentPolyTerraPLA,
		"GFB00": FilamentBambuABS,
		"GFN03": FilamentBambuPACF,
		"GFC00": FilamentBambuPC,
		"GFA00": FilamentBambuPLABasic,
		"GFA01": FilamentBambuPLAMatte,
		"GFS01": FilamentSupportG,
		"GFS00": FilamentSupportW,
		"GFU01": FilamentBambuTPU95A,
		"GFB02": FilamentBambuASAAero,
		"GFA02": FilamentBambuPLAMetal,
		"GFG01": FilamentBambuPETGTranslucent,
		"GFA07": FilamentBambuPLAMarble,
		"GFA16": FilamentBambuPLAWood,
		"GFA06": FilamentBambuPLASilkPlus,
		"GFG02": FilamentBambuPETGHF,
		"GFU02": FilamentBambuTPUForAMS,
		"GFS06": FilamentBambuSupportForABS,
		"GFC01": FilamentBambuPCFR,
		"GFA15": FilamentBambuPLAGalaxy,
		"GFN08": FilamentBambuPA6GF,
		"GFA11": FilamentBambuPLAAero,
		"GFB51": FilamentBambuASACF,
		"GFG50": FilamentBambuPETGCF,
		"GFS03": FilamentBambuSupportForPAPET,
		"GFA08": FilamentBambuPLASparkle,
		"GFB50": FilamentBambuABSGF,
		"GFN04": FilamentBambuPAHTCF,
		"GFN05": FilamentBambuPA6CF,
		"GFA05": FilamentBambuPLASilk,
		"GFS04": FilamentBambuPVA,
		"GFA50": FilamentBambuPLACF,
		"GFS05": FilamentBambuSupportForPLAPETG,
		"GFU00": FilamentBambuTPU95AHF,
		"GFN06": FilamentBambuPPACF,
		"GFB01": FilamentBambuASA,
		"GFA12": FilamentBambuPLAGlow,
		"GFB99": FilamentABS,
		"GFB98": FilamentGenericASA,
		"GFN99": FilamentPA,
		"GFN98": FilamentPACF,
		"GFC99": FilamentPC,
		"GFG99": FilamentPETG,
		"GFL99": FilamentPLA,
		"GFL98": FilamentPLACF,
		"GFS99": FilamentGenericPVA,
		"GFU99": FilamentTPU,
	}

	if f, ok := idxToFilament[idx]; ok {
		return f, nil
	}
	return 0, fmt.Errorf("filament index %s not found", idx)
}

// Name returns the human-readable name of a Filament.
func (f Filament) Name() string {
	names := map[Filament]string{
		FilamentPolyLitePLA:            "PolyLite PLA",
		FilamentPolyTerraPLA:           "PolyTerra PLA",
		FilamentBambuABS:               "Bambu ABS",
		FilamentBambuPACF:              "Bambu PA-CF",
		FilamentBambuPC:                "Bambu PC",
		FilamentBambuPLABasic:          "Bambu PLA Basic",
		FilamentBambuPLAMatte:          "Bambu PLA Matte",
		FilamentSupportG:               "Support G",
		FilamentSupportW:               "Support W",
		FilamentBambuTPU95A:            "Bambu TPU 95A",
		FilamentBambuASAAero:           "Bambu ASA Aero",
		FilamentBambuPLAMetal:          "Bambu PLA Metal",
		FilamentBambuPETGTranslucent:   "Bambu PETG Translucent",
		FilamentBambuPLAMarble:         "Bambu PLA Marble",
		FilamentBambuPLAWood:           "Bambu PLA Wood",
		FilamentBambuPLASilkPlus:       "Bambu PLA Silk Plus",
		FilamentBambuPETGHF:            "Bambu PETG HF",
		FilamentBambuTPUForAMS:         "Bambu TPU for AMS",
		FilamentBambuSupportForABS:     "Bambu Support for ABS",
		FilamentBambuPCFR:              "Bambu PC-FR",
		FilamentBambuPLAGalaxy:         "Bambu PLA Galaxy",
		FilamentBambuPA6GF:             "Bambu PA6-GF",
		FilamentBambuPLAAero:           "Bambu PLA Aero",
		FilamentBambuASACF:             "Bambu ASA-CF",
		FilamentBambuPETGCF:            "Bambu PETG-CF",
		FilamentBambuSupportForPAPET:   "Bambu Support for PA/PET",
		FilamentBambuPLASparkle:        "Bambu PLA Sparkle",
		FilamentBambuABSGF:             "Bambu ABS-GF",
		FilamentBambuPAHTCF:            "Bambu PAHT-CF",
		FilamentBambuPA6CF:             "Bambu PA6-CF",
		FilamentBambuPLASilk:           "Bambu PLA Silk",
		FilamentBambuPVA:               "Bambu PVA",
		FilamentBambuPLACF:             "Bambu PLA-CF",
		FilamentBambuSupportForPLAPETG: "Bambu Support for PLA/PETG",
		FilamentBambuTPU95AHF:          "Bambu TPU 95A HF",
		FilamentBambuPPACF:             "Bambu PPA-CF",
		FilamentBambuASA:               "Bambu ASA",
		FilamentBambuPLAGlow:           "Bambu PLA Glow",
		FilamentABS:                    "ABS",
		FilamentGenericASA:             "ASA",
		FilamentPA:                     "PA",
		FilamentPACF:                   "PA-CF",
		FilamentPC:                     "PC",
		FilamentPETG:                   "PETG",
		FilamentPLA:                    "PLA",
		FilamentPLACF:                  "PLA-CF",
		FilamentGenericPVA:             "PVA",
		FilamentTPU:                    "TPU",
	}

	if name, ok := names[f]; ok {
		return name
	}
	return "Unknown"
}

// GetSettings returns the AMSFilamentSettings for a Filament.
func (f Filament) GetSettings() AMSFilamentSettings {
	settings := map[Filament]AMSFilamentSettings{
		FilamentPolyLitePLA:            {"GFL00", 190, 250, "PLA"},
		FilamentPolyTerraPLA:           {"GFL01", 190, 250, "PLA"},
		FilamentBambuABS:               {"GFB00", 240, 270, "ABS"},
		FilamentBambuPACF:              {"GFN03", 270, 300, "PA-CF"},
		FilamentBambuPC:                {"GFC00", 260, 280, "PC"},
		FilamentBambuPLABasic:          {"GFA00", 190, 250, "PLA"},
		FilamentBambuPLAMatte:          {"GFA01", 190, 250, "PLA"},
		FilamentSupportG:               {"GFS01", 190, 250, "PA-S"},
		FilamentSupportW:               {"GFS00", 190, 250, "PLA-S"},
		FilamentBambuTPU95A:            {"GFU01", 200, 250, "TPU"},
		FilamentBambuASAAero:           {"GFB02", 240, 280, "ASA"},
		FilamentBambuPLAMetal:          {"GFA02", 190, 230, "PLA"},
		FilamentBambuPETGTranslucent:   {"GFG01", 230, 260, "PETG"},
		FilamentBambuPLAMarble:         {"GFA07", 190, 230, "PLA"},
		FilamentBambuPLAWood:           {"GFA16", 190, 240, "PLA"},
		FilamentBambuPLASilkPlus:       {"GFA06", 210, 240, "PLA"},
		FilamentBambuPETGHF:            {"GFG02", 230, 260, "PETG"},
		FilamentBambuTPUForAMS:         {"GFU02", 230, 230, "TPU"},
		FilamentBambuSupportForABS:     {"GFS06", 190, 220, "Support"},
		FilamentBambuPCFR:              {"GFC01", 260, 280, "PC"},
		FilamentBambuPLAGalaxy:         {"GFA15", 190, 230, "PLA"},
		FilamentBambuPA6GF:             {"GFN08", 260, 290, "PA6"},
		FilamentBambuPLAAero:           {"GFA11", 220, 260, "PLA"},
		FilamentBambuASACF:             {"GFB51", 250, 280, "ASA"},
		FilamentBambuPETGCF:            {"GFG50", 240, 270, "PETG"},
		FilamentBambuSupportForPAPET:   {"GFS03", 280, 300, "Support"},
		FilamentBambuPLASparkle:        {"GFA08", 190, 230, "PLA"},
		FilamentBambuABSGF:             {"GFB50", 240, 270, "ABS"},
		FilamentBambuPAHTCF:            {"GFN04", 260, 290, "PAHT"},
		FilamentBambuPA6CF:             {"GFN05", 260, 290, "PA6"},
		FilamentBambuPLASilk:           {"GFA05", 210, 230, "PLA"},
		FilamentBambuPVA:               {"GFS04", 220, 250, "PVA"},
		FilamentBambuPLACF:             {"GFA50", 210, 240, "PLA"},
		FilamentBambuSupportForPLAPETG: {"GFS05", 190, 220, "Support"},
		FilamentBambuTPU95AHF:          {"GFU00", 230, 230, "TPU"},
		FilamentBambuPPACF:             {"GFN06", 280, 310, "PPA"},
		FilamentBambuASA:               {"GFB01", 240, 270, "ASA"},
		FilamentBambuPLAGlow:           {"GFA12", 190, 230, "PLA"},
		FilamentABS:                    {"GFB99", 240, 270, "ABS"},
		FilamentGenericASA:             {"GFB98", 240, 270, "ASA"},
		FilamentPA:                     {"GFN99", 270, 300, "PA"},
		FilamentPACF:                   {"GFN98", 270, 300, "PA"},
		FilamentPC:                     {"GFC99", 260, 280, "PC"},
		FilamentPETG:                   {"GFG99", 220, 260, "PETG"},
		FilamentPLA:                    {"GFL99", 190, 250, "PLA"},
		FilamentPLACF:                  {"GFL98", 190, 250, "PLA"},
		FilamentGenericPVA:             {"GFS99", 190, 250, "PVA"},
		FilamentTPU:                    {"GFU99", 200, 250, "TPU"},
	}

	if s, ok := settings[f]; ok {
		return s
	}
	return AMSFilamentSettings{}
}

// FilamentTray represents a filament tray.
type FilamentTray struct {
	K             float64  `json:"k"`
	N             int      `json:"n"`
	TagUID        string   `json:"tag_uid"`
	TrayIDName    string   `json:"tray_id_name"`
	TrayInfoIdx   string   `json:"tray_info_idx"`
	TrayType      string   `json:"tray_type"`
	TraySubBrands string   `json:"tray_sub_brands"`
	TrayColor     string   `json:"tray_color"`
	TrayWeight    string   `json:"tray_weight"`
	TrayDiameter  string   `json:"tray_diameter"`
	TrayTemp      string   `json:"tray_temp"`
	TrayTime      string   `json:"tray_time"`
	BedTempType   string   `json:"bed_temp_type"`
	BedTemp       string   `json:"bed_temp"`
	NozzleTempMax int      `json:"nozzle_temp_max"`
	NozzleTempMin int      `json:"nozzle_temp_min"`
	XCamInfo      string   `json:"xcam_info"`
	TrayUUID      string   `json:"tray_uuid"`
	Cols          []string `json:"cols,omitempty"`
}

// FilamentTrayFromDict creates a FilamentTray from a map.
func FilamentTrayFromDict(d map[string]any) FilamentTray {
	tray := FilamentTray{}

	if v, ok := d["k"]; ok {
		switch val := v.(type) {
		case float64:
			tray.K = val
		case int:
			tray.K = float64(val)
		}
	}
	if v, ok := d["n"]; ok {
		switch val := v.(type) {
		case float64:
			tray.N = int(val)
		case int:
			tray.N = val
		}
	}
	if v, ok := d["tag_uid"]; ok {
		tray.TagUID = fmt.Sprintf("%v", v)
	}
	if v, ok := d["tray_id_name"]; ok {
		tray.TrayIDName = fmt.Sprintf("%v", v)
	}
	if v, ok := d["tray_info_idx"]; ok {
		tray.TrayInfoIdx = fmt.Sprintf("%v", v)
	}
	if v, ok := d["tray_type"]; ok {
		tray.TrayType = fmt.Sprintf("%v", v)
	}
	if v, ok := d["tray_sub_brands"]; ok {
		tray.TraySubBrands = fmt.Sprintf("%v", v)
	}
	if v, ok := d["tray_color"]; ok {
		tray.TrayColor = fmt.Sprintf("%v", v)
	}
	if v, ok := d["tray_weight"]; ok {
		tray.TrayWeight = fmt.Sprintf("%v", v)
	}
	if v, ok := d["tray_diameter"]; ok {
		tray.TrayDiameter = fmt.Sprintf("%v", v)
	}
	if v, ok := d["tray_temp"]; ok {
		tray.TrayTemp = fmt.Sprintf("%v", v)
	}
	if v, ok := d["tray_time"]; ok {
		tray.TrayTime = fmt.Sprintf("%v", v)
	}
	if v, ok := d["bed_temp_type"]; ok {
		tray.BedTempType = fmt.Sprintf("%v", v)
	}
	if v, ok := d["bed_temp"]; ok {
		tray.BedTemp = fmt.Sprintf("%v", v)
	}
	if v, ok := d["nozzle_temp_max"]; ok {
		switch val := v.(type) {
		case float64:
			tray.NozzleTempMax = int(val)
		case int:
			tray.NozzleTempMax = val
		}
	}
	if v, ok := d["nozzle_temp_min"]; ok {
		switch val := v.(type) {
		case float64:
			tray.NozzleTempMin = int(val)
		case int:
			tray.NozzleTempMin = val
		}
	}
	if v, ok := d["xcam_info"]; ok {
		tray.XCamInfo = fmt.Sprintf("%v", v)
	}
	if v, ok := d["tray_uuid"]; ok {
		tray.TrayUUID = fmt.Sprintf("%v", v)
	}
	if v, ok := d["cols"]; ok {
		if cols, ok := v.([]any); ok {
			tray.Cols = make([]string, len(cols))
			for i, col := range cols {
				tray.Cols[i] = fmt.Sprintf("%v", col)
			}
		}
	}

	return tray
}

// FilamentTrayFromJSON creates a FilamentTray from JSON.
func FilamentTrayFromJSON(data []byte) (FilamentTray, error) {
	var d map[string]any
	if err := json.Unmarshal(data, &d); err != nil {
		return FilamentTray{}, err
	}
	return FilamentTrayFromDict(d), nil
}

// GetFilament returns the Filament for this tray.
func (f *FilamentTray) GetFilament() AMSFilamentSettings {
	return AMSFilamentSettings{
		TrayInfoIdx:   f.TrayInfoIdx,
		NozzleTempMin: f.NozzleTempMin,
		NozzleTempMax: f.NozzleTempMax,
		TrayType:      f.TrayType,
	}
}
