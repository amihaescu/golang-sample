package model

type Device struct {
	Name      string `json:"name"`
	Long      int    `json:"long"`
	Lat       int    `json:"lat"`
	IpAddress string `json:"ipAddress"`
}
