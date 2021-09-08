package tesla

/// POST Get Access Token
// Auth is an authorization structure for the Tesla API.
var AuthURL = "/oauth/token"

type RefreshAuthToken struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
	ClientID     string `json:"client_id"`
	Scope        string `json:"scope"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	CreatedAt    int    `json:"created_at"`
}

var ListVehiclesURL = "/api/1/vehicles"

// Vehicle is a structure that describes a single Tesla vehicle.
type Vehicle struct {
	ID                     int         `json:"id"`
	VehicleID              int         `json:"vehicle_id"`
	Vin                    string      `json:"vin"`
	DisplayName            string      `json:"display_name"`
	OptionCodes            string      `json:"option_codes"`
	Color                  interface{} `json:"color"`
	Tokens                 []string    `json:"tokens"`
	State                  string      `json:"state"`
	InService              interface{} `json:"in_service"`
	IDS                    string      `json:"id_s"`
	CalendarEnabled        bool        `json:"calendar_enabled"`
	BackseatToken          interface{} `json:"backseat_token"`
	BackseatTokenUpdatedAt interface{} `json:"backseat_token_updated_at"`
}

// Vehicles encapsulates a collection of Tesla Vehicles.
type Vehicles []struct {
	*Vehicle
}

// VehiclesResponse is the response to a vehicles API query.
type VehiclesResponse struct {
	Response Vehicles `json:"response"`
	Count    int      `json:"count"`
}

// Mobile Enabled
var MobileEnabledURL = "api/1/vehicles/%d/mobile_enabled"

type MobileEnabledResponse struct {
	Response bool `json:"response"`
}

// Charge State
var ChargStateURL = "/api/1/vehicles/%d/data_request/charge_state"

// ChargeState is the actual charge_state data
type ChargeState struct {
	BatteryHeaterOn              bool        `json:"battery_heater_on"`
	BatteryLevel                 int         `json:"battery_level"`
	BatteryRange                 float64     `json:"battery_range"`
	ChargeCurrentRequest         int         `json:"charge_current_request"`
	ChargeCurrentRequestMax      int         `json:"charge_current_request_max"`
	ChargeEnableRequest          bool        `json:"charge_enable_request"`
	ChargeLimitSoc               int         `json:"charge_limit_soc"`
	ChargeLimitSocMax            int         `json:"charge_limit_soc_max"`
	ChargeLimitSocMin            int         `json:"charge_limit_soc_min"`
	ChargeLimitSocStd            int         `json:"charge_limit_soc_std"`
	ChargeMilesAddedIdeal        float64     `json:"charge_miles_added_ideal"`
	ChargeMilesAddedRated        float64     `json:"charge_miles_added_rated"`
	ChargePortColdWeatherMode    bool        `json:"charge_port_cold_weather_mode"`
	ChargePortDoorOpen           bool        `json:"charge_port_door_open"`
	ChargePortLatch              string      `json:"charge_port_latch"` // "Engaged", "Disengaged"
	ChargeRate                   float64     `json:"charge_rate"`
	ChargeToMaxRange             bool        `json:"charge_to_max_range"`
	ChargerActualCurrent         int         `json:"charge_actual_current"`
	ChargerPhases                int         `json:"charge_phases"` // 1?
	ChargerPilotCurrent          int         `json:"charger_pilot_current"`
	ChargerPower                 int         `json:"charger_power"`
	ChargerVoltage               int         `json:"charger_voltage"`
	ChargingState                string      `json:"charging_state"` // "Stopped", "Starting", "Charging", "Disconnected"
	ConnChargeCable              string      `json:"conn_charge_cable"`
	EstBatteryRange              float64     `json:"est_battery_range"`
	FastChargerBrand             string      `json:"fast_charger_brand"`
	FastChargerPresent           bool        `json:"fast_charger_present"`
	FastChargerType              string      `json:"fast_charger_type"`
	IdealBatteryRange            float64     `json:"ideal_battery_range"`
	ManagedChargingActive        bool        `json:"managed_charging_active"`
	ManagedChargingStartTime     interface{} `json:"managed_charging_start_time"`
	ManagedChargingUserCancelled bool        `json:"managed_charging_user_cancelled"`
	MaxRangeChargeCounter        int         `json:"max_range_charge_counter"`
	NotEnoughPowerToHeat         bool        `json:"not_enough_power_to_heat"`
	ScheduledChargingPending     bool        `json:"scheduled_charging_pending"`
	ScheduledChargingStartTime   int         `json:"scheduled_charging_start_time"` // seconds
	TimeToFullCharge             float64     `json:"time_to_full_charge"`           // in hours
	TimeStamp                    int         `json:"timestamp"`                     // ms
	TripCharging                 bool        `json:"trip_charging"`
	UsableBatteryLevel           int         `json:"usable_battery_level"`
	UserChargeEnableRequest      bool        `json:"user_charge_enable_request"`
}

type ChargeStateResponse struct {
	Response ChargeState
}

// ClimateState returns the state of the climate control
type ClimateState struct {
	BatteryHeater              bool    `json:"battery_heater"`
	BatteryHeaterNoPower       bool    `json:"battery_heater_no_power"`
	DriverTempSetting          float64 `json:"driver_temp_setting"`
	FanStatus                  int     `json:"fan_status"`
	InsideTemp                 float64 `json:"inside_temp"`
	IsAutoConditioningOn       bool    `json:"is_auto_conditioning_on"`
	IsClimateOn                bool    `json:"is_climate_on"`
	IsFrontDefrosterOn         bool    `json:"is_front_defroster_on"`
	IsPreconditioning          bool    `json:"is_preconditioning"`
	IsRearDefrosterOn          bool    `json:"is_rear_defroster_on"`
	LeftTempDirection          int     `json:"left_temp_direction"`
	MaxAvailTemp               float64 `json:"max_avail_temp"`
	MinAvailTemp               float64 `json:"min_avail_temp"`
	OutsideTemp                float64 `json:"outside_temp"`
	PassengerTempSetting       float64 `json:"passenger_temp_setting"`
	RemoteHeaterControlEnabled bool    `json:"remote_heater_control_enabled"`
	RightTempDirection         int     `json:"right_temp_direction"`
	SeatHeaterLeft             int     `json:"seat_heater_left"`
	SeatHeaterRearCenter       int     `json:"seat_heater_rear_center"`
	SeatHeaterRearLeft         int     `json:"seat_heater_rear_left"`
	SeatHeaterRearLeftBack     int     `json:"seat_heater_rear_left_back"`
	SeatHeaterRearRight        int     `json:"seat_heater_rear_right"`
	SeatHeaterRearRightBack    int     `json:"seat_heater_rear_right_back"`
	SeatHeaterRight            int     `json:"seat_heater_right"`
	SideMirrorHeaters          bool    `json:"side_mirror_heaters"`
	SmartPreconditioning       bool    `json:"smart_preconditioning"`
	SteeringWheelHeater        bool    `json:"steering_wheel_heater"`
	TimeStamp                  int     `json:"timestamp"` // ms
	WiperBladeHeater           bool    `json:"wiper_blade_heater"`
}

var ClimateStateURL = "/api/1/vehicles/%d/data_request/climate_state"

type ClimateStateResponse struct {
	Response ClimateState
}

// DriveState is the result of the drive_state call, and includes information
// about vehicle position and speed
type DriveState struct {
	GpsAsOf                 int         `json:"gps_as_of"`
	Heading                 int         `json:"heading"`
	Latitude                float64     `json:"latitude"`
	Longitude               float64     `json:"longitude"`
	NativeLatitude          float64     `json:"native_latitude"`
	NativeLocationSupported int         `json:"native_location_supported"`
	NativeLongitude         float64     `json:"native_longitude"`
	NativeType              string      `json:"native_type"`
	Power                   int         `json:"power"`
	ShiftState              interface{} `json:"shift_state"`
	Speed                   interface{} `json:"speed"`
	TimeStamp               int         `json:"timestamp"` // ms
}

var DriveStateURL = "/api/1/vehicles/%d/data_request/drive_state"

// DriveStateResponse encapsulates a DriveState object.
type DriveStateResponse struct {
	Response DriveState
}

// GuiSettings return a number of settings regarding the GUI on the CID
type GuiSettings struct {
	Gui24HourTime       bool   `json:"gui_24_hour_time"`
	GuiChargeRateUnits  string `json:"gui_charge_rate_units"`
	GuiDistanceUnits    string `json:"gui_distance_units"`
	GuiRangeDisplay     string `json:"gui_range_display"`
	GuiTemperatureUnits string `json:"gui_temperature_units"`
	TimeStamp           int    `json:"timestamp"` // ms
}

var GuiStateURL = "/api/1/vehicles/%d/data_request/gui_settings"

// GuiSettingsResponse encapsulates a GuiSettings object
type GuiSettingsResponse struct {
	Response GuiSettings
}

// A VehicleStateMediaState returns the state of media control
type VehicleStateMediaState struct {
	RemoteControlEnabled bool `json:"remote_control_enabled"`
}

// A VehicleStateSoftwareUpdate returns information on pending software updates
type VehicleStateSoftwareUpdate struct {
	ExpectedDurationSec int    `json:"expected_duration_sec"`
	Status              string `json:"status"`
}

// A VehicleStateSpeedLimitMode returns the speed limiting parameters
type VehicleStateSpeedLimitMode struct {
	Active          bool    `json:"active"`
	CurrentLimitMph float64 `json:"current_limit_mph"`
	MaxLimitMph     int     `json:"max_limit_mph"`
	MinLimitMph     int     `json:"min_limit_mph"`
	PinCodeSet      bool    `json:"pin_code_set"`
}

// VehicleState is the return value from a vehicle_state call
type VehicleState struct {
	APIVersion              int                        `json:"api_version"`
	AutoparkStateV2         string                     `json:"autopark_state_v2"`
	AutoparkStyle           string                     `json:"autopark_style"`
	CalendarSupported       bool                       `json:"calendar_supported"`
	CarVersion              string                     `json:"car_version"`
	CenterDisplayState      int                        `json:"center_display_state"`
	Df                      int                        `json:"df"`
	Dr                      int                        `json:"dr"`
	Ft                      int                        `json:"ft"`
	HomelinkNearby          bool                       `json:"homelink_nearby"`
	IsUserPresent           bool                       `json:"is_user_present"`
	LastAutoparkError       string                     `json:"last_autopark_error"`
	Locked                  bool                       `json:"locked"`
	MediaState              VehicleStateMediaState     `json:"media_state"`
	NotificationsSupported  bool                       `json:"notifications_supported"`
	Odometer                float64                    `json:"odometer"`
	ParsedCalendarSupported bool                       `json:"parsed_calendar_supported"`
	Pf                      int                        `json:"pf"`
	Pr                      int                        `json:"pr"`
	RemoteStart             bool                       `json:"remote_start"`
	RemoteStartSupported    bool                       `json:"remote_start_started"`
	Rt                      int                        `json:"rt"`
	SoftwareUpdate          VehicleStateSoftwareUpdate `json:"software_update"`
	SpeedLimitMode          VehicleStateSpeedLimitMode `json:"speed_limit_mode"`
	SunRoofPercentOpen      int                        `json:"sun_roof_percent_open"`
	SunRoofState            string                     `json:"sun_roof_state"`
	TimeStamp               int                        `json:"timestamp"` // ms
	ValetMode               bool                       `json:"valet_mode"`
	ValetPinNeeded          bool                       `json:"valet_pin_needed"`
	VehicleName             string                     `json:"vehicle_name"`
}

var VehicleStateURL = "/api/1/vehicles/%d/data_request/vehicle_state"

// VehicleStateResponse encapsulates a VehicleState object
type VehicleStateResponse struct {
	Response VehicleState
}

// VehicleConfig is the return data from a vehicle_config call
type VehicleConfig struct {
	CanAcceptNavigationRequests bool   `json:"can_accept_navigation_requests"`
	CanActuateTrunks            bool   `json:"can_actuate_trunks"`
	CarSpecialType              string `json:"car_special_type"` // "base"
	CarType                     string `json:"car_type"`         // "models"
	ChargePortType              string `json:"charge_port_type"`
	EuVehicle                   bool   `json:"eu_vehicle"`
	ExteriorColor               string `json:"exterior_color"`
	HasAirSuspension            bool   `json:"has_air_suspension"`
	HasLudicrousMode            bool   `json:"has_ludicrous_mode"`
	MotorizedChargePort         bool   `json:"motorized_charge_port"`
	PerfConfig                  string `json:"perf_config"`
	Plg                         bool   `json:"plg"`
	RearSeatHeaters             int    `json:"rear_seat_heaters"`
	RearSeatType                int    `json:"rear_seat_type"`
	Rhd                         bool   `json:"rhd"`
	RoofColor                   string `json:"roof_color"` // "Colored"
	SeatType                    int    `json:"seat_type"`
	SpoilerType                 string `json:"spoiler_type"`
	SunRoofInstalled            int    `json:"sun_roof_installed"`
	ThirdRowSeats               string `json:"third_row_seats"`
	TimeStamp                   int    `json:"timestamp"` // ms
	TrimBadging                 string `json:"trim_badging"`
	WheelType                   string `json:"wheel_type"`
}

var VehicleConfigURL = "/api/1/vehicles/%d/data_request/vehicle_config"

// VehicleConfigResponse encapsulates a VehicleConfig
type VehicleConfigResponse struct {
	Response VehicleConfig
}

// VehicleData is the actual data structure for a vehicle_data call
type VehicleData struct {
	Vehicle
	UserID int           `json:"user_id"`
	Ds     DriveState    `json:"drive_state"`
	Cls    ClimateState  `json:"climate_state"`
	Chs    ChargeState   `json:"charge_state"`
	Gs     GuiSettings   `json:"gui_settings"`
	Vs     VehicleState  `json:"vehicle_state"`
	Vc     VehicleConfig `json:"vehicle_config"`
}

var VehicleDataURL = "/api/1/vehicles/%d/vehicle_data"

// VehicleDataResponse is the return from a vehicle_data call
type VehicleDataResponse struct {
	Response VehicleData
}
