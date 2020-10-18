package domain

type GoogleState struct {
	BaseEntity `bson:",inline"`
	State      string `json:"state" bson:"state"`
}
