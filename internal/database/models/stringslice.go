package models

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/watchedsky-social/core/internal/utils"
)

type StringSlice []string

func (s *StringSlice) Scan(src any) error {
	var slice []string

	switch asserted := src.(type) {
	case []byte:
		str := strings.Map(func(r rune) rune {
			if r == '{' || r == '}' {
				return -1
			}

			return r
		}, string(asserted))
		slice = strings.Split(str, ",")
	case []string:
		slice = asserted
	case string:
		asserted = strings.Map(func(r rune) rune {
			if r == '{' || r == '}' {
				return -1
			}

			return r
		}, asserted)
		slice = strings.Split(asserted, ",")
	case []any:
		slice = utils.FromAnySlice[string](asserted)
	default:
		return fmt.Errorf("type %T not supported", asserted)
	}

	*s = slice
	return nil
}

func (s StringSlice) Value() (driver.Value, error) {
	quoted := utils.Map(s, func(str string) string {
		return fmt.Sprintf("%q", str)
	})

	return fmt.Sprintf("{%s}", strings.Join(quoted, ",")), nil
}
