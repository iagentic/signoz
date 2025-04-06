package types

import "errors"

var (
	ErrColumnNotFound = errors.New("column not found")
	ErrBetweenValues  = errors.New("(not) between operator requires two values")
	ErrInValues       = errors.New("(not) in operator requires a list of values")
)

// FilterOperator is the operator for the filter.
type FilterOperator int

const (
	FilterOperatorUnknown FilterOperator = iota
	FilterOperatorEqual
	FilterOperatorNotEqual
	FilterOperatorGreaterThan
	FilterOperatorGreaterThanOrEq
	FilterOperatorLessThan
	FilterOperatorLessThanOrEq

	FilterOperatorLike
	FilterOperatorNotLike
	FilterOperatorILike
	FilterOperatorNotILike

	FilterOperatorBetween
	FilterOperatorNotBetween

	FilterOperatorIn
	FilterOperatorNotIn

	FilterOperatorExists
	FilterOperatorNotExists

	FilterOperatorRegexp
	FilterOperatorNotRegexp

	FilterOperatorContains
	FilterOperatorNotContains
)
