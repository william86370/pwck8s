package rancher

import (
	"math/rand"
	"time"
)

func GenerateId() string {
	// Generate a new User ID similar to Rancher User ID
	// u-<random 5 char string>
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	const letterBytes = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, 5)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return "pwck8s-" + string(b)
}
