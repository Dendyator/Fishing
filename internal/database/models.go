package database

import "time"

type FishingData struct {
	ID             int                `db:"id"`
	Date           time.Time          `db:"date"`
	Location       string             `db:"location"`
	Coordinates    string             `db:"coordinates"`
	Comment        string             `db:"comment"`
	FishingMethods []string           `db:"fishing_methods"`
	CaughtFish     []CaughtFish       `db:"-"`
	TrophyFish     *FishingDataTrophy `db:"-"`
}

type CaughtFish struct {
	Species string  `db:"species"`
	Weight  float64 `db:"weight"`
}

type FishingDataTrophy struct {
	Species      string  `db:"species"`
	Weight       float64 `db:"weight"`
	CatchMethod  string  `db:"catch_method"`
	Bait         string  `db:"bait"`
	BaitDetails  string  `db:"bait_details"`
	PhotoURL     string  `db:"photo_url"`
	LurePhotoURL string  `db:"lure_photo_url"`
}
