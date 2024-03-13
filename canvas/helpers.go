package canvas

import (
	"strings"
	"time"

	"github.com/ninja-software/terror/v2"
)

func ReplaceSpaceInStr(str string, replacer string) string {
	str = strings.TrimSpace(str)
	str = strings.Replace(str, " ", replacer, -1)
	return str
}

func UTCToPerthTime(utc string) (string, error) {
	if utc == "" {
		return "", nil
	}

	t, err := time.Parse(time.RFC3339, utc)
	if err != nil {
		return "", terror.Error(err, "error parsing time")
	}

	perthLoc, err := time.LoadLocation("Australia/Perth")
	if err != nil {
		return "", terror.Error(err, "error loading location")
	}

	perthTime := t.In(perthLoc)
	return perthTime.String(), nil
}
