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
	// email := ""
	// password := ""
	// token, err = GetAndCacheToken(client, &email, &password)
	token, err = LoadCachedToken()
	if err != nil {
		return
	}
	vehicles, err := ListVehicles(client, token)
	for _, v := range *vehicles {
		log.Println("VehicleID", v.VehicleID, "ID", v.ID)
		config, err := GetVehicleConfig(client, token, v.ID)
		if err != nil {
			log.Println("err", err)
			return
		}
		log.Println("config", *config)

		vs, err := GetVehicleState(client, token, v.ID)
		if err != nil {
			log.Println("err", err)
			return
		}
		log.Println("vehicle_state", *vs)

		setting, err := GetGuiSettings(client, token, v.ID)
		if err != nil {
			log.Println("err", err)
			return
		}
		log.Println("setting", *setting)

		drive, err := GetDriveState(client, token, v.ID)
		if err != nil {
			log.Println("err", err)
			return
		}
		log.Println("drive", *drive)

		climate, err := GetClimateState(client, token, v.ID)
		if err != nil {
			log.Println("err", err)
			return
		}
		log.Println("climate", *climate)

		charge, err := GetChargeState(client, token, v.ID)
		if err != nil {
			log.Println("err", err)
			return
		}
		log.Println("Charge state", *charge)

		ok, err := GetMobileEnabled(client, token, v.ID)
		log.Println("ok", ok)

		if err != nil {
			log.Println("err", err)
			return
		}

		vv, err := GetVehicleData(client, token, v.ID)
		if err != nil {
			log.Println("err", err)
			return
		}
		log.Println("Version", vv.Vs.CarVersion)
	}
}
