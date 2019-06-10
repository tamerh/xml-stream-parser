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

	start := time.Now()

	file, _ := os.Open("uniprot_test.xml.gz")

	defer file.Close()

	gz, _ := gzip.NewReader(file)

	br := bufio.NewReader(gz)

	p := NewXmlParser(br, "entry").SkipTags([]string{"comment", "gene", "protein", "feature", "sequence"})

	totalentry := 0
	for range *p.Stream() {
		totalentry++
	}

	if totalentry != 2 {
		panic("Expected result count is 2 but found ->" + string(totalentry))
	}

	elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)

}

func TestXml(t *testing.T) {

	file, _ := os.Open("books.xml")
	defer file.Close()

	br := bufio.NewReader(file)

	p := NewXmlParser(br, "book")

	var resultEntryCount int
	for range *p.Stream() {
		resultEntryCount++
	}

	if resultEntryCount != 12 {
		panic("Expected result count is 12 but found ->" + string(resultEntryCount))
	}

}

func TestXmlUniref(t *testing.T) {

	file, _ := os.Open("uniref.xml")
	defer file.Close()

	br := bufio.NewReader(file)

	p := NewXmlParser(br, "entry")

	var resultEntryCount int
	for range *p.Stream() {
		resultEntryCount++
	}

	if resultEntryCount != 1 {
		panic("Expected result count is 1 but found ->" + string(resultEntryCount))
	}

}

func TestArticleXml(t *testing.T) {

	file, _ := os.Open("article.xml")
	defer file.Close()

	br := bufio.NewReader(file)

	p := NewXmlParser(br, "article-meta")

	for entry := range *p.Stream() {
		if len(entry.Elements["article-id"]) != 3 {
			panic("Article should have 3 article id  ")
		}
	}

}

func TestTaxonomyXml(t *testing.T) {

	file, _ := os.Open("taxonomy.xml")
	defer file.Close()

	br := bufio.NewReader(file)
	p := NewXmlParser(br, "taxon")

	for entry := range *p.Stream() {
		fmt.Println(entry)
	}

}

func TestHmdb(t *testing.T) {

	file, _ := os.Open("hmdb.xml")
	defer file.Close()

	br := bufio.NewReader(file)

	p := NewXmlParser(br, "metabolite").SkipTags([]string{"taxonomy", "ontology"})

	entrycount := 0
	for range *p.Stream() {
		entrycount++
	}

	if entrycount != 1 {
		panic("Expected result count is 1 but found ->" + string(entrycount))
	}

}

func TestBooks2Xml(t *testing.T) {

	file, _ := os.Open("books2.xml")
	defer file.Close()

	br := bufio.NewReader(file)

	p := NewXmlParser(br, "book").SkipTags([]string{"description"})

	for book := range *p.Stream() {

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

	file, _ := os.Open("books_invalid.xml")
	defer file.Close()

	br := bufio.NewReader(file)

	p := NewXmlParser(br, "book").SkipTags([]string{"description"})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r == nil {
				panic("The code did not panic")
			} else {
				os.Exit(0)
			}
			wg.Wait()
		}()

	}()

	for range *p.Stream() {
	}
	wg.Done()

}
