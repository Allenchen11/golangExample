package models

import "time"

//Award is a award
type Award struct {
	AwardID           int       `json:"award_id"`
	AwardName         string    `json:"award_name"`
	AwardSerialNumber string    `json:"award_serial_number"`
	AwardStatus       string    `json:"award_status"`
	AwardUpdateTime   time.Time `json:"award_update_time"`
}

var award Award
