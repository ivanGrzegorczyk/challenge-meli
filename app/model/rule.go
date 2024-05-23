package model

type Rule struct {
	Ip          string `json:"ip" bson:"ip"`
	Path        string `json:"path" bson:"path"`
	MaxRequests int    `json:"max_requests" bson:"max_requests"`
	Time        int    `json:"time" bson:"time"`
}
