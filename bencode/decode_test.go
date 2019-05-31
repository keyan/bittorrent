package bencode

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadIntUntil(t *testing.T) {
	reader := strings.NewReader("3:cat")
	d := decoder{*bufio.NewReader(reader)}

	val, err := d.readIntUntil(':')
	assert.Nil(t, err)
	assert.Equal(t, uint64(3), val)
}

func TestReadString(t *testing.T) {
	reader := strings.NewReader("3:cat")
	d := decoder{*bufio.NewReader(reader)}

	val, err := d.readString()
	assert.Nil(t, err)
	assert.Equal(t, "cat", val)
}

func TestReadList(t *testing.T) {
	reader := strings.NewReader("l4:spam4:eggse")
	d := decoder{*bufio.NewReader(reader)}

	val, err := d.readInterface()
	assert.Nil(t, err)
	assert.Equal(t, []interface{}{"spam", "eggs"}, val)
}

func TestDecode(t *testing.T) {
	reader := strings.NewReader("d3:cow3:moo4:spam4:eggse")

	val, err := Decode(reader)
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{"cow": "moo", "spam": "eggs"}, val)
}
