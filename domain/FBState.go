package domain

type FBState struct {
	BaseEntity `bson:",inline"`
	State      string `bson:"state"`
}
