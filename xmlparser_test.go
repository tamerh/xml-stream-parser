package xmlparser

import (
	"bufio"
	"os"
	"testing"
)

func getparser(prop string) *XMLParser {

	return getparserFile("sample.xml", prop)
}

func getparserFile(filename, prop string) *XMLParser {

	file, _ := os.Open(filename)

	br := bufio.NewReader(file)

	p := NewXMLParser(br, prop)

	return p

}

func TestBasics(t *testing.T) {

	p := getparser("tag1")

	var results []*XMLElement
	for xml := range p.Stream() {
		results = append(results, xml)
	}
	if len(results) != 2 {
		panic("Test failed result must be 2")
	}

	if len(results[0].Childs) != 4 || len(results[1].Childs) != 4 {
		panic("Test failed")
	}
	// result 1
	if results[0].Attrs["att1"] != "<att0>" || results[0].Attrs["att2"] != "att0" {
		panic("Test failed")
	}

	if results[0].Childs["tag11"][0].Attrs["att1"] != "att0" {
		panic("Test failed")
	}

	if results[0].Childs["tag11"][0].InnerText != "InnerText110" {
		panic("Test failed")
	}

	if results[0].Childs["tag11"][1].InnerText != "InnerText111" {
		panic("Test failed")
	}

	if results[0].Childs["tag12"][0].Attrs["att1"] != "att0" {
		panic("Test failed")
	}

	if results[0].Childs["tag12"][0].InnerText != "" {
		panic("Test failed")
	}

	if results[0].Childs["tag13"][0].Attrs != nil && results[0].Childs["tag13"][0].InnerText != "InnerText13" {
		panic("Test failed")
	}

	if results[0].Childs["tag14"][0].Attrs != nil && results[0].Childs["tag14"][0].InnerText != "" {
		panic("Test failed")
	}

	//result 2
	if results[1].Attrs["att1"] != "<att1>" || results[1].Attrs["att2"] != "att1" {
		panic("Test failed")
	}

	if results[1].Childs["tag11"][0].Attrs["att1"] != "att1" {
		panic("Test failed")
	}

	if results[1].Childs["tag11"][0].InnerText != "InnerText2" {
		panic("Test failed")
	}

	if results[1].Childs["tag12"][0].Attrs["att1"] != "att1" {
		panic("Test failed")
	}

	if results[1].Childs["tag12"][0].InnerText != "" {
		panic("Test failed")
	}
	if results[1].Childs["tag13"][0].Attrs != nil && results[1].Childs["tag13"][0].InnerText != "InnerText213" {
		panic("Test failed")
	}

	if results[1].Childs["tag14"][0].Attrs != nil && results[1].Childs["tag14"][0].InnerText != "" {
		panic("Test failed")
	}

}

func TestTagWithNoChild(t *testing.T) {

	p := getparser("tag2")

	var results []*XMLElement
	for xml := range p.Stream() {
		results = append(results, xml)
	}
	if len(results) != 2 {
		panic("Test failed")
	}
	if results[0].Childs != nil || results[1].Childs != nil {
		panic("Test failed")
	}
	if results[0].Attrs["att1"] != "testattr<" || results[1].Attrs["att1"] != "testattr<2" {
		panic("Test failed")
	}
	// with inner text
	p = getparser("tag3")

	results = results[:0]
	for xml := range p.Stream() {
		results = append(results, xml)
	}

	if len(results) != 2 {
		panic("Test failed")
	}
	if results[0].Childs != nil || results[1].Childs != nil {
		panic("Test failed")
	}

	if results[0].Attrs != nil || results[0].InnerText != "tag31" {
		panic("Test failed")
	}

	if results[1].Attrs["att1"] != "testattr<2" || results[1].InnerText != "tag32 " {
		panic("Test failed")
	}

}

