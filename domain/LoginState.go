package domain

type LoginState struct {
	State string `bson:"state"`
	Key   string `bson:"key"`
}
