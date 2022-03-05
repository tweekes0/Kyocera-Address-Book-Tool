package element

import (
	"encoding/xml"
	"errors"
	"log"
	"regexp"

	db "github.com/kitar0s/kyocera-ab-tool/db"
)

var (
	ErrCannotCreateElement = errors.New("element could not be create")
)

const (
	xmlPattern = "></[a-zA-Z]*>"
)

/*
	contactElement models how Kyocera's see contacts within their address books
	as XML files. 
	Each of the properties within the struct correspond to an attritbute within
	the each of the XML elements.

	XMLName: name of the XML element
	Id: the position of the contact within the address book
	Type: defines the element as a contact
	DisplayName: the name for the contact
	MailAddress: the email for the contact
	SendKeisyou: an attr that was always set to 0 
	SMB*: attributes for scanning via SMB
	FTP*: attributes for scanning via FTP
	Fax*: attributes for scanning via Fax
	InetFax*:  attributes for scanning via InternetFax

*/

type contactElement struct {
	XMLName               xml.Name `xml:"Item"`
	Id                    int64    `xml:"Id,attr"`
	Type                  string   `xml:"Type,attr"`
	DisplayName           string   `xml:"DisplayName,attr"`
	SendKeisyou           string   `xml:"SendKeisyou,attr"`
	MailAddress           string   `xml:"MailAddress,attr"`
	SendCorpName          string   `xml:"SendCorpName,attr"`
	SendPostName          string   `xml:"SendPostName,attr"`
	SmbHostName           string   `xml:"SmbHostName,attr"`
	SmbPath               string   `xml:"SmbPath,attr"`
	SmbLoginName          string   `xml:"SmbLoginName,attr"`
	SmbLoginPasswd        string   `xml:"SmbLoginPasswd,attr"`
	SmbPort               string   `xml:"SmbPort,attr"`
	FtpPath               string   `xml:"FtpPath,attr"`
	FtpHostName           string   `xml:"FtpHostName,attr"`
	FtpLoginName          string   `xml:"FtpLoginName,attr"`
	FtpLoginPasswd        string   `xml:"FtpLoginPasswd,attr"`
	FtpPort               string   `xml:"FtpPort,attr"`
	FaxNumber             string   `xml:"FaxNumber,attr"`
	FaxSubaddress         string   `xml:"FaxSubaddress,attr"`
	FaxPassword           string   `xml:"FaxPassword,attr"`
	FaxCommSpeed          string   `xml:"FaxCommSpeed,attr"`
	FaxECM                string   `xml:"FaxECM,attr"`
	FaxEncryptKeyNumber   string   `xml:"FaxEncryptKeyNumber,attr"`
	FaxEncryption         string   `xml:"FaxEncryption,attr"`
	FaxEncryptBoxEnabled  string   `xml:"FaxEncryptBoxEnabled,attr"`
	FaxEncryptBoxID       string   `xml:"FaxEncryptBoxID,attr"`
	InetFAXAddr           string   `xml:"InetFAXAddr,attr"`
	InetFAXMode           string   `xml:"InetFAXMode,attr"`
	InetFAXResolution     string   `xml:"InetFAXResolution,attr"`
	InetFAXFileType       string   `xml:"InetFAXFileType,attr"`
	IFaxSendModeType      string   `xml:"IFaxSendModeType,attr"`
	InetFAXDataSize       string   `xml:"InetFAXDataSize,attr"`
	InetFAXPaperSize      string   `xml:"InetFAXPaperSize,attr"`
	InetFAXResolutionEnum string   `xml:"InetFAXResolutionEnum,attr"`
	InetFAXPaperSizeEnum  string   `xml:"InetFAXPaperSizeEnum,attr"`
}

/*
	contactElement constructor that returns a new contactElement when given a
	valid Entry and id.
*/

