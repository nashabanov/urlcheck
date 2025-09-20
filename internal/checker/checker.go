package checker

import "github.com/nashabanov/urlcheck/internal/types"

type Checker interface {
	Check(url string) *types.Result
}
