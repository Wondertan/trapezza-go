package session

import (
	"github.com/Wondertan/trapezza-go/utils"
)

func RandID(l int) ID {
	return ID(utils.RandString(l))
}
