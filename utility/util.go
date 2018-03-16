package utility

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

//Semester - calculates current semester
func Semester() (round int) {
	NowMonth := time.Now().UTC().Month()

	if NowMonth == 1 || NowMonth == 2 || NowMonth == 3 || NowMonth == 4 || NowMonth == 5 || NowMonth == 6 {
		round = 1
	} else {
		round = 3
	}
	return round
}

//GetServiceXML - get data from an URL xml service
func GetServiceXML(T interface{}, url string) error {

	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("GET error: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Status error: %v", response.StatusCode)
	}
	data, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return fmt.Errorf("Read body: %v", err)
	}
	xml.Unmarshal(data, &T)
	return err
}

//GetInitEnd - returns the start and end of a semester in time format
func GetInitEnd() (fromDate time.Time, toDate time.Time) {
	Semester := Semester()
	var Inicial time.Month
	var Final time.Month
	if Semester == 1 {
		Inicial = time.January
		Final = time.June
	} else {
		Inicial = time.July
		Final = time.December
	}
	fromDate = time.Date(time.Now().UTC().Year(), Inicial, 1, 0, 0, 0, 0, time.UTC)
	toDate = time.Date(time.Now().UTC().Year(), Final, 30, 0, 0, 0, 0, time.UTC)
	return fromDate, toDate
}

//SendJsonToRuler - send Ruler
func SendJsonToRuler(url string, trequest string, datajson interface{}) (string, error) {
	b := new(bytes.Buffer)
	if datajson != nil {
		json.NewEncoder(b).Encode(datajson)
	}
	client := &http.Client{}
	req, err := http.NewRequest(trequest, url, b)
	r, err := client.Do(req)
	if err != nil {
		beego.Error("error", err)
		return "na", err
	}
	defer r.Body.Close()
	contents, err := ioutil.ReadAll(r.Body)
	return strings.Replace(string(contents), "\"", "", -1), err
}

// //Getipaddres - returns the current ip server of api
// func Getipaddres() net.IP {
// 	conn, err := net.Dial("udp", "8.8.8.8:80")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer conn.Close()
// 	localAddr := conn.LocalAddr().(*net.UDPAddr)
// 	return localAddr.IP
// }
