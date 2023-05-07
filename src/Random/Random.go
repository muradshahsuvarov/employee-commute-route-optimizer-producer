package Random

import (
	"math/rand"
	"time"
)

func GenerateRandomID(idLength int, characters string) string {

	rand.Seed(time.Now().UnixNano())

	id := make([]byte, idLength)

	for i := 0; i < idLength; i++ {
		id[i] = characters[rand.Intn(len(characters))]
	}

	return string(id)
}
