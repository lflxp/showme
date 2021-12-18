package efmt

import "reflect"

const (
	ldigits = "0123456789abcdefx"
	udigits = "0123456789ABCDEFX"
)

const (
	signed   = true
	unsigned = false
)

const (
	commaSpaceString  = ", "
	nilAngleString    = "<nil>"
	nilParenString    = "(nil)"
	nilString         = "nil"
	mapString         = "map["
	percentBangString = "%!"
	missingString     = "(MISSING)"
	badIndexString    = "(BADINDEX)"
	panicString       = "(PANIC="
	extraString       = "%!(EXTRA "
	badWidthString    = "%!(BADWIDTH)"
	badPrecString     = "%!(BADPREC)"
	noVerbString      = "%!(NOVERB)"
	invReflectString  = "<invalid reflect.Value>"
)

type Printer struct {
	pp
}

// Sprintf formats according to a format specifier and returns the result as []byte.
func (t *Printer) Sprintf(format string, a ...interface{}) []byte {
	t.reset()
	t.doPrintf(format, a)
	return t.buf
}

// Sprint formats using the default formats for its operands and returns the result as []byte.
// Spaces are added between operands when neither is a string.
func (t *Printer) Sprint(a ...interface{}) []byte {
	t.reset()
	t.doPrint(a)
	return t.buf
}

func (p *pp) reset() {
	p.buf = p.buf[:0]
	p.arg = nil
	p.value = reflect.Value{}
	p.panicking = false
	p.erroring = false
	p.wrapErrs = false
	p.wrappedErr = nil

	p.fmt.init(&p.buf)
}
