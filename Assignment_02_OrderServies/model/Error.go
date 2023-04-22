package model

type Error struct {

	// Code specific to Incentiveweb API (nothing to do with the HTTP error-code)
	Code string `json:"code"`

	// Message targeted to developers. Not intended for users
	Message string `json:"message"`

	// milliseconds since epoch
	Timestamp int64 `json:"timestamp"`
}
