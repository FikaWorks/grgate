//go:build integration || integrationgithub || integrationgitlab

package tests

import (
	"fmt"
	"math/rand"
	"time"
)

var characters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

const hashSize = 5

func randomString(length int) string {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))

	b := make([]rune, length)
	for i := range b {
		b[i] = characters[r.Intn(len(characters))]
	}
	return string(b)
}

func generateRandomRepositoryName(prefix string) string {
	return fmt.Sprintf("%s-%s", prefix, randomString(hashSize))
}
