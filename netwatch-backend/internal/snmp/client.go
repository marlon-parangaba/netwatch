package snmp

import (
	"fmt"
	"time"

	"github.com/ertel/netwatch-backend/internal/models"
	"github.com/gosnmp/gosnmp"
)

// Client é um wrapper sobre gosnmp adaptado ao modelo de Device
type Client struct {
	snmp    *gosnmp.GoSNMP
	device  *models.Device
}

// Result representa um único resultado OID → valor
type Result struct {
	OID   string
	Value interface{}
	Type  gosnmp.Asn1BER
}

// NewClient cria um cliente SNMP a partir de um Device
func NewClient(device *models.Device, timeout time.Duration, retries int) (*Client, error) {
	g := &gosnmp.GoSNMP{
		Target:    device.IPAddress,
		Port:      uint16(device.SNMPPort),
		Timeout:   timeout,
		Retries:   retries,
		MaxOids:   gosnmp.MaxOids,
	}

	switch device.SNMPVersion {
	case models.SNMPv1:
		g.Version = gosnmp.Version1
		g.Community = device.SNMPCommunity
	case models.SNMPv2c:
		g.Version = gosnmp.Version2c
		g.Community = device.SNMPCommunity
	case models.SNMPv3:
		g.Version = gosnmp.Version3
		g.SecurityModel = gosnmp.UserSecurityModel
		g.MsgFlags = gosnmp.AuthPriv
		g.SecurityParameters = &gosnmp.UsmSecurityParameters{
			UserName:                 device.SNMPv3Username,
			AuthenticationProtocol:   authProtoFromString(device.SNMPv3AuthProto),
			AuthenticationPassphrase: device.SNMPv3AuthKey,
			PrivacyProtocol:          privProtoFromString(device.SNMPv3PrivProto),
			PrivacyPassphrase:        device.SNMPv3PrivKey,
		}
	default:
		return nil, fmt.Errorf("versão SNMP não suportada: %s", device.SNMPVersion)
	}

	return &Client{snmp: g, device: device}, nil
}

// Connect abre a sessão UDP
func (c *Client) Connect() error {
	return c.snmp.Connect()
}

// Close encerra a sessão
func (c *Client) Close() {
	if c.snmp.Conn != nil {
		c.snmp.Conn.Close()
	}
}

// Get busca um ou mais OIDs específicos
func (c *Client) Get(oids []string) ([]Result, error) {
	pkt, err := c.snmp.Get(oids)
	if err != nil {
		return nil, fmt.Errorf("SNMP GET falhou: %w", err)
	}
	return packetToResults(pkt), nil
}

// GetSingle busca um único OID e retorna o valor como float64
func (c *Client) GetSingle(oid string) (float64, error) {
	results, err := c.Get([]string{oid})
	if err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, fmt.Errorf("OID %s não retornou valor", oid)
	}
	return toFloat64(results[0].Value), nil
}

// GetString busca um único OID e retorna como string
func (c *Client) GetString(oid string) (string, error) {
	results, err := c.Get([]string{oid})
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", fmt.Errorf("OID %s não retornou valor", oid)
	}
	return toString(results[0].Value), nil
}

// Walk faz um SNMP Walk a partir de um OID raiz
func (c *Client) Walk(rootOID string) ([]Result, error) {
	var results []Result

	walkFn := func(pdu gosnmp.SnmpPDU) error {
		results = append(results, Result{
			OID:   pdu.Name,
			Value: pdu.Value,
			Type:  pdu.Type,
		})
		return nil
	}

	var err error
	if c.snmp.Version == gosnmp.Version1 {
		err = c.snmp.Walk(rootOID, walkFn)
	} else {
		err = c.snmp.BulkWalk(rootOID, walkFn)
	}

	if err != nil {
		return nil, fmt.Errorf("SNMP Walk falhou em %s: %w", rootOID, err)
	}
	return results, nil
}

// GetSysInfo coleta informações básicas do sistema
func (c *Client) GetSysInfo() (*SysInfo, error) {
	results, err := c.Get(SysInfoOIDs)
	if err != nil {
		return nil, err
	}

	info := &SysInfo{}
	for _, r := range results {
		switch r.OID {
		case OIDSysDescr:
			info.Descr = toString(r.Value)
		case OIDSysOID:
			info.OID = toString(r.Value)
		case OIDSysUptime:
			info.Uptime = toFloat64(r.Value)
		case OIDSysContact:
			info.Contact = toString(r.Value)
		case OIDSysName:
			info.Name = toString(r.Value)
		case OIDSysLocation:
			info.Location = toString(r.Value)
		}
	}

	info.Vendor = DetectVendorFromOID(info.OID)
	return info, nil
}

