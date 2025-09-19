package checker

import "urlcheck/types"

type Checker interface {
	Check(url string) *types.Result
}
