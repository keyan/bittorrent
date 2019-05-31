package bencode

import (
	"bufio"
	"errors"
	"io"
	"strconv"
)

type dictionary = map[string]interface{}

type decoder struct {
	bufio.Reader
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func (d *decoder) readIntUntil(until byte) (uint64, error) {
	value, err := d.ReadBytes(until)
	check(err)

	str := string(value[:len(value)-1])

	i, err := strconv.ParseUint(str, 10, 64)
	check(err)

	return i, nil
}

func (d *decoder) readString() (string, error) {
	numBytes, err := d.readIntUntil(':')
	check(err)

	buffer := make([]byte, numBytes)
	_, err = io.ReadFull(d, buffer)
	check(err)

	return string(buffer), nil
}

func (d *decoder) readList() ([]interface{}, error) {
	var list []interface{}

	for {
		nextByte, err := d.ReadByte()
		check(err)

		if nextByte == 'e' {
			break
		} else if err = d.UnreadByte(); err != nil {
			return nil, err
		}

		item, err := d.readInterface()
		check(err)

		list = append(list, item)
	}

	return list, nil
}

func (d *decoder) readInterface() (interface{}, error) {
	var value interface{}
	var err error

	nextByte, err := d.ReadByte()
	check(err)

	switch nextByte {
	case 'i':
		value, err = d.readIntUntil('e')
	case 'd':
		value, err = d.readDictionary()
	case 'l':
		value, err = d.readList()
	default:
		if err = d.UnreadByte(); err != nil {
			return nil, err
		}
		value, err = d.readString()
	}

	return value, err
}

func (d *decoder) readDictionary() (dictionary, error) {
	dict := make(dictionary)

	for {
		key, err := d.readString()
		check(err)

		value, err := d.readInterface()
		check(err)

		dict[key] = value

		nextByte, err := d.ReadByte()
		check(err)

		if nextByte == 'e' {
			break
		} else {
			d.UnreadByte()
		}
	}

	return dict, nil
}

func Decode(reader io.Reader) (dictionary, error) {
	d := decoder{*bufio.NewReader(reader)}

	firstByte, err := d.ReadByte()
	if err != nil {
		return make(dictionary), nil
	} else if firstByte != 'd' {
		return nil, errors.New("becode data must begin with a dictionary")
	}

	dict, err := d.readDictionary()
	check(err)

	return dict, nil
}
