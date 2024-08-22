package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/bluesky-social/indigo/api/atproto"
)

type SkeetInfo atproto.RepoCreateRecord_Output

func (s *SkeetInfo) Scan(src any) error {
	var data []byte
	switch d := src.(type) {
	case []byte:
		data = d
	case string:
		data = []byte(d)
	default:
		return fmt.Errorf("didn't expect %T", src)
	}

	return json.Unmarshal(data, s)
}

func (s SkeetInfo) Value() (driver.Value, error) {
	return json.Marshal(s)
}
