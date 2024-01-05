package helpers

import (
	"fmt"

	"github.com/lithammer/shortuuid/v4"
)

func GenerateID(prefix string) string {
	id := shortuuid.New()
	return fmt.Sprintf("%s_%s", prefix, id)
}
