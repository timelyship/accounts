package application

var (
	StringConst StringConstants = StringConstants{
		LogLevelInfo:    "INFO",
		LogLevelDebug:   "DEBUG",
		Email:           "Email",
		Phone:           "Phone",
		VerifyEmail:     "VerifyEmail",
		PasswordPattern: "^[a-zA-Z0-9]{8,}$",
		//nolint:lll    // https://stackoverflow.com/questions/19605150/regex-for-password-must-contain-at-least-eight-characters-at-least-one-number-a
		EmailPattern: "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?" +
			"(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
	}

	IntConst IntConstants = IntConstants{
		FirstNameMaxLen: 32, //nolint:gomnd
		LastNameMaxLen:  32, //nolint:gomnd
		EmailNameMaxLen: 32, //nolint:gomnd
	}
)

type StringConstants struct {
	LogLevelInfo    string
	LogLevelDebug   string
	EmailPattern    string
	PasswordPattern string
	Email           string
	Phone           string
	VerifyEmail     string
}

type IntConstants struct {
	FirstNameMaxLen int
	LastNameMaxLen  int
	EmailNameMaxLen int
}
