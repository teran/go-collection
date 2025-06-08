package helm

import (
	"strings"
)

type Resources []Resource

func (r Resources) FilterByKind(kind string) Resources {
	result := Resources{}
	for _, r := range r {
		if strings.EqualFold(r.GetString("kind"), kind) {
			result = append(result, r)
		}
	}
	return result
}
