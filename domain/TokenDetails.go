package domain

type TokenDetails struct {
	BaseEntity   `bson:",inline"`
	AccessToken  string `bson:"access_token"`
	RefreshToken string `bson:"refresh_token"`
	AccessUuid   string `bson:"access_token_id"`
	RefreshUuid  string `bson:"refresh_token_id"`
	AtExpires    int64  `bson:"access_token_exp"`
	RtExpires    int64  `bson:"refresh_token_exp"`
}
