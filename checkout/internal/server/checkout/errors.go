package checkout

import "errors"

var (
	ErrEmptyCount = errors.New("empty Count")
	ErrEmptySKU   = errors.New("empty Sku")
	ErrEmptyUser  = errors.New("empty User")
)
