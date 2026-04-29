package selectronic

type Board struct {
	BoardID         string `json:"board_id"`
	BoardShortID    string `json:"board_short_id"`
	FirmwareVersion string `json:"firmware_version"`
	HardwareVersion string `json:"hardware_version"`
	OEMBranding     string `json:"oem_branding"`
}

type Device struct {
	ID           string `json:"id"`
	ShortID      string `json:"short_id"`
	SerialNum    string `json:"serialnum"`
	Manufacturer string `json:"manufacturer"`
	DeviceType   string `json:"device_type"`
	Type         string `json:"type"`
	Firmware     string `json:"firmware"`
	Poll         int    `json:"poll"`
	PowerRating  string `json:"power_rating"`
	FaultHref    string `json:"fault_href"`
}

type Point struct {
	ItemCount int        `json:"item_count"`
	Items     PointItems `json:"items"`
	Now       int64      `json:"now"`
}

type PointItems struct {
	BatteryInWhToday  float64 `json:"battery_in_wh_today"`
	BatteryInWhTotal  float64 `json:"battery_in_wh_total"`
	BatteryOutWhToday float64 `json:"battery_out_wh_today"`
	BatteryOutWhTotal float64 `json:"battery_out_wh_total"`
	BatterySOC        float64 `json:"battery_soc"`
	BatteryW          float64 `json:"battery_w"`
	FaultCode         float64 `json:"fault_code"`
	FaultTS           int64   `json:"fault_ts"`
	GenStatus         float64 `json:"gen_status"`
	GridInWhToday     float64 `json:"grid_in_wh_today"`
	GridInWhTotal     float64 `json:"grid_in_wh_total"`
	GridOutWhToday    float64 `json:"grid_out_wh_today"`
	GridOutWhTotal    float64 `json:"grid_out_wh_total"`
	GridW             float64 `json:"grid_w"`
	LoadW             float64 `json:"load_w"`
	LoadWhToday       float64 `json:"load_wh_today"`
	LoadWhTotal       float64 `json:"load_wh_total"`
	ShuntW            float64 `json:"shunt_w"`
	SolarWhToday      float64 `json:"solar_wh_today"`
	SolarWhTotal      float64 `json:"solar_wh_total"`
	SolarInverterW    float64 `json:"solarinverter_w"`
	Timestamp         int64   `json:"timestamp"`
}

type Snapshot struct {
	Board  Board
	Device Device
	Point  Point
}
