package main

import log "github.com/sirupsen/logrus"

func (c *Conn) WriteString(s string) {
	s = "+" + s
	n , err := c.Write([]byte(s))
	if n != len(s) || err != nil {
		log.Error("Could not write string " , err)
	}

}
func (c *Conn) WriteBulk(s string) {

}
func (c *Conn) WriteInt(s string) {

}
func (c *Conn) WriteArray(s string) {

}
func (c *Conn) WriteNull() {
	c.WriteString("\"$-1\\r\\n\"")
}

