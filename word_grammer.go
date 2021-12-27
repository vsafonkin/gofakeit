package gofakeit

import (
	rand "math/rand"
	"unicode"
)

// simple, compound, complex, and compound-complex

// SentenceSimple will generate a random simple sentence
func SentenceSimple() string { return sentenceSimple(globalFaker.Rand) }

// SentenceSimple will generate a random simple sentence
func (f *Faker) SentenceSimple() string { return sentenceSimple(f.Rand) }

func sentenceSimple(r *rand.Rand) string {
	// simple sentence consists of a noun phrase and a verb phrase
	str := phraseNoun(r) + " " + phraseVerb(r) + "."

	// capitalize the first letter
	strR := []rune(str)
	strR[0] = unicode.ToUpper(strR[0])
	return string(strR)
}
