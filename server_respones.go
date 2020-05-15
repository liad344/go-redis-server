package main

const (
	INT         = ":"
	STRING      = "+"
	ERROR       = "-"
	BULK_STRING = "$"
	ARRAY       = "*"
)

func (c *Conn) FWrite(prefix string, data []byte) []byte {
	b := []byte(prefix)
	b = append(b, data...)
	b = append(b, '\r', '\n')
	return b
}

func (c *Conn) WriteBulk(s string) {
}
func (c *Conn) WriteArray(s string) {
}
