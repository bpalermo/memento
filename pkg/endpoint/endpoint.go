package endpoint

import "net"

// Endpoint service endpoint
type Endpoint struct {
	Ip       string   `json:"ip,omitempty"`
	Port     uint16   `json:"port,omitempty"`
	Locality string   `json:"locality,omitempty"`
	Stage    string   `json:"stage,omitempty"`
	Metadata []string `json:"metadata,omitempty"`
}

func NewEndpoint(ip net.IP, port uint16) *Endpoint {
	return &Endpoint{
		Ip:   ip.String(),
		Port: port,
	}
}