func TestTagWithSpaceAndSkipOutElement(t *testing.T) {

	p := getparser("tag4").SkipElements([]string{"skipOutsideTag"}).SkipOuterElements()

	var results []*XMLElement
	for xml := range p.Stream() {
		results = append(results, xml)
	}

	if len(results) != 1 {
		panic("Test failed")
	}

	if results[0].Childs["tag11"][0].Attrs["att1"] != "att0 " {
		panic("Test failed")
	}

	if results[0].Childs["tag11"][0].InnerText != "InnerText0 " {
		panic("Test failed")
	}

}

func TestQuote(t *testing.T) {

	p := getparser("quotetest")

	var results []*XMLElement
	for xml := range p.Stream() {
		results = append(results, xml)
	}

	if len(results) != 1 {
		panic("Test failed")
	}

	if results[0].Attrs["att1"] != "test" || results[0].Attrs["att2"] != "test\"" || results[0].Attrs["att3"] != "test'" {
		panic("Test failed")
	}

}

func TestSkip(t *testing.T) {

	p := getparser("tag1").SkipElements([]string{"tag11", "tag13"})

	var results []*XMLElement
	for xml := range p.Stream() {
		results = append(results, xml)
	}

	if len(results[0].Childs) != 2 {
		panic("Test failed")
	}

	if len(results[1].Childs) != 2 {
		panic("Test failed")
	}

	if results[0].Childs["tag11"] != nil {
		panic("Test failed")
	}

	if results[0].Childs["tag13"] != nil {
		panic("Test failed")
	}

	if results[1].Childs["tag11"] != nil {
		panic("Test failed")
	}

	if results[1].Childs["tag13"] != nil {
		panic("Test failed")
	}

}

