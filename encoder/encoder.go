// Package encoder is a package that is primarily responsible for encoding a string url
// which is passed to it from SaveUrl. In this case, encoding is the process of
// taking a normal, human readable url such as http://google.com and converting
// it into a smaller, more easily transported, typed, and/or remembered url.
package encoder

const (
	alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	length   = int64(len(alphabet))
)

// Encode number to base62.
func Encode(n int64) string {
	if n == 0 {
		return string(alphabet[0])
	}

	s := ""
	for ; 0 < n; n = n / length {
		s = string(alphabet[n%length]) + s
	}

	return s
}
