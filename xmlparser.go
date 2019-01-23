package xmlparser

import (
	"bufio"
	"fmt"
)

const (
	findLoopTag int = iota
	findElement
)

const eof = rune(0)
const elementOpen rune = '<'
const elementClose rune = '>'
const slash rune = '/'
const equal = '='
const quote = '"'

/*
XMLParser parser/scrapper of xml file
For more improvment following can be done
1- skip tags inside the element for now a tag can be skipped only root element of looptag
2- to make it more parallel maybe first just get the looptag content and send it for processing.
3- change slices size and append if applicabale.
*/
type XMLParser struct {
	R             *bufio.Reader
	LoopTag       string
	OutChannel    *chan XMLEntry
	SkipTags      []string
	FinishMessage string
	//	ProgBar       *mpb.Bar
	//	ProgByEntry   bool
	//	ProgBySize    bool
	//internal
	skipTagNames    map[string]bool
	state           int
	totalParsed     uint32
	progSizeCounter int32
}

// internal struct
type xmlTag struct {
	name         string
	attrs        map[string]string
	specialClose bool
}

/*
XMLEntry is a result of each parsed loop
*/
type XMLEntry struct {
	Attrs    map[string]string
	Elements map[string][]XMLElement
}

/*
XMLElement is a typical xml elements which keeps the parsed data
*/
type XMLElement struct {
	Attrs     map[string]string
	InnerText string
	Childs    map[string][]XMLElement
}

const errorMsg = "Parsing error check if document is valid or contact for help."
const errorMsg2 = " tag is not properly closed."
const errorMsg3 = "Main loop tag must have element inside."

/*
Parse starts parsing xml document
*/
func (x *XMLParser) Parse() {

	x.init()
	var c, n rune
	var entry *XMLEntry
	//start := time.Now()
	//progresEntryCount := 0
	//maxProgressEntry := 1000
	for {

		switch x.state {

		case findLoopTag:

			for {
				c = x.read()

				if c == eof {
					// finish
					if len(x.FinishMessage) > 0 {
						fmt.Println("Parsing done ", x.FinishMessage, " total parsed ->", x.totalParsed)
					}

					/**
					if x.ProgBar != nil {
						x.ProgBar.IncrBy(progresEntryCount)
					}
					**/
					close(*x.OutChannel)
					return
				}

				if c == elementOpen {
					atag := x.startTag()
					if atag.name == x.LoopTag {
						if atag.specialClose { // maybe main loop only has attribute???
							panic(errorMsg3 + "\n" + errorMsg)
						}
						entry = &XMLEntry{
							Attrs:    atag.attrs,
							Elements: map[string][]XMLElement{},
						}
						x.state = findElement
						break
					}

				}
			}

		case findElement:

			for {
				c = x.read()

				if c == eof {
					// exit completly
					if len(x.FinishMessage) > 0 {
						fmt.Println("Parsing done ", x.FinishMessage, " total parsed ->", x.totalParsed)
					}
					/**
					if x.ProgBar != nil {
						x.ProgBar.IncrBy(progresEntryCount)
					}
					**/

					close(*x.OutChannel)
					return
				}
				if c == elementOpen {

					n = x.read()
					//first check that if loop tag is closing
					if n == slash {
						close := x.closeTagName()
						if close == x.LoopTag {
							// loop tag is closing exit from this state
							x.state = findLoopTag
							*x.OutChannel <- *entry

							/**
							if x.ProgByEntry && progresEntryCount == maxProgressEntry {
								progresEntryCount++
								x.ProgBar.IncrBy(progresEntryCount)
								progresEntryCount = 0
							} else {
								progresEntryCount++
							}
							**/
							x.totalParsed++
							//start = time.Now()
							break
						} else { //this means some other tag is being close in loop tag
							continue
						}
					} else {
						x.unread()
					}

					atag := x.startTag()
					if _, ok := x.skipTagNames[atag.name]; !ok {

						// build element tree.
						childs := map[string][]XMLElement{}
						el := &XMLElement{
							Attrs:  atag.attrs,
							Childs: childs,
						}
						if !atag.specialClose {
							x.getElementTree(atag, el)
						}

						if _, ok = entry.Elements[atag.name]; ok {
							entry.Elements[atag.name] = append(entry.Elements[atag.name], *el)
						} else {
							var elements []XMLElement
							elements = append(elements, *el)
							entry.Elements[atag.name] = elements
						}

					} else { // we don't interested in this tag so move until end of it
						if !atag.specialClose {
							x.skipTag(atag)
						}
					}

				}

			}

		}

	}

}

func (x *XMLParser) init() {

	if x.OutChannel == nil {
		panic("Result channel is missing.")
	}

	x.skipTagNames = map[string]bool{}
	for _, tag := range x.SkipTags {
		x.skipTagNames[tag] = true
	}

	/**
	if x.ProgByEntry && x.ProgBySize {
		panic("progress bar must be based on either entry or size")
	}

	if x.ProgBar != nil && !x.ProgByEntry && !x.ProgBySize {
		x.ProgByEntry = true // this is default
	}
	**/

}

