package service

import "crypto/rand"

type StateGenerator interface {
	Generate() (string, error)
}

func NewStateGenerator(length int) StateGenerator {
	return &randomStateGenerator{length: length}
}

type randomStateGenerator struct {
	length int
}

func (g *randomStateGenerator) Generate() (string, error) {
	b := make([]byte, g.length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	for i, v := range b {
		b[i] = randTable[v%byte(len(randTable))]
	}

	return string(b), nil
}