func TestError(t *testing.T) {

	p := getparserFile("error.xml", "tag1")

	for xml := range p.Stream() {
		if xml.Err == nil {
			panic("It must give error")
		}
	}

}
func TestGetAllNodes(t *testing.T) {
	p := getparser("examples")
	for xml := range p.Stream() {
		nodes := xml.GetAllNodes("father.son.grandson")
		if len(nodes) != 8 {
			t.Errorf("Lenght of xml.GetAllNodes is not the expected \n\t Expected: %d \n\t Found: %d", 8, len(nodes))
		} else {
			values := []string{"grandson111", "grandson112", "grandson121", "grandson122", "grandson131", "grandson132", "grandson211", "grandson212"}
			for i, node := range nodes {
				if node.GetValue(".") != values[i] {
					t.Errorf("The value of the grandson %d doesn´t match with the expected \n\t Expected: %s \n\t Found: %s", i, values[i], node.GetValue("."))
				}
			}
		}
	}
}
func TestGetValue(t *testing.T) {
	var found string
	p := getparser("examples")
	for xml := range p.Stream() {
		found = xml.GetValue("@inittag")
		if found != "initial_attr" {
			t.Errorf("@inittag doesn´t match with expected \n\t Expected: %s \n\t Found: %s", "initial_attr", found)
		}
		found = xml.GetValue("tag1.tag11")
		if found != "InnerText110" {
			t.Errorf("tag1>tag11 doesn´t match with expected \n\t Expected: %s \n\t Found: %s", "InnerText110", found)
		}
		found = xml.GetValue("tag1.tag11[1]")
		if found != "InnerText111" {
			t.Errorf("tag1>tag11[1] doesn´t match with expected \n\t Expected: %s \n\t Found: %s", "InnerText111", found)
		}
		found = xml.GetValue("tag1[1].tag11")
		if found != "InnerText2" {
			t.Errorf("tag1[1]>tag11 doesn´t match with expected \n\t Expected: %s \n\t Found: %s", "InnerText2", found)
		}
		found = xml.GetValue("tag1[10].tag11")
		if found != "" {
			t.Errorf("tag1[10]>tag11 doesn´t match with expected \n\t Expected: %s \n\t Found: %s", "InnerText2", found)
		}
		found = xml.GetValue("tag1.tag11[10]")
		if found != "" {
			t.Errorf("tag1>tag11[10] doesn´t match with expected \n\t Expected: %s \n\t Found: %s", "InnerText2", found)
		}
		found = xml.GetValue("tag1.tag12@att1")
		if found != "att0" {
			t.Errorf("tag1>tag12>@att1 doesn´t match with expected \n\t Expected: %s \n\t Found: %s", "att1", found)
		}
		found = xml.GetValue("tag1[1].tag12@att1")
		if found != "att1" {
			t.Errorf("tag1[1]>tag12>@att1 doesn´t match with expected \n\t Expected: %s \n\t Found: %s", "att1", found)
		}
		found = xml.GetValue("tag1[1].tag12@missingatt")
		if found != "" {
			t.Errorf("tag1[1]>tag12>@missingatt doesn´t match with expected \n\t Expected: %s \n\t Found: %s", "att1", found)
		}
		found = xml.GetValue("missingtag.tag12.tag13")
		if found != "" {
			t.Errorf("missingtag>tag12>tag13 doesn´t match with expected \n\t Expected: %s \n\t Found: %s", "att1", found)
		}
		found = xml.GetValue("tag1[1].tag12.missingtag@att1")
		if found != "" {
			t.Errorf("tag1[1]>tag12>missingtag>@att1 doesn´t match with expected \n\t Expected: %s \n\t Found: %s", "att1", found)
		}

		node := xml.GetNode("tag1[1].tag13")
		found = node.GetValue(".")
		if found != "InnerText213" {
			t.Errorf("tag1[1]>tag13 doesn´t match with expected \n\t Expected: %s \n\t Found: %s", "InnerText213", found)
		}
	}

}
func TestGetValueNumeric(t *testing.T) {
	var i int
	var f float64
	p := getparser("examples")
	for xml := range p.Stream() {
		i = xml.GetValueInt("numeric.int")
		if i != 8 {
			t.Errorf("numeric.int doesn´t match with expected \n\t Expected: %d \n\t Found: %d", 8, i)
		}
		i = xml.GetValueInt("numeric.int[1]")
		if i != 18 {
			t.Errorf("numeric.int[1] doesn´t match with expected \n\t Expected: %d \n\t Found: %d", 18, i)
		}
		i = xml.GetValueInt("numeric.int[2]")
		if i != 0 {
			t.Errorf("numeric.int[2] doesn´t match with expected \n\t Expected: %d \n\t Found: %d", 0, i)
		}
		i = xml.GetValueInt("numeric.int[2]@realInt")
		if i != 9 {
			t.Errorf("numeric.int[2]@realInt doesn´t match with expected \n\t Expected: %d \n\t Found: %d", 9, i)
		}
		f = xml.GetValueF64("numeric.float")
		if f != 39.9 {
			t.Errorf("numeric.float doesn´t match with expected \n\t Expected: %f \n\t Found: %f", 39.9, f)
		}
	}
}
func TestGetValueDeep(t *testing.T) {
	p := getparser("examples")
	for xml := range p.Stream() {
		i := xml.GetValueIntDeep("numericDeep.deep.int")
		if i != 2 {
			t.Errorf("numericDeep.deep.int doesn´t match with expected \n\t Expected: %d \n\t Found: %d", 2, i)
		}
		f := xml.GetValueF64Deep("numericDeep.deep.float")
		if f != 1.2 {
			t.Errorf("numericDeep.deep.float doesn´t match with expected \n\t Expected: %f \n\t Found: %f", 1.2, f)
		}
		v := xml.GetValueDeep("father.son.grandson")
		if v != "grandson111" {
			t.Errorf("father.son.grandson Deep search doesn´t match with expected \n\t Expected: %s \n\t Found: %s", "grandson111", v)
		}
	}
}
func Benchmark1(b *testing.B) {

	for n := 0; n < b.N; n++ {
		p := getparser("tag4").SkipElements([]string{"skipOutsideTag"}).SkipOuterElements()
		for xml := range p.Stream() {
			nothing(xml)
		}
	}
}

func Benchmark2(b *testing.B) {

	for n := 0; n < b.N; n++ {
		p := getparser("tag4")
		for xml := range p.Stream() {
			nothing(xml)
		}
	}

}

func nothing(...interface{}) {
}
