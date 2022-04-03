package domain

import "time"

type HomepageLinks struct {
	Cards   int     `json:"cards,omitempty" bson:"cards,omitempty"`
	Details details `json:"details,omitempty" bson:"details,omitempty"`
}

type details []struct {
	System    string    `json:"system,omitempty" bson:"system,omitempty"`
	Board     string    `json:"board,omitempty" bson:"board,omitempty"`
	Subject   string    `json:"subject,omitempty" bson:"subject,omitempty"`
	StartYear time.Time `json:"startyear,omitempty" bson:"startyear"`
	EndYear   time.Time `json:"endyear,omitempty" bson:"endyear"`
}
