package application

var (
	StringConst StringConstants = StringConstants{
		LogLevelInfo:  "INFO",
		LogLevelDebug: "DEBUG",
		Email:         "Email",
		Phone:         "PhoneNumber",
		// this string is important, do not change it, if you change it you have to change the email service,
		//utility.js if(payload.context === 'VERIFY_EMAIL')
		VerifyEmail: "VERIFY_EMAIL",
		//nolint:lll    // https://stackoverflow.com/questions/19605150/regex-for-password-must-contain-at-least-eight-characters-at-least-one-number-a
		PasswordPattern: "^[a-zA-Z0-9]{8,}$",
		EmailPattern: "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?" +
			"(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
		DefaultPicture: "16100991235ff80b9e",
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
	DefaultPicture  string
}

type IntConstants struct {
	FirstNameMaxLen int
	LastNameMaxLen  int
	EmailNameMaxLen int
}
