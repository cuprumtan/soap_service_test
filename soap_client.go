package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	Host string `yaml:"soap_host"`
	Port string `yaml:"soap_port"`
}

type GoogleTranslate struct {
	XMLName xml.Name
	Body    struct {
		XMLName           xml.Name
		GoogleTranslateResponse struct {
			XMLName xml.Name
			Return  []string `xml:"GoogleTranslateResult"`
		} `xml:"GoogleTranslateResponse"`
	}
}

func (c *Config ) getConf() *Config  {

	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func main() {
	
	var c Config
	c.getConf()
	
	fmt.Println("\nGoogle Translate console service\n")
	
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Source text: ")
	text, _ := reader.ReadString('\n')
	text = text[:len(text)-2]
	fmt.Print("Source language: ")
	sourceLanguage, _ := reader.ReadString('\n')
	sourceLanguage = sourceLanguage[:len(sourceLanguage)-2]
	fmt.Print("Destination language: ")
	destinationLanguage, _ := reader.ReadString('\n')
	destinationLanguage = destinationLanguage[:len(destinationLanguage)-2]
	fmt.Print("\n")
		
	// wsdl service url
	url := fmt.Sprintf("%s%s%s%s%s",
		"http://",
		c.Host,
		":",
		c.Port,
		"/?wsdl",
	)
	
	// payload
	payload := []byte(strings.TrimSpace(`
		<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:tran="Translator">
		   <soapenv:Body>
		      <tran:GoogleTranslate>
		         <!--Optional:-->
		         <tran:source_text>` + text + `</tran:source_text>
		         <!--Optional:-->
		         <tran:source_language>` + sourceLanguage + `</tran:source_language>
		         <!--Optional:-->
		         <tran:destination_language>` + destinationLanguage + `</tran:destination_language>
		      </tran:GoogleTranslate>
		   </soapenv:Body>
		</soapenv:Envelope>`,
	))

	httpMethod := "POST"

	// soap action
	soapAction := "GoogleTranslate"

	log.Println("-> Preparing the request")

	// prepare the request
	req, err := http.NewRequest(httpMethod, url, bytes.NewReader(payload))
	if err != nil {
		log.Fatal("Error on creating request object. ", err.Error())
		return
	}

	// set the content type header, as well as the oter required headers
	req.Header.Set("Content-type", "text/xml")
	req.Header.Set("SOAPAction", soapAction)

	// prepare the client request
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	log.Println("-> Dispatching the request")

	// dispatch the request
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Error on dispatching request. ", err.Error())
		return
	}

	log.Println("-> Retrieving and parsing the response")

	// read and parse the response body
	result := new(GoogleTranslate)
	err = xml.NewDecoder(res.Body).Decode(result)
	if err != nil {
		log.Fatal("Error on unmarshaling xml. ", err.Error())
		return
	}

	log.Println("-> Everything is good, printing users data")

	// print the users data
	users := result.Body.GoogleTranslateResponse.Return
	fmt.Println("\nTranslated text:")
	fmt.Println(strings.Join(users, ", "))
}