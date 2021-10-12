package errs

import (
	"errors"
)

var (
	TxBegin  = errors.New("error starting transaction")
	TxCommit = errors.New("error committing transaction")
	Saving   = errors.New("error on save")
)
