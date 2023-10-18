package loms

import "errors"

var (
	ErrEmptyCount = errors.New("empty Count")
	ErrEmptyItems = errors.New("empty Items")
	ErrEmptyOrder = errors.New("empty OrderId")
	ErrEmptySKU   = errors.New("empty Sku")
	ErrEmptyUser  = errors.New("empty User")
)
