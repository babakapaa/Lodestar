package schema

import (
	"fmt"
	"time"
	_ "time/tzdata"
)

func validZone(name string) error {
	if name == "Local" {
		return fmt.Errorf("%q is not a time zone", name)
	}
	_, err := time.LoadLocation(name)
	return err
}
