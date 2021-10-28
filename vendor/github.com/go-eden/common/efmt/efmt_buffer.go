package efmt

import "unicode/utf8"

// Use simple []byte instead of bytes.Buffer to avoid large dependency.
type buffer []byte

func (t *buffer) write(p []byte) {
	*t = append(*t, p...)
}

func (t *buffer) writeString(s string) {
	*t = append(*t, s...)
}

func (t *buffer) writeByte(c byte) {
	*t = append(*t, c)
}

func (t *buffer) writeRune(r rune) {
	if r < utf8.RuneSelf {
		*t = append(*t, byte(r))
		return
	}

	b := *t
	n := len(b)
	for n+utf8.UTFMax > cap(b) {
		b = append(b, 0)
	}
	w := utf8.EncodeRune(b[n:n+utf8.UTFMax], r)
	*t = b[:n+w]
}