func (x *XMLParser) getElementTree(tag *xmlTag, result *XMLElement) *XMLElement {

	var c rune
	var n rune
	var innerText []rune

	for {
		c = x.read()
		if c == elementOpen {

			n = x.read()
			//first check if this is the end given tag
			if n == slash {
				close := x.closeTagName()
				if close == tag.name {
					//if there is no element and not special close getText
					if !tag.specialClose && len(result.Childs) == 0 {
						result.InnerText = string(innerText)
					}
					return result
				}
				//need this? else { x.unreadSize(len(close))}
			} else {
				x.unread()
			}

			currenttag := x.startTag()
			childs := map[string][]XMLElement{}
			currentElement := &XMLElement{
				Attrs:  currenttag.attrs,
				Childs: childs,
			}

			if !currenttag.specialClose {
				x.getElementTree(currenttag, currentElement)
			}

			if _, ok := result.Childs[currenttag.name]; ok {
				result.Childs[currenttag.name] = append(result.Childs[currenttag.name], *currentElement)
			} else {
				var childs []XMLElement
				childs = append(childs, *currentElement)
				result.Childs[currenttag.name] = childs
			}

		} else { // keep for innerText
			innerText = append(innerText, c)
		}
	}

}

func (x *XMLParser) skipTag(tag *xmlTag) {

	var c, n rune
	tagname := []rune(tag.name)
start:
	for {

		c = x.read()

		if c == elementOpen {
			n = x.read()
			if n == slash {

				for i := 0; i < len(tag.name); i++ {
					c = x.read()
					if c != tagname[i] {
						goto start
					}
				}

				for {
					c = x.read()
					if x.isWS(c) {
						continue
					}
					if c == elementClose {
						return
					}
				}

			}

		}
	}
}

func (x *XMLParser) closeTagName() string {

	var s []rune
	for {
		c := x.read()
		if c == elementClose {
			return string(s)
		}
		if !x.isWS(c) {
			s = append(s, c)
		}
	}
}

func (x *XMLParser) startTag() *xmlTag {

	var result = &xmlTag{}

	//1- get tag name
	// a tag have 3 forms * <abc> ** <abc type="foo" val="bar"/> *** <abc/>
	var s []rune
	var tagname string
	var c rune
	var prev = rune(0)
	var alreadyClosed bool
	var alreadySpecialClosed bool
	for {
		c = x.read()
		if c == eof {
			///TODO return error.
			//return
		}
		if x.isWS(c) { //form2
			tagname = string(s)
			break
		}

		if c == elementClose { //form1 and form3
			if prev == slash {
				tagname = string(s[:len(s)-1])
				alreadySpecialClosed = true
			} else {
				tagname = string(s)
			}
			alreadyClosed = true
			break
		}

		s = append(s, c)
		prev = c
	}

	result.name = tagname

	if alreadyClosed {
		result.specialClose = alreadySpecialClosed
		return result
	}

	//2- is special closew
	var in []rune
	var prevRune = rune(0)
	var specailClose bool
	for {
		c := x.read()
		if c != elementClose {
			in = append(in, c)
		} else {
			if prevRune == slash {
				specailClose = true
			}
			break
		}
		if c == eof {
			///x.exit()? TODO
			///return
		}
		prevRune = c
	}
	result.specialClose = specailClose

	//3- get attributes if needed
	if len(in) > 3 {

		var r = map[string]string{}
		var lastAttrEnd = 0

		for i := 1; i < len(in)-1; {
			if in[i] == equal && in[i+1] == quote {
				key := string(x.removeWS(in[lastAttrEnd:i]))
				valStartIndex := i + 2
				for i = i + 2; i < len(in); i++ {
					if in[i] == quote {
						r[key] = string(in[valStartIndex:i])
						lastAttrEnd = i + 1
						break
					}
				}
			} else {
				i++
			}

		}
		result.attrs = r
	}

	return result

}

func (*XMLParser) removeWS(in []rune) []rune {

	for i := len(in) - 1; i >= 0; i-- {
		if in[i] == ' ' || in[i] == '\t' {
			in = append(in[:i], in[i+1:]...)
		}
	}
	return in

}

func (*XMLParser) isWS(in rune) bool {

	if in == ' ' || in == '\t' {
		return true
	}

	return false

}

func (x *XMLParser) read() rune {
	ch, _, err := x.R.ReadRune()

	/**
	if x.ProgBySize {

		if x.progSizeCounter == 10000 { // to minimize the performance affect
			x.ProgBar.IncrBy(10000)
			x.progSizeCounter = 0
		} else {
			x.progSizeCounter++
		}

	}
	**/
	if err != nil {
		return eof
	}
	return ch
}

func (x *XMLParser) unread() {
	err := x.R.UnreadRune()
	if err != nil {
		panic(errorMsg)
	}
}

func (x *XMLParser) unreadSize(size int) {

	for size == 0 {

		err := x.R.UnreadRune()
		if err != nil {
			panic(errorMsg)
		}
		size--
	}

}
