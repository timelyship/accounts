package domain

type GoogleState struct {
	BaseEntity `bson:",inline"`
	State      string `bson:"state"`
}
