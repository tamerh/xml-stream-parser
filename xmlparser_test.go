package xmlparser

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

func TestGzXml(t *testing.T) {

	var resultChannel = make(chan XMLEntry)

	var wg sync.WaitGroup

	wg.Add(1)
	file, _ := os.Open("uniprot_test.xml.gz")

	defer file.Close()

	gz, err := gzip.NewReader(file)

	if err != nil {
		fmt.Println("Error,", err)
	}

	br := bufio.NewReader(gz)

	var parser = XMLParser{
		R:          br,
		LoopTag:    "entry",
		OutChannel: &resultChannel,
		SkipTags:   []string{"comment", "gene", "protein", "feature", "sequence"},
	}

	start := time.Now()
	fmt.Println("Started...")
	go func() {
		parser.Parse()
		wg.Done()
	}()

	totalentry := 0
	for range resultChannel {
		totalentry++
	}

	if totalentry != 2 {
		panic("Expected result count is 2 but found ->" + string(totalentry))
	}

	elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)

}

func TestXml(t *testing.T) {

	var resultChannel = make(chan XMLEntry)

	file, _ := os.Open("books.xml")
	defer file.Close()

	br := bufio.NewReader(file)

	var parser = XMLParser{
		R:          br,
		LoopTag:    "book",
		OutChannel: &resultChannel,
	}

	go parser.Parse()

	var resultEntryCount int
	for range resultChannel {
		resultEntryCount++
	}

	if resultEntryCount != 12 {
		panic("Expected result count is 12 but found ->" + string(resultEntryCount))
	}

}

func TestXmlUniref(t *testing.T) {

	var resultChannel = make(chan XMLEntry)

	file, _ := os.Open("uniref.xml")
	defer file.Close()

	br := bufio.NewReader(file)

	var parser = XMLParser{
		R:          br,
		LoopTag:    "entry",
		OutChannel: &resultChannel,
	}

	go parser.Parse()

	var resultEntryCount int
	for range resultChannel {
		resultEntryCount++
	}

	if resultEntryCount != 1 {
		panic("Expected result count is 1 but found ->" + string(resultEntryCount))
	}

}

func TestArticleXml(t *testing.T) {

	var resultChannel = make(chan XMLEntry)

	file, _ := os.Open("article.xml")
	defer file.Close()

	br := bufio.NewReader(file)

	var parser = XMLParser{
		R:          br,
		LoopTag:    "article-meta",
		OutChannel: &resultChannel,
	}

	go parser.Parse()

	for entry := range resultChannel {
		if len(entry.Elements["article-id"]) != 3 {
			panic("Article should have 3 article id  ")
		}

	}

}

func TestTaxonomyXml(t *testing.T) {

	var resultChannel = make(chan XMLEntry)

	file, _ := os.Open("taxonomy.xml")
	defer file.Close()

	br := bufio.NewReader(file)

	var parser = XMLParser{
		R:          br,
		LoopTag:    "taxon",
		OutChannel: &resultChannel,
	}

	go parser.Parse()

	for entry := range resultChannel {
		fmt.Println(entry)
	}

}

func TestHmdb(t *testing.T) {

	var resultChannel = make(chan XMLEntry)

	file, _ := os.Open("hmdb.xml")
	defer file.Close()

	br := bufio.NewReader(file)

	var parser = XMLParser{
		R:          br,
		LoopTag:    "metabolite",
		SkipTags:   []string{"taxonomy,ontology"},
		OutChannel: &resultChannel,
	}

	go parser.Parse()

	entrycount := 0
	for range resultChannel {
		entrycount++
	}

	if entrycount != 1 {
		panic("Expected result count is 1 but found ->" + string(entrycount))
	}

}

func TestBooks2Xml(t *testing.T) {

	var resultChannel = make(chan XMLEntry)

	file, _ := os.Open("books2.xml")
	defer file.Close()

	br := bufio.NewReader(file)

	var parser = XMLParser{
		R:          br,
		LoopTag:    "book",
		OutChannel: &resultChannel,
		SkipTags:   []string{"description"},
	}

	go parser.Parse()

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

}

// this test must throw panic
func TestBooksInvalidXml(t *testing.T) {

	var resultChannel = make(chan XMLEntry)

	file, _ := os.Open("books_invalid.xml")
	defer file.Close()

	br := bufio.NewReader(file)

	var parser = XMLParser{
		R:          br,
		LoopTag:    "book",
		OutChannel: &resultChannel,
		SkipTags:   []string{"description"},
	}

	go func() {
		defer func() {
			if r := recover(); r == nil {
				panic("The code did not panic")
			} else {
				os.Exit(0)
			}
		}()
		parser.Parse()
	}()

	for book := range resultChannel {

		// print ISBN value
		isbn := book.Attrs["ISBN"]
		fmt.Println(isbn)

	}

}
