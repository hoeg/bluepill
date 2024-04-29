package values

import "strings"

type Whitelist []string

func (w Whitelist) Value() string {
	return strings.Join(w, ", ")
}
