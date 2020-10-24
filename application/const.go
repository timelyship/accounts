package application

var (
	STRING_CONST StringConstants = StringConstants{
		LOG_LEVEL_INFO:  "INFO",
		LOG_LEVEL_DEBUG: "DEBUG",
		EMAIL_PATTERN:   "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
	}

	INT_CONST IntConstants = IntConstants{
		FIRST_NAME_MAX_LEN: 32,
		LAST_NAME_MAX_LEN:  32,
		EMAIL_NAME_MAX_LEN: 32,
	}
)

type StringConstants struct {
	LOG_LEVEL_INFO  string
	LOG_LEVEL_DEBUG string
	EMAIL_PATTERN   string
}

type IntConstants struct {
	FIRST_NAME_MAX_LEN int
	LAST_NAME_MAX_LEN  int
	EMAIL_NAME_MAX_LEN int
}
