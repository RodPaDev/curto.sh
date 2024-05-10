package nanoid

import (
	"math/rand/v2"
	"strings"
)

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz~"

func Gen(size int) string {
	str := strings.Builder{}
	for i := 0; i < size; i += 1 {
		j := rand.IntN(len(alphabet))
		str.WriteByte(alphabet[j])
	}
	return str.String()
}
