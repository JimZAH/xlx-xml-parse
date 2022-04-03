package xlxxml

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type Station struct {
	Stationx  []Stations `xml:"STATION"`
	Timestamp int64
}

type Nodes struct {
	Nodex []Node `xml:"NODE"`
}

type Node struct {
	Callsign      string `xml:"Callsign"`
	IP            string `xml:"IP"`
	LinkedModule  string `xml:"LinkedModule"`
	Protocol      string `xml:"Protocol"`
	ConnectTime   string `xml:"ConnectTime"`
	LastHeardTime string `xml:"LastHeardTime"`
}

type Stations struct {
	XMLName       xml.Name `xml:"STATION"`
	Callsign      string   `xml:"Callsign"`
	Vianode       string   `xml:"Via-node"`
	Onmodule      string   `xml:"On-module"`
	Viapeer       string   `xml:"Via-peer"`
	LastHeardTime string   `xml:"LastHeardTime"`
}

func Parse(file string) (Nodes, Station) {

	var nodes Nodes
	var stations Station

	// Line counter
	var lc int = 0

	// Station buffer
	var buffer2 []byte

	xmlFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}

	// Main buffer
	byteValue, _ := ioutil.ReadAll(xmlFile)

	// Clean up the XLX XML mess
	for i := 0; i < len(byteValue); i++ {

		// Remove first two lines, replace with ?
		if byteValue[i] == 10 && lc < 2 {
			lc++
			for j := 0; j < i; j++ {
				byteValue[j] = 63
			}
			continue
		}

		// Remove any spaces in tags
		if byteValue[i] == 60 && lc > 1 {
			lc++
			for k := i + 1; true; k++ {
				if byteValue[k] == 32 {
					byteValue[k] = 45
				}
				if byteValue[k] == 62 {
					break
				}
			}
		}

		// Remove <XLXWVV  linked peers> we don't care about those
		if byteValue[i] == 60 && byteValue[i+1] == 88 && byteValue[i+16] == 112 {
			for l := i; true; l++ {
				byteValue[l] = 63
				if byteValue[l+1] == 60 && byteValue[l+2] == 88 {
					break
				}
			}
		}
	}

	// Store station tags in buffer 2
	for i := 0; i < len(byteValue); i++ {
		if byteValue[i] == 60 {
			if byteValue[i+1] == 88 && byteValue[i+2] == 76 && byteValue[i+15] == 117 {
				for l := i; l < len(byteValue); l++ {
					buffer2 = append(buffer2, byteValue[l])
				}
			}
		}
	}

	err = xml.Unmarshal(byteValue, &nodes)
	if err != nil {
		panic(err)
	}

	err = xml.Unmarshal(buffer2, &stations)
	if err != nil {
		panic(err)
	}

	defer xmlFile.Close()

	return nodes, stations
}
