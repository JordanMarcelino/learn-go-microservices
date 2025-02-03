package tokenutils

import (
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/rand"
)

func GenerateOTPCode() string {
	rand.Seed(uint64(time.Now().UnixNano()))
	var sb strings.Builder
	for i := 0; i < 6; i++ {
		sb.WriteString(strconv.Itoa(rand.Intn(10)))
	}
	return sb.String()
}
