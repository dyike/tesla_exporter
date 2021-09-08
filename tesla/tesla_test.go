package tesla

import (
	"crypto/tls"
	"log"
	"net/http"
	"testing"
)

func TestVehicleData(t *testing.T) {
	var (
		token *Token
		err   error
	)
	// Don't verify TLS certs...
	tls := &tls.Config{InsecureSkipVerify: true}
	// Get TLS transport
	tr := &http.Transport{TLSClientConfig: tls}
	// Make an HTTPS client
	client := &http.Client{Transport: tr}
	email := ""
	password := ""
	token, err = GetAndCacheToken(client, &email, &password)
	if err != nil {
		return
	}
	vehicles, err := ListVehicles(client, token)
	for _, v := range *vehicles {
		log.Println("VehicleID", v.VehicleID, "ID", v.ID)
		vv, err := GetVehicleData(client, token, v.ID)
		if err != nil {
			log.Println("err", err)
			return
		}
		log.Println("Version", vv.Vs.CarVersion)
	}
}
