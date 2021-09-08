package tesla

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var (
	// BaseURL is tesla api url.
	BaseURL = "https://owner-api.teslamotors.com"

	UserAgent     = "ityike.lab"
	teslaClientID = "81527cff06843c8634fdc09e8ac0abefb46ac849f38fe1e431c2ef2106796384"
	// TokenCachePath is the location (UNIX specific?) to cache API credentials.
	TokenCachePath = os.Getenv("HOME") + "/.tesla.cache"

	// TokenCachePathNewSuffix is the suffix to add to a new cache file when updating.
	TokenCachePathNewSuffix = ".new"
)

// GetToken
func GetToken(client *http.Client, username *string, password *string) (*Token, error) {
	accessToken := Login(*username, *password)
	var res Token
	err := json.Unmarshal([]byte(accessToken), &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// RefreshToken
func RefreshToken(client *http.Client, token *Token) (*Token, error) {
	var auth RefreshAuthToken
	auth.GrantType = "refresh_token"
	auth.RefreshToken = token.RefreshToken
	auth.ClientID = teslaClientID
	auth.Scope = "openid email offline_access"

	return authCommon(client, &auth)
}

// List Tesla vehicles
func ListVehicles(client *http.Client, token *Token) (*Vehicles, error) {
	vehicleJson, err := GetRequest(client, token, ListVehiclesURL)
	if err != nil {
		return nil, err
	}
	var res VehiclesResponse
	err = json.Unmarshal(vehicleJson, &res)
	if err != nil {
		return nil, err
	}

	// check response count
	if len(res.Response) != res.Count {
		return nil, fmt.Errorf("List vehicles response length %d != Count %d", len(res.Response), res.Count)
	}
	return &(res.Response), nil
}

// Whether mobile access is enabled.
func GetMobileEnabled(client *http.Client, token *Token, id int) (bool, error) {
	resJson, err := GetRequest(client, token, fmt.Sprintf(MobileEnabledURL, id))
	if err != nil {
		return false, err
	}
	var res MobileEnabledResponse
	err = json.Unmarshal(resJson, &res)
	if err != nil {
		return false, err
	}
	return res.Response, nil
}

// Charge State
func GetChargeState(client *http.Client, token *Token, id int) (*ChargeState, error) {
	resJson, err := GetRequest(client, token, fmt.Sprintf(ChargStateURL, id))
	if err != nil {
		return nil, err
	}
	var res ChargeStateResponse
	err = json.Unmarshal(resJson, &res)
	if err != nil {
		return nil, err
	}
	return &(res.Response), nil
}

// GetClimateState
// returns information on the current internal temperature and climate control system.
func GetClimateState(client *http.Client, token *Token, id int) (*ClimateState, error) {
	resJson, err := GetRequest(client, token, fmt.Sprintf(ClimateStateURL, id))
	if err != nil {
		return nil, err
	}
	var res ClimateStateResponse
	err = json.Unmarshal(resJson, &res)
	if err != nil {
		return nil, err
	}
	return &(res.Response), nil
}

// GetDriveState returns the driving and position state of the vehicle
func GetDriveState(client *http.Client, token *Token, id int) (*DriveState, error) {
	resJson, err := GetRequest(client, token, fmt.Sprintf(DriveStateURL, id))
	if err != nil {
		return nil, err
	}
	var res DriveStateResponse
	err = json.Unmarshal(resJson, &res)
	if err != nil {
		return nil, err
	}
	return &(res.Response), nil
}

// GetGuiSettings
func GetGuiSettings(client *http.Client, token *Token, id int) (*GuiSettings, error) {
	resJson, err := GetRequest(client, token, fmt.Sprintf(GuiStateURL, id))
	if err != nil {
		return nil, err
	}
	var res GuiSettingsResponse
	err = json.Unmarshal(resJson, &res)
	if err != nil {
		return nil, err
	}
	return &(res.Response), nil
}

// GetVehicleState returns the vehicle's physical state, such as which doors are open
func GetVehicleState(client *http.Client, token *Token, id int) (*VehicleState, error) {
	resJson, err := GetRequest(client, token, fmt.Sprintf(VehicleStateURL, id))
	if err != nil {
		return nil, err
	}
	var res VehicleStateResponse
	err = json.Unmarshal(resJson, &res)
	if err != nil {
		return nil, err
	}
	return &(res.Response), nil
}

// GetVehicleConfig
func GetVehicleConfig(client *http.Client, token *Token, id int) (*VehicleConfig, error) {
	resJson, err := GetRequest(client, token, fmt.Sprintf(VehicleConfigURL, id))
	if err != nil {
		return nil, err
	}
	var res VehicleConfigResponse
	err = json.Unmarshal(resJson, &res)
	if err != nil {
		return nil, err
	}
	return &(res.Response), nil
}

// GetVehicleData
func GetVehicleData(client *http.Client, token *Token, id int) (*VehicleData, error) {
	resJson, err := GetRequest(client, token, fmt.Sprintf(VehicleDataURL, id))
	if err != nil {
		return nil, err
	}
	var res VehicleDataResponse
	err = json.Unmarshal(resJson, &res)
	if err != nil {
		return nil, err
	}
	return &(res.Response), nil
}

func authCommon(client *http.Client, auth *RefreshAuthToken) (*Token, error) {
	authJson, err := json.Marshal(auth)
	if err != nil {
		return nil, err
	}
	body, err := PostRequest(client, nil, AuthURL, authJson)
	if err != nil {
		return nil, err
	}
	var res Token
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// SaveCachedToken saves a Token structure (JSON representation)
func saveCachedToken(t *Token) error {
	tokenJson, err := json.Marshal(t)
	if err != nil {
		return nil
	}

	// write to local file
	err = ioutil.WriteFile(TokenCachePath+TokenCachePathNewSuffix, tokenJson, 0600)
	if err != nil {
		return err
	}
	// Move into place
	err = os.Rename(TokenCachePath+TokenCachePathNewSuffix, TokenCachePath)
	if err != nil {
		return err
	}
	return nil
}

// GetAndCacheToken gets a new token and saves it in the local filesystem.
func GetAndCacheToken(client *http.Client, username *string, password *string) (*Token, error) {
	t, err := GetToken(client, username, password)
	if err != nil {
		return t, err
	}

	err = saveCachedToken(t)
	if err != nil {
		return t, err
	}
	return t, nil
}

// RefreshAndCacheToken refreshs and saves the returned token in the local filesystem
func RefreshAndCacheToken(client *http.Client, token *Token) (*Token, error) {
	t, err := RefreshToken(client, token)
	if err != nil {
		return t, err
	}
	err = saveCachedToken(t)
	if err != nil {
		return t, err
	}

	return t, nil
}

// LoadCachedToken returns the token (if any) from the cache file.
func LoadCachedToken() (*Token, error) {
	var t Token

	body, err := ioutil.ReadFile(TokenCachePath)
	if err != nil {
		return nil, err
	}

	// Parse response, get token structure
	err = json.Unmarshal(body, &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// DeleteCachedToken removes the cached token file.
func DeleteCachedToken() error {
	err := os.Remove(TokenCachePath)
	return err
}

// CheckToken returns true if a token is valid.
func CheckToken(t *Token) bool {
	start, end := tokenTimes(t)
	now := time.Now()
	return (start.Before(now) && now.Before(end))
}

// TokenTimes returns the start and end times for a token.
func tokenTimes(t *Token) (start, end time.Time) {
	start = time.Unix(int64(t.CreatedAt), 0)
	end = time.Unix(int64(t.CreatedAt)+int64(t.ExpiresIn), 0)
	return
}