func newContactElement(e *db.Entry, id int64) (*contactElement, error) {
	if e == nil {
		return nil, ErrCannotCreateElement
	}

	p := new(contactElement)
	p.Id = id
	p.Type = "Contact"
	p.DisplayName = e.Name
	p.MailAddress = e.Email
	p.SendCorpName = ""
	p.SendPostName = ""
	p.SmbHostName = ""
	p.SmbPath = ""
	p.SmbLoginPasswd = ""
	p.SmbLoginName = ""
	p.SmbPort = ""
	p.FtpPath = ""
	p.FtpHostName = ""
	p.FtpLoginName = ""
	p.FtpLoginPasswd = ""
	p.FtpPort = "21"
	p.FaxNumber = ""
	p.FaxSubaddress = ""
	p.FaxPassword = ""
	p.FaxCommSpeed = "BPS_33600"
	p.FaxECM = "On"
	p.FaxEncryptKeyNumber = "0"
	p.FaxEncryption = "Off"
	p.FaxEncryptBoxEnabled = "Off"
	p.FaxEncryptBoxID = "0000"
	p.InetFAXAddr = ""
	p.InetFAXMode = "Simple"
	p.InetFAXResolution = "3"
	p.InetFAXFileType = "TIFF_MH"
	p.IFaxSendModeType = "IFAX"
	p.InetFAXDataSize = "1"
	p.InetFAXPaperSize = "1"
	p.InetFAXResolutionEnum = "Default"
	p.InetFAXPaperSizeEnum = "Default"

	return p, nil
}

/*
	oneTouchKeyElement models how Kyocera's abstract OneTouchKeys (otk) or 
	scanner shortcuts within an XML file.

	XMLName: name of the XML element
	Id: is the order in which they appear on the scanner
	AddresdId: the ID of a contactElement. This defines where the OTK should get 
	the information to scan via the addressType.
	Type: defines that the element will be a OneTouchKey
	AddressType: defines the scan method ie Email/SMB/FTP/etc
	DisplayName: The name that will appear on the OneTouchKey
*/

type oneTouchKeyElement struct {
	XMLName     xml.Name `xml:"Item"`
	Id          int64    `xml:"Id,attr"`
	AddressId   int64    `xml:"AddressId,attr"`
	Type        string   `xml:"Type,attr"`
	AddressType string   `xml:"AddressType,attr"`
	DisplayName string   `xml:"DisplayName,attr"`
}

/*
 	oneTouchKey constructor that returns a new oneTouchKey element when given 
	the proper parameters
*/

func newOneTouchKeyElement(id, addressId int64, displayName,
	addressType string) (*oneTouchKeyElement, error) {
	p := new(oneTouchKeyElement)
	p.Id = id
	p.Type = "OneTouchKey"
	p.DisplayName = displayName
	p.AddressId = addressId
	p.AddressType = addressType

	return p, nil
}

/*
	converts the an XML struct into a self closing XML string

	Logs error and exits execution if there is an error with conversion
*/

func elementToString(e interface{}) string {
	xml, err := xml.Marshal(e)
	if err != nil {
		log.Fatalf("cannot convert to xml: %q", err)
	}

	r := regexp.MustCompile(xmlPattern)
	s := r.ReplaceAllString(string(xml), "/>")

	return s
}

/*
	Converts an Entry into a string of a contactElement.

	Logs error and exits execution if there is an error with conversion
*/

func EntryToContact(e *db.Entry, id int64) string {
	contact, err := newContactElement(e, id)
	if err != nil {
		log.Fatalf("cannot create XML entry: %q", err)
	}

	return elementToString(contact)
}

/*
	Converts an Entry into a string of a oneTouchKeyElement. 

	Logs error and exits execution if there is an error with conversion
*/

func EntryToOTK(e *db.Entry, id, addressId int64, addressType string) string {
	otk, err := newOneTouchKeyElement(id, addressId, e.Name, addressType)
	if err != nil {
		log.Fatalf("cannot convert to xml: %q", err)
	}

	return elementToString(otk)
}
