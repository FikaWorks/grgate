//go:build integration || integrationgithub || integrationgitlab

package tests

import (
	"fmt"
	"math/rand"
)

var characters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

func randomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = characters[rand.Intn(len(characters))]
	}
	return string(b)
}

func generateRandomRepositoryName(prefix string) string {
	return fmt.Sprintf("%s-%s", prefix, randomString(hashSize))
}
