package snmp

// OIDs padrão MIB-II (RFC 1213) e extensões de fabricantes
// Referência: https://www.iana.org/assignments/enterprise-numbers

// ---- Sistema (MIB-II) ----
const (
	OIDSysDescr    = "1.3.6.1.2.1.1.1.0"
	OIDSysOID      = "1.3.6.1.2.1.1.2.0"
	OIDSysUptime   = "1.3.6.1.2.1.1.3.0"
	OIDSysContact  = "1.3.6.1.2.1.1.4.0"
	OIDSysName     = "1.3.6.1.2.1.1.5.0"
	OIDSysLocation = "1.3.6.1.2.1.1.6.0"
)

// ---- Interfaces (IF-MIB) ----
const (
	OIDIfNumber      = "1.3.6.1.2.1.2.1.0"
	OIDIfTable       = "1.3.6.1.2.1.2.2"
	OIDIfIndex       = "1.3.6.1.2.1.2.2.1.1"
	OIDIfDescr       = "1.3.6.1.2.1.2.2.1.2"
	OIDIfType        = "1.3.6.1.2.1.2.2.1.3"
	OIDIfSpeed       = "1.3.6.1.2.1.2.2.1.5"
	OIDIfAdminStatus = "1.3.6.1.2.1.2.2.1.7"
	OIDIfOperStatus  = "1.3.6.1.2.1.2.2.1.8"
	OIDIfInOctets    = "1.3.6.1.2.1.2.2.1.10"
	OIDIfInErrors    = "1.3.6.1.2.1.2.2.1.14"
	OIDIfOutOctets   = "1.3.6.1.2.1.2.2.1.16"
	OIDIfOutErrors   = "1.3.6.1.2.1.2.2.1.20"

	// IF-MIB 64-bit counters (RFC 2863)
	OIDIfHCInOctets  = "1.3.6.1.2.1.31.1.1.1.6"
	OIDIfHCOutOctets = "1.3.6.1.2.1.31.1.1.1.10"
	OIDIfAlias       = "1.3.6.1.2.1.31.1.1.1.18"
	OIDIfHighSpeed   = "1.3.6.1.2.1.31.1.1.1.15" // Mbps
)

// ---- IP (MIB-II) ----
const (
	OIDIPAddrTable  = "1.3.6.1.2.1.4.20"
	OIDIPAdEntAddr  = "1.3.6.1.2.1.4.20.1.1"
	OIDIPAdEntIfIdx = "1.3.6.1.2.1.4.20.1.2"
)

// ---- HOST-RESOURCES-MIB (RFC 2790) ----
const (
	// CPU
	OIDHRProcessorLoad = "1.3.6.1.2.1.25.3.3.1.2"

	// Memória
	OIDHRStorageTable       = "1.3.6.1.2.1.25.2.3"
	OIDHRStorageDescr       = "1.3.6.1.2.1.25.2.3.1.3"
	OIDHRStorageAllocationUnits = "1.3.6.1.2.1.25.2.3.1.4"
	OIDHRStorageSize        = "1.3.6.1.2.1.25.2.3.1.5"
	OIDHRStorageUsed        = "1.3.6.1.2.1.25.2.3.1.6"
)

// ---- UCD-SNMP-MIB (Linux/Net-SNMP) ----
const (
	OIDUCDCPUUser    = "1.3.6.1.4.1.2021.11.9.0"
	OIDUCDCPUSystem  = "1.3.6.1.4.1.2021.11.10.0"
	OIDUCDCPUIdle    = "1.3.6.1.4.1.2021.11.11.0"
	OIDUCDMemTotal   = "1.3.6.1.4.1.2021.4.5.0"  // kB
	OIDUCDMemFree    = "1.3.6.1.4.1.2021.4.11.0" // kB
	OIDUCDMemBuffers = "1.3.6.1.4.1.2021.4.14.0"
	OIDUCDMemCached  = "1.3.6.1.4.1.2021.4.15.0"
)

// ---- Mikrotik (enterprise 14988) ----
const (
	OIDMikrotikCPU         = "1.3.6.1.4.1.14988.1.1.3.14.0" // CPU load %
	OIDMikrotikMemTotal    = "1.3.6.1.4.1.14988.1.1.8.1.0"  // bytes
	OIDMikrotikMemUsed     = "1.3.6.1.4.1.14988.1.1.8.2.0"  // bytes
	OIDMikrotikTemperature = "1.3.6.1.4.1.14988.1.1.3.10.0" // °C * 10
	OIDMikrotikFirmware    = "1.3.6.1.4.1.14988.1.1.7.4.0"
	OIDMikrotikModel       = "1.3.6.1.4.1.14988.1.1.7.7.0"
	OIDMikrotikSerial      = "1.3.6.1.4.1.14988.1.1.7.3.0"
	OIDMikrotikVoltage     = "1.3.6.1.4.1.14988.1.1.3.8.0"  // mV
)

// ---- Enterprise OIDs (para identificar fabricante) ----
var EnterpriseOIDs = map[string]string{
	"1.3.6.1.4.1.14988": "mikrotik",
	"1.3.6.1.4.1.9":     "cisco",
	"1.3.6.1.4.1.2011":  "huawei",
	"1.3.6.1.4.1.2636":  "juniper",
	"1.3.6.1.4.1.41112": "ubiquiti",
}

// DetectVendorFromOID identifica o fabricante a partir do sysOID
func DetectVendorFromOID(sysOID string) string {
	for prefix, vendor := range EnterpriseOIDs {
		if len(sysOID) >= len(prefix) && sysOID[:len(prefix)] == prefix {
			return vendor
		}
	}
	return "generic"
}

// SysInfoOIDs são os OIDs coletados em todo GET inicial
var SysInfoOIDs = []string{
	OIDSysDescr,
	OIDSysOID,
	OIDSysUptime,
	OIDSysContact,
	OIDSysName,
	OIDSysLocation,
}

// MetricOIDsByVendor mapeia tipo de métrica para OID por fabricante
var MetricOIDsByVendor = map[string]map[string]string{
	"mikrotik": {
		"cpu":         OIDMikrotikCPU,
		"memory_used": OIDMikrotikMemUsed,
		"memory_total": OIDMikrotikMemTotal,
		"temperature": OIDMikrotikTemperature,
	},
	"generic": {
		"cpu_user":   OIDUCDCPUUser,
		"cpu_system": OIDUCDCPUSystem,
		"cpu_idle":   OIDUCDCPUIdle,
		"mem_total":  OIDUCDMemTotal,
		"mem_free":   OIDUCDMemFree,
	},
}
