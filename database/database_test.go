package database

import (
	"os"
)

func init() {
	os.Setenv("WAVE_ENV", "testing")
}