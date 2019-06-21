package xmlparser

import (
	"bufio"
	"fmt"
	"unicode/utf8"
)

type XMLParser struct {
	reader        *bufio.Reader
	loopElement   string
	resultChannel chan *XMLElement
	skipElements  map[string]bool
	scratch       *scratch
	scratch2      *scratch
	TotalReadSize uint64
}

type XMLElement struct {
	Attrs     map[string]string
	InnerText string
	Childs    map[string][]XMLElement
	Err       error
}

func NewXMLParser(reader *bufio.Reader, loopElement string) *XMLParser {

	x := &XMLParser{
		reader:        reader,
		loopElement:   loopElement,
		resultChannel: make(chan *XMLElement, 256),
		skipElements:  map[string]bool{},
		scratch:       &scratch{data: make([]byte, 1024)},
		scratch2:      &scratch{data: make([]byte, 1024)},
	}
	return x
}

func (x *XMLParser) SkipElements(skipElements []string) *XMLParser {

	if len(skipElements) > 0 {
		for _, s := range skipElements {
			x.skipElements[s] = true
		}
	}
	return x

}

func (x *XMLParser) Stream() chan *XMLElement {

	go x.parse()

	return x.resultChannel

}

func (x *XMLParser) parse() {

	defer close(x.resultChannel)
	var element *XMLElement
	var tagName string
	var tagClosed bool
	var err error
	var b byte
	var iscomment bool

	err = x.skipDeclerations()

	if err != nil {
		x.sendError()
		return
	}

	for {
		b, err = x.readByte()

		if err != nil {
			return
		}

		if x.isWS(b) {
			continue
		}

		if b == '<' {

			iscomment, err = x.isComment()

			if err != nil {
				x.sendError()
				return
			}

			if iscomment {
				continue
			}

			tagName, element, tagClosed, err = x.startElement()

			if err != nil {
				x.sendError()
				return
			}

			if tagName == x.loopElement {
				if tagClosed {
					x.resultChannel <- element
					continue
				}

				element = x.getElementTree(tagName, element)
				x.resultChannel <- element
				if element.Err != nil {
					return
				}

			}

		}
	}

}

func (x *XMLParser) getElementTree(tagName string, result *XMLElement) *XMLElement {

	if result.Err != nil {
		return result
	}

	var cur byte
	var next byte
	var err error
	var element *XMLElement
	var tagClosed bool
	x.scratch2.reset() // this hold the inner text
	var tagName2 string
	var iscomment bool

	for {

		cur, err = x.readByte()

		if err != nil {
			result.Err = err
			return result
		}

		if cur == '<' {

			iscomment, err = x.isComment()

			if err != nil {
				result.Err = err
				return result
			}

			if iscomment {
				continue
			}

			next, err = x.readByte()

			if err != nil {
				result.Err = err
				return result
			}

			if next == '/' { // close tag
				tag, err := x.closeTagName()

				if err != nil {
					result.Err = err
					return result
				}

				if tag == tagName {
					if len(result.Childs) == 0 { // check special tag???
						result.InnerText = string(x.scratch2.bytes())
					}
					return result
				}
			} else {
				x.unreadByte()
			}

			tagName2, element, tagClosed, err = x.startElement()

			if err != nil {
				result.Err = err
				return result
			}

			if _, ok := x.skipElements[tagName2]; ok && !tagClosed {
				x.skipElement(tagName2)
				continue
			}
			if !tagClosed {
				element = x.getElementTree(tagName2, element)
			}

			if _, ok := result.Childs[tagName2]; ok {
				result.Childs[tagName2] = append(result.Childs[tagName2], *element)
			} else {
				var childs []XMLElement
				childs = append(childs, *element)
				if result.Childs == nil {
					result.Childs = map[string][]XMLElement{}
				}
				result.Childs[tagName2] = childs
			}

		} else {
			x.scratch2.add(cur)
		}

	}
}

func (x *XMLParser) skipElement(elname string) error {

	var c byte
	var next byte
	var err error
	var curname string
	for {

		c, err = x.readByte()

		if err != nil {
			return err
		}
		if c == '<' {

			next, err = x.readByte()

			if err != nil {
				return err
			}

			if next == '/' {
				curname, err = x.closeTagName()
				if err != nil {
					return err
				}
				if curname == elname {
					return nil
				}
			}

		}

	}
}

func (x *XMLParser) startElement() (string, *XMLElement, bool, error) {

	x.scratch.reset()

	var cur byte
	var prev byte
	var err error
	var result = &XMLElement{}
	// a tag have 3 forms * <abc > ** <abc type="foo" val="bar"/> *** <abc />
	var attr string
	var attrVal string
	var tagName string
	for {

		cur, err = x.readByte()

		if err != nil {
			return "", nil, false, x.defaultError()
		}

		if x.isWS(cur) {
			tagName = string(x.scratch.bytes())
			x.scratch.reset()
			goto search_close_tag
		}

		if cur == '>' {
			if prev == '/' {
				return string(x.scratch.bytes()[:len(x.scratch.bytes())-1]), result, true, nil
			}
			return string(x.scratch.bytes()), result, false, nil
		}
		x.scratch.add(cur)
		prev = cur
	}

search_close_tag:
	for {

		cur, err = x.readByte()

		if err != nil {
			return "", nil, false, x.defaultError()
		}

		if x.isWS(cur) {
			continue
		}

		if cur == '=' {
			if result.Attrs == nil {
				result.Attrs = map[string]string{}
			}

			cur, err = x.readByte()

			if err != nil {
				return "", nil, false, x.defaultError()
			}

			if !(cur == '"' || cur == '\'') {
				return "", nil, false, x.defaultError()
			}

			attr = string(x.scratch.bytes())
			attrVal, err = x.string(cur)
			if err != nil {
				return "", nil, false, x.defaultError()
			}
			result.Attrs[attr] = attrVal
			x.scratch.reset()
			continue
		}

		if cur == '>' { //if tag name not found
			if prev == '/' { //tag special close
				return tagName, result, true, nil
			}
			return tagName, result, false, nil
		}

		x.scratch.add(cur)
		prev = cur

	}

}

