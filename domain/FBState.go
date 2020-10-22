package domain

type FBState struct {
	BaseEntity `bson:",inline"`
	State      string `json:"state" bson:"state"`
}
