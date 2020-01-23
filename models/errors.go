package models

import "errors"

// ErrInvalidMonth error for invalid month
var ErrInvalidMonth = errors.New("invalid month")

// ErrEndCondition error for wrong end condition
var ErrEndCondition = errors.New("wrong end condition")

// ErrSendToAPI error for failed sending stats to API
var ErrSendToAPI = errors.New("failed sending stats to API")
