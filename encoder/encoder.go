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