// GetInterfaces retorna a tabela de interfaces
func (c *Client) GetInterfaces() ([]InterfaceData, error) {
	ifMap := make(map[int]*InterfaceData)

	walkAndFill := func(oid string, setter func(*InterfaceData, gosnmp.SnmpPDU)) error {
		results, err := c.Walk(oid)
		if err != nil {
			return err
		}
		for _, r := range results {
			idx := lastOIDIndex(r.OID)
			if _, ok := ifMap[idx]; !ok {
				ifMap[idx] = &InterfaceData{Index: idx}
			}
			setter(ifMap[idx], gosnmp.SnmpPDU{Name: r.OID, Value: r.Value, Type: r.Type})
		}
		return nil
	}

	_ = walkAndFill(OIDIfDescr, func(d *InterfaceData, p gosnmp.SnmpPDU) { d.Descr = toString(p.Value) })
	_ = walkAndFill(OIDIfType, func(d *InterfaceData, p gosnmp.SnmpPDU) { d.Type = int(toFloat64(p.Value)) })
	_ = walkAndFill(OIDIfAdminStatus, func(d *InterfaceData, p gosnmp.SnmpPDU) { d.AdminStatus = int(toFloat64(p.Value)) })
	_ = walkAndFill(OIDIfOperStatus, func(d *InterfaceData, p gosnmp.SnmpPDU) { d.OperStatus = int(toFloat64(p.Value)) })
	_ = walkAndFill(OIDIfAlias, func(d *InterfaceData, p gosnmp.SnmpPDU) { d.Alias = toString(p.Value) })

	// Tentar 64-bit primeiro, fallback para 32-bit
	err64In := walkAndFill(OIDIfHCInOctets, func(d *InterfaceData, p gosnmp.SnmpPDU) { d.InOctets = uint64(toFloat64(p.Value)) })
	if err64In != nil {
		_ = walkAndFill(OIDIfInOctets, func(d *InterfaceData, p gosnmp.SnmpPDU) { d.InOctets = uint64(toFloat64(p.Value)) })
	}
	err64Out := walkAndFill(OIDIfHCOutOctets, func(d *InterfaceData, p gosnmp.SnmpPDU) { d.OutOctets = uint64(toFloat64(p.Value)) })
	if err64Out != nil {
		_ = walkAndFill(OIDIfOutOctets, func(d *InterfaceData, p gosnmp.SnmpPDU) { d.OutOctets = uint64(toFloat64(p.Value)) })
	}
	_ = walkAndFill(OIDIfHighSpeed, func(d *InterfaceData, p gosnmp.SnmpPDU) { d.SpeedMbps = uint64(toFloat64(p.Value)) })
	_ = walkAndFill(OIDIfInErrors, func(d *InterfaceData, p gosnmp.SnmpPDU) { d.InErrors = uint64(toFloat64(p.Value)) })
	_ = walkAndFill(OIDIfOutErrors, func(d *InterfaceData, p gosnmp.SnmpPDU) { d.OutErrors = uint64(toFloat64(p.Value)) })

	var ifaces []InterfaceData
	for _, v := range ifMap {
		ifaces = append(ifaces, *v)
	}
	return ifaces, nil
}

// ---- Tipos de dados ----

type SysInfo struct {
	Descr    string
	OID      string
	Uptime   float64 // centiseconds
	Contact  string
	Name     string
	Location string
	Vendor   string
}

type InterfaceData struct {
	Index       int
	Descr       string
	Type        int
	SpeedMbps   uint64
	AdminStatus int
	OperStatus  int
	Alias       string
	InOctets    uint64
	OutOctets   uint64
	InErrors    uint64
	OutErrors   uint64
}

// ---- Helpers internos ----

func packetToResults(pkt *gosnmp.SnmpPacket) []Result {
	results := make([]Result, 0, len(pkt.Variables))
	for _, v := range pkt.Variables {
		results = append(results, Result{
			OID:   v.Name,
			Value: v.Value,
			Type:  v.Type,
		})
	}
	return results
}

func toFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case uint:
		return float64(val)
	case uint32:
		return float64(val)
	case uint64:
		return float64(val)
	case int:
		return float64(val)
	case int32:
		return float64(val)
	case int64:
		return float64(val)
	case float32:
		return float64(val)
	case float64:
		return val
	}
	return 0
}

func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case []byte:
		return string(val)
	case gosnmp.ObjectIdentifier:
		return val.String()
	}
	return fmt.Sprintf("%v", v)
}

// lastOIDIndex extrai o último número de um OID (ex: .1.3.6.1...1.5 → 5)
func lastOIDIndex(oid string) int {
	var idx int
	for i := len(oid) - 1; i >= 0; i-- {
		if oid[i] == '.' {
			fmt.Sscanf(oid[i+1:], "%d", &idx)
			return idx
		}
	}
	return 0
}

func authProtoFromString(s string) gosnmp.SnmpV3AuthProtocol {
	switch s {
	case "SHA":
		return gosnmp.SHA
	case "SHA224":
		return gosnmp.SHA224
	case "SHA256":
		return gosnmp.SHA256
	case "SHA384":
		return gosnmp.SHA384
	case "SHA512":
		return gosnmp.SHA512
	default:
		return gosnmp.MD5
	}
}

func privProtoFromString(s string) gosnmp.SnmpV3PrivProtocol {
	switch s {
	case "AES":
		return gosnmp.AES
	case "AES192":
		return gosnmp.AES192
	case "AES256":
		return gosnmp.AES256
	case "DES":
		return gosnmp.DES
	default:
		return gosnmp.AES
	}
}
