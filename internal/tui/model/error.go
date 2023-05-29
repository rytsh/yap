package model

import (
	"fmt"
	"time"
)

func TimeErr(err error, t time.Time) error {
	return fmt.Errorf("%s - %w", t.Format("15:04:05"), err)
}
