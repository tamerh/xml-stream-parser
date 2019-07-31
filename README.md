## xml stream parser 
xml-stream-parser is xml parser for GO. It is efficient to parse large xml data with streaming fashion. 

### Usage

```xml
<?xml version="1.0" encoding="UTF-8"?>
<bookstore>
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
</bookstore>
```

<b>Stream</b> over books
```go


f, _ := os.Open("input.xml")
br := bufio.NewReaderSize(f,65536)
parser := xmlparser.NewXMLParser(br, "book")

for xml := range parser.Stream() {
	fmt.Println(xml.Childs["title"][0].InnerText)
	fmt.Println(xml.Childs["comments"][0].Childs["userComment"][0].Attrs["rating"])
	fmt.Println(xml.Childs["comments"][0].Childs["userComment"][0].InnerText)
}
   
```

<b>Skip</b> tags for speed
```go
parser := xmlparser.NewXMLParser(br, "book").SkipElements([]string{"price", "comments"})
```

<b>Error</b> handlings
```go
for xml := range parser.Stream() {
   if xml.Err !=nil { 
      // handle error
   }
}
```

<b>Progress</b> of parsing
```go
// total byte read to calculate the progress of parsing
parser.TotalReadSize
```

<b>Using GetValue</b> function from a XMLElement instance:
```
func TestGetValue(t *testing.T) {
	var found string
	p := getparser("examples")
	for xml := range p.Stream() {
		found = xml.GetValue("tag1.tag11")
   }
}
```
and to get an attribute value:
```
found = xml.GetValue("tag1.tag12@att0")
```
even you can get the value of a different node than the first using indexation like this:
```
found = xml.GetValue("tag1[1].tag11@att0")
```
_note: this index works in the same way that the indexation of the arrays in functional language, that means that it won´t work as the indexation in the xslt. In this way 0 will be the fisrt element, 1 the second, and so on._
The last point about this function: it will return an empty string when you ask for an non existing node, a index out of the rage, or a non existing attribute, so all of this invokations:
```
found = xml.GetValue("missingTag.tag11@att")
found = xml.GetValue("tag1[1].missingTag@att0")
found = xml.GetValue("tag1[999999999].missingTag@att0")
found = xml.GetValue("tag1[1].tag11@missingAttribute")
```
will generate an empty string.


If you interested check also [json parser](https://github.com/tamerh/jsparser) which works similarly
