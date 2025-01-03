package kxml

import (
	"encoding/xml"
	"io"
	"os"
	"regexp"

	l "github.com/micst/kxgctl/kxg/logging"
)

const (
	GroupAddressDocumentHeader = `<?xml version="1.0" encoding="utf-8" standalone="yes"?>` + "\n"
	XmlNamespace               = `http://knx.org/xml/ga-export/01`
	XmlNameDocument            = `GroupAddress-Export`
	XmlNameGroup               = `GroupRange`
	XmlNameAddress             = `GroupAddress`
)

type GroupAddress struct {
	XMLName     xml.Name `xml:"GroupAddress"`
	Name        string   `xml:"Name,attr"`
	Address     string   `xml:"Address,attr"`
	Description string   `xml:"Description,attr,omitempty"`
	Security    string   `xml:"Security,attr,omitempty"`
	DPTs        string   `xml:"DPTs,attr,omitempty"`
}

type GroupRangeMiddle struct {
	XMLName     xml.Name       `xml:"GroupRange"`
	Name        string         `xml:"Name,attr"`
	RangeStart  string         `xml:"RangeStart,attr"`
	RangeEnd    string         `xml:"RangeEnd,attr"`
	Description string         `xml:"Description,attr,omitempty"`
	Addresses   []GroupAddress `xml:"GroupAddress"`
}

type GroupRangeMain struct {
	XMLName      xml.Name           `xml:"GroupRange"`
	Name         string             `xml:"Name,attr"`
	RangeStart   string             `xml:"RangeStart,attr"`
	RangeEnd     string             `xml:"RangeEnd,attr"`
	Description  string             `xml:"Description,attr,omitempty"`
	MiddleGroups []GroupRangeMiddle `xml:"GroupRange"`
}

type GroupAddressDocument struct {
	XMLName    xml.Name         `xml:"GroupAddress-Export"`
	Xmlns      string           `xml:"xmlns,attr"`
	MainGroups []GroupRangeMain `xml:"GroupRange"`
}

func (d *GroupAddressDocument) ReadXml(file string) {
	if _, exist_err := os.Stat(file); exist_err == nil {
		if xmlFile, err := os.Open(file); err == nil {
			defer xmlFile.Close()
			byteValue, _ := io.ReadAll(xmlFile)
			if unmarshall_err := xml.Unmarshal(byteValue, d); unmarshall_err != nil {
				l.Error("could not unmarshall file " + file)
			}
		} else {
			l.Error("error reading file ")
		}
	} else {
		l.Error("file does not exist: " + file)
	}
}

func (d *GroupAddressDocument) GetXml() string {
	r := regexp.MustCompile("></[A-Za-z-]+>\n")
	if bytes, marshal_err := xml.MarshalIndent(d, "", "  "); marshal_err == nil {
		bytes_str := string(bytes)
		bytes_str = r.ReplaceAllString(bytes_str, " />\n")
		bytes_str += "\n"
		return bytes_str
	} else {
		l.Error("could not marshal file struct")
	}
	return ""
}

func (d *GroupAddressDocument) WriteXml(file string) {
	r := regexp.MustCompile("></[A-Za-z-]+>\n")
	if bytes, marshal_err := xml.MarshalIndent(d, "", "  "); marshal_err == nil {
		bytes_str := string(bytes)
		bytes_str = r.ReplaceAllString(bytes_str, " />\n")
		if xmlFile, file_err := os.Create(file); file_err == nil {
			defer xmlFile.Close()
			if _, write_err := xmlFile.Write([]byte(GroupAddressDocumentHeader)); write_err == nil {

				xmlFile.Write([]byte(bytes_str))
			} else {
				l.Error("could write to file")
			}
		} else {
			l.Error("could not create xml file")
		}
	} else {
		l.Error("could not marshal file struct")
	}
}
