package request

import "timelyship.com/accounts/utility"

type Request interface {
	ApplyUiValidation() *utility.RestError
}