func (x *XMLParser) isComment() (bool, error) {

	var c byte
	var err error

	c, err = x.readByte()

	if err != nil {
		return false, err
	}

	if c != '!' {
		x.unreadByte()
		return false, nil
	}

	var d, e byte

	d, err = x.readByte()

	if err != nil {
		return false, err
	}

	e, err = x.readByte()

	if err != nil {
		return false, err
	}

	if d != '-' || e != '-' {
		err = x.defaultError()
		return false, err
	}

	// skip part
	x.scratch.reset()
	for {

		c, err = x.readByte()

		if err != nil {
			return false, err
		}

		if c == '>' && len(x.scratch.bytes()) > 1 && x.scratch.bytes()[len(x.scratch.bytes())-1] == '-' && x.scratch.bytes()[len(x.scratch.bytes())-2] == '-' {
			return true, nil
		}

		x.scratch.add(c)

	}

}

func (x *XMLParser) skipDeclerations() error {

	var a, b []byte
	var c, d byte
	var err error

scan_declartions:
	for {

		// when identifying a xml declaration we need to know 2 bytes ahead. Unread works 1 byte at a time so we use Peek and read together.
		a, err = x.reader.Peek(1)

		if err != nil {
			return err
		}

		if a[0] == '<' {

			b, err = x.reader.Peek(2)

			if err != nil {
				return err
			}

			if b[1] == '!' || b[1] == '?' { // either comment or decleration

				// read 2 peaked byte
				_, err = x.readByte()

				if err != nil {
					return err
				}

				_, err = x.readByte()
				if err != nil {
					return err
				}

				c, err = x.readByte()

				if err != nil {
					return err
				}

				d, err = x.readByte()

				if err != nil {
					return err
				}

				if c == '-' && d == '-' {
					goto skipComment
				} else {
					goto skipDecleration
				}

			} else { // declerations ends.

				return nil

			}

		}

		// read peaked byte
		_, err = x.readByte()

		if err != nil {
			return err
		}

	}

skipComment:
	x.scratch.reset()
	for {

		c, err = x.readByte()

		if err != nil {
			return err
		}

		if c == '>' && len(x.scratch.bytes()) > 1 && x.scratch.bytes()[len(x.scratch.bytes())-1] == '-' && x.scratch.bytes()[len(x.scratch.bytes())-2] == '-' {
			goto scan_declartions
		}

		x.scratch.add(c)

	}

skipDecleration:
	depth := 1
	for {

		c, err = x.readByte()

		if err != nil {
			return err
		}

		if c == '>' {
			depth--
			if depth == 0 {
				goto scan_declartions
			}
			continue
		}
		if c == '<' {
			depth++
		}

	}

}

func (x *XMLParser) closeTagName() (string, error) {

	x.scratch.reset()
	var c byte
	var err error
	for {
		c, err = x.readByte()

		if err != nil {
			return "", err
		}

		if c == '>' {
			return string(x.scratch.bytes()), nil
		}
		if !x.isWS(c) {
			x.scratch.add(c)
		}
	}
}

func (x *XMLParser) readByte() (byte, error) {

	by, err := x.reader.ReadByte()

	x.TotalReadSize++

	if err != nil {
		return 0, err
	}
	return by, nil

}

func (x *XMLParser) unreadByte() error {

	err := x.reader.UnreadByte()
	if err != nil {
		return err
	}
	x.TotalReadSize = x.TotalReadSize - 1
	return nil

}

func (x *XMLParser) isWS(in byte) bool {

	if in == ' ' || in == '\n' || in == '\t' || in == '\r' {
		return true
	}

	return false

}

func (x *XMLParser) sendError() {
	err := fmt.Errorf("Invalid xml")
	x.resultChannel <- &XMLElement{Err: err}
}

func (x *XMLParser) defaultError() error {
	err := fmt.Errorf("Invalid xml")
	return err
}

func (x *XMLParser) string(start byte) (string, error) {

	x.scratch.reset()

	var err error
	var c byte
	for {

		c, err = x.readByte()
		if err != nil {
			if err != nil {
				return "", err
			}
		}

		if c == start {
			return string(x.scratch.bytes()), nil
		}

		x.scratch.add(c)

	}

}

// scratch taken from
//https://github.com/bcicen/jstream
// string called by `any` or `object`(for map keys) after reading `"`

type scratch struct {
	data []byte
	fill int
}

// reset scratch buffer
func (s *scratch) reset() { s.fill = 0 }

// bytes returns the written contents of scratch buffer
func (s *scratch) bytes() []byte { return s.data[0:s.fill] }

// grow scratch buffer
func (s *scratch) grow() {
	ndata := make([]byte, cap(s.data)*2)
	copy(ndata, s.data[:])
	s.data = ndata
}

// append single byte to scratch buffer
func (s *scratch) add(c byte) {
	if s.fill+1 >= cap(s.data) {
		s.grow()
	}

	s.data[s.fill] = c
	s.fill++
}

// append encoded rune to scratch buffer
func (s *scratch) addRune(r rune) int {
	if s.fill+utf8.UTFMax >= cap(s.data) {
		s.grow()
	}

	n := utf8.EncodeRune(s.data[s.fill:], r)
	s.fill += n
	return n
}
