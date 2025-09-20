package checker

import "urlcheck/internal/types"

type Checker interface {
	Check(url string) *types.Result
}
