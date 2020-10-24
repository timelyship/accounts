package request

import "timelyship.com/accounts/utility"

type Request interface {
	Validate() *utility.RestError
}
