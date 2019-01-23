## XML Stream Parser for GO
xml-stream-parser is a GO library to parse xml files effectively. It is written to addres the performance [issue](https://github.com/golang/go/issues/21823) in default xml package when that issue is resolved you can skip this library.

Right now it requires to write a bit too much code to parse which will be simplified but it is working faster compare to default xml package. It also works with very low memory footprint. This library used in my [project](https://github.com/tamerh/biobtree) with large xml files like size of more than 100GB. But be aware that it may not cover every case create issue if you found any.

### Install

```
go get -u github.com/tamerh/xml-stream-parser
```


### Usage

Let say you have following xml and you want to loop over book as a stream
and parse various elements and attributes

```xml
<?xml version="1.0" encoding="UTF-8"?>
<bookstore>
   <book ISBN="10-000000-001">
      <title>The Iliad and The Odyssey</title>
      <price>12.95</price>
      <comments>
         <userComment rating="4">Best translation I've read.</userComment>
         <userComment rating="2">I like other versions better.</userComment>
      </comments>
      <description>Homer's two epics of the ancient world, The Iliad & The Odyssey, tell stories as riveting today as when they were written between the eighth and ninth century B.C.</description>
   </book>
   <book ISBN="10-000000-999">
      <title>Anthology of World Literature</title>
      <price>24.95</price>
      <comments>
         <userComment rating="3">Needs more modern literature.</userComment>
         <userComment rating="4">Excellent overview of world literature.</userComment>
      </comments>
      <description>The anthology includes epic and lyric poetry, drama, and prose narrative, with many complete works and a focus on the most influential pieces and authors from each region and time period.</description>
   </book>
</bookstore>
```

you can use the library like so

```go
//First open your file and create reader. You can also use gzip file check tests
file, _ := os.Open("books2.xml")
defer file.Close()
br := bufio.NewReader(file)

// then create  following channel to read your parsed data from.
var resultChannel = make(chan XMLEntry)

// init parser
var parser = XMLParser{
R:          br, 
// define tag to loop over
LoopTag:    "book",
OutChannel: &resultChannel,
// you can skip tags that you are not interested it relatively speeds up the process
SkipTags:   []string{"description"}, 
}

// start parsing with a go routine
go parser.Parse()

// and finally read parsed data 
for book := range resultChannel {
// print ISBN value
isbn := book.Attrs["ISBN"]
fmt.Println(isbn)

// print title
title := book.Elements["title"][0].InnerText
fmt.Println(title)

// print a user commet which has rating 4
// basically you can walk on all the sub nodes if you have
for _, userComments := range book.Elements["comments"][0].Childs {
	for _, comment := range userComments {
		if comment.Attrs["rating"] == "4" {
			  // print the user comment
			  fmt.Println(comment.InnerText)
			}
		}
	}
}
```
