package keygen

import "github.com/jaevor/go-nanoid"

var canonic, _ = nanoid.Standard(21)

func Generate() string {
	return canonic()
}
