package domain

import "encoding/xml"

type portGroup struct {
	Text string   `xml:",chardata"`
	ID   string   `xml:"id"`
	Port []uint16 `xml:"port"`
}
type portGroups struct {
	Text      string      `xml:",chardata"`
	ID        string      `xml:"id"`
	PortGroup []portGroup `xml:"PortGroup"`
}

func (p portGroups) GroupById(id string) []uint16 {
	for _, group := range p.PortGroup {
		if group.ID == id {
			return group.Port
		}
	}
	panic("group not found")
}

type rule_to struct {
	Text  string `xml:",chardata"`
	Group string `xml:"group"`
	Port  uint16 `xml:"port"`
}
type rule struct {
	Text string  `xml:",chardata"`
	ID   string  `xml:"id"`
	From string  `xml:"from"`
	To   rule_to `xml:"to"`
}
type rules struct {
	Text string `xml:",chardata"`
	ID   string `xml:"id"`
	Rule []rule `xml:"Rule"`
}

type LoadBalancerConfig struct {
	XMLName    xml.Name   `xml:"LoadBalancerConfig"`
	Text       string     `xml:",chardata"`
	Xmlns      string     `xml:"xmlns,attr"`
	ID         string     `xml:"id"`
	PortGroups portGroups `xml:"PortGroups"`
	Rules      rules      `xml:"Rules"`
}
