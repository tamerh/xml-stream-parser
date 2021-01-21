# xml stream parser

xml-stream-parser is xml parser for GO. It is efficient to parse large xml data with streaming fashion.

## Usage

```xml
<?xml version="1.0" encoding="UTF-8"?>
<bookstore number="2" loc="273456">
   <book>
      <title>The Iliad and The Odyssey</title>
      <price>12.95</price>
      <comments>
         <userComment rating="4">Best translation I've read.</userComment>
         <userComment rating="2">I like other versions better.</userComment>
      </comments>
   </book>
   <book>
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
</bookstore>
```

**Stream** over books and journals

```go
f, _ := os.Open("input.xml")
br := bufio.NewReaderSize(f,65536)
parser := xmlparser.NewXMLParser(br, "book", "journal")

for xml := range parser.Stream() {
   fmt.Println(xml.Childs["title"][0].InnerText)
   if xml.Name == "book" {
      fmt.Println(xml.Childs["comments"][0].Childs["userComment"][0].Attrs["rating"])
      fmt.Println(xml.Childs["comments"][0].Childs["userComment"][0].InnerText)
   }
}
```

**Skip** tags for speed

```go
parser := xmlparser.NewXMLParser(br, "book").SkipElements([]string{"price", "comments"})
```

**Attributes** only

```go
parser := xmlparser.NewXMLParser(br, "bookstore", "book").ParseAttributesOnly("bookstore")
```

**Error** handlings

```go
for xml := range parser.Stream() {
   if xml.Err !=nil {
      // handle error
   }
}
```

**Progress** of parsing

```go
// total byte read to calculate the progress of parsing
parser.TotalReadSize
```

**Xpath** query provides alternative to default fast access for different usecases 
```go

parser := xmlparser.NewXMLParser(bufreader, "bookstore").EnableXpath()

for xml := range p.Stream() {
   // select books 
   xml.SelectElements("//book")
   xml.SelectElements("./book")
   xml.SelectElements("book")
   // select titles
   xml.SelectElements("./book/title")
   // select book with price condition
   xml.SelectElements("//book[price>=20.95]"))
   //comments with rating 4
   xml.SelectElements("//book/comments/userComment[@rating='4']")
}
// for evaluate function or reuse existing xpath expression
// sum of all the book price
expr, err := p.CompileXpath("sum(//book/price)")
price := expr.Evaluate(p.CreateXPathNavigator(xml)).(float64)

```
xpath functionality implemented via [xpath](https://github.com/antchfx/xpath) library check more 
examples in its documentation

If you interested check also [json parser](https://github.com/tamerh/jsparser) which works similarly
