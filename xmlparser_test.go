package xmlparser

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func getparser(prop ...string) *XMLParser {

	return getparserFile("sample.xml", prop...)
}

func getparserFile(filename string, prop ...string) *XMLParser {

	file, _ := os.Open(filename)

	br := bufio.NewReader(file)

	p := NewXMLParser(br, prop...)

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

	if results[0].Childs["tag11"][0].InnerText != "Hello                                                你好            Gür" {
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

func TestMultipleTags(t *testing.T) {
	p := getparser("tag1", "tag2")

	tagCount := map[string]int{}
	for xml := range p.Stream() {
		if xml.Name != "tag1" && xml.Name != "tag2" {
			t.Errorf("Only 'tag1' and 'tag2' expected, but '%s' returned", xml.Name)
		}
		tagCount[xml.Name]++
	}

	if tagCount["tag1"] != 2 {
		t.Errorf("There should be 2 parsed 'tag1', but %d found", tagCount["tag1"])
	}
	if tagCount["tag2"] != 2 {
		t.Errorf("There should be 2 parsed 'tag2', but %d found", tagCount["tag2"])
	}
}

func TestMultipleTagsNested(t *testing.T) {
	p := getparser("tag1", "tag11")

	tagCount := map[string]int{}
	for xml := range p.Stream() {
		if xml.Name != "tag1" && xml.Name != "tag11" {
			t.Errorf("Only 'tag1' and 'tag11' expected, but '%s' returned", xml.Name)
		}
		tagCount[xml.Name]++
	}

	if tagCount["tag1"] != 2 {
		t.Errorf("There should be 2 parsed 'tag1', but %d found", tagCount["tag1"])
	}
	if tagCount["tag11"] != 1 {
		if tagCount["tag11"] == 4 {
			t.Errorf("There should be only 1 parsed 'tag11', but 'tag11' nested under 'tag1' were parsed too")
		}
		t.Errorf("There should be 1 parsed 'tag11', but %d found", tagCount["tag11"])
	}
}

func TestXpath(t *testing.T) {
	xmlDoc := `<?xml version="1.0" encoding="UTF-8"?>
	<bookstore>
		 <book id="bk101">
				<title>The Iliad and The Odyssey</title>
				<price>12.95</price>
				<comments>
					 <userComment rating="4">Best translation I've read.</userComment>
					 <userComment rating="2">I like other versions better.</userComment>
				</comments>
		 </book>
		 <book id="bk102">
				<title>Anthology of World Literature</title>
				<price>24.95</price>
				<comments>
					 <userComment rating="3">Needs more modern literature.</userComment>
					 <userComment rating="4">Excellent overview of world literature.</userComment>
				</comments>
		 </book>
		 <journal>
				<title>Journal of XML parsing</title>
				<issue>1</issue>
		 </journal>
	</bookstore>`

	sreader := strings.NewReader(xmlDoc)

	bufreader := bufio.NewReader(sreader)

	p := NewXMLParser(bufreader, "bookstore").EnableXpath()

	for xml := range p.Stream() {

		if list, err := xml.SelectElements("//book"); len(list) != 2 || err != nil {
			t.Fatal("//book != 2")
		}

		if list, err := xml.SelectElements("./book"); len(list) != 2 || err != nil {
			t.Fatal("./book != 2")
		}

		if list, err := xml.SelectElements("book"); len(list) != 2 || err != nil {
			t.Fatal("book != 2")
		}

		list, err := xml.SelectElements("./book/title")
		if len(list) != 2 || err != nil {
			t.Fatal("book != 2")
		}

		title, err := xml.SelectElement("./book/title")
		if err != nil && title.InnerText != "The Iliad and The Odyssey" {
			t.Fatal("./book/title")
		}

		el, err := xml.SelectElement("//book[@id='bk101']")
		if el == nil || err != nil {
			t.Fatal("//book[@id='bk101] is not found")
		}
		list, err = xml.SelectElements("//book[price>=10.95]")
		if list == nil || err != nil || len(list) != 2 {
			t.Fatal("//book[price>=10.95]")
		}

		list, err = xml.SelectElements("//book/comments/userComment[@rating='2']")
		if len(list) != 1 || err != nil {
			t.Fatal("//book/comments/userComment[@rating='2']")
		}

		// all books total price
		expr, err := p.CompileXpath("sum(//book/price)")
		if err != nil {
			t.Fatal("sum(//book/price) xpath expression compile error")
		}
		price := expr.Evaluate(p.CreateXPathNavigator(xml)).(float64)

		if fmt.Sprintf("%.2f", price) != "37.90" {
			t.Fatal("invalid total price->", price)
		}

	}
}

func TestXpathNS(t *testing.T) {

	br := bufio.NewReader(bytes.NewReader([]byte(`
		<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
		<soap:Body>
			<soap:BodyNest1>
				<soap:BodyNest2>
				</soap:BodyNest2>
			</soap:BodyNest1>
		</soap:Body>

		<soap:Body>
		<soap:BodyNest1>
			<soap:BodyNest3 nestatt3="nestatt3val">
			</soap:BodyNest3>
		</soap:BodyNest1>
	</soap:Body>

		</soap:Envelope>
	`)))

	str := NewXMLParser(br, "soap:Envelope").EnableXpath()
	for xml := range str.Stream() {

		if list, err := xml.SelectElements("soap:Body"); len(list) != 2 || err != nil {
			t.Fatal("soap:Body != 2")
		}

		if list, err := xml.SelectElements("./soap:Body/soap:BodyNest1"); len(list) != 2 || err != nil {
			t.Fatal("/soap:Body/soap:BodyNest1 != 2")
		}

		if list, err := xml.SelectElements("./soap:Body/soap:BodyNest1/soap:BodyNest2"); len(list) != 1 || err != nil {
			t.Fatal("/soap:Body/soap:BodyNest1/soap:BodyNest2 != 1")
		}

		list, err := xml.SelectElements("./soap:Body/soap:BodyNest1/soap:BodyNest3")
		if len(list) != 1 || err != nil {
			t.Fatal("/soap:Body/soap:BodyNest1/soap:BodyNest3 != 1")
		}

		if list[0].Attrs["nestatt3"] != "nestatt3val" {
			t.Fatal("nestatt3 attiribute test failed")
		}

	}

}

func TestAttrOnly(t *testing.T) {
	p := getparser("examples", "tag1").ParseAttributesOnly("examples")
	for xml := range p.Stream() {
		if xml.Err != nil {
			t.Fatal(xml.Err)
		}
		if xml.Name == "examples" {
			if len(xml.Childs) != 0 {
				t.Fatal("Childs not empty for ParseAttributesOnly tags")
			}
			fmt.Printf("Name: \t%s\n", xml.Name)
			fmt.Printf("Attrs: \t%v\n\n", xml.Attrs)
		}
		if xml.Name == "tag1" {
			if len(xml.Childs) == 0 {
				t.Fatal("Childs not empty for ParseAttributesOnly tags")
			}
			fmt.Printf("Name: \t%s\n", xml.Name)
			fmt.Printf("Attrs: \t%v\n", xml.Attrs)
			fmt.Printf("Childs: %v\n", xml.Childs)
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

func Benchmark3(b *testing.B) {

	for n := 0; n < b.N; n++ {
		p := getparser("tag4").EnableXpath()
		for xml := range p.Stream() {
			nothing(xml)
		}
	}

}

func nothing(...interface{}) {
}
