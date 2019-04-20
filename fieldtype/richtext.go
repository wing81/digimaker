//Author xc, Created on 2019-03-25 20:00
//{COPYRIGHTS}

package fieldtype

import (
	"database/sql/driver"
	"errors"
)

type RichTextField struct {
	data string
}

//when update db
func (t RichTextField) Value() (driver.Value, error) {
	return t.data, nil
}

func (t *RichTextField) Scan(src interface{}) error {
	var source string
	switch src.(type) {
	case string:
		source = src.(string)
	case []byte:
		source = string(src.([]byte))
	default:
		return errors.New("Incompatible type for GzippedText")
	}

	t.data = source
	return nil
}
