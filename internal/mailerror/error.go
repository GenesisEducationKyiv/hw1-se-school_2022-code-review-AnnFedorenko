package mailerror

import "errors"

var ErrEmailSubscribed = errors.New("email already subscribed")
var ErrEmailNotValid = errors.New("email address is not valid")
