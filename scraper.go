package main

import (
	"fmt"

	"encoding/json"
	"log"

	"time"

	"github.com/gocolly/colly"
)

type Item struct {
	Link	string `json:"link"`
	Name	string `json:"name"`
	Price	string `json:"price"`
	Instock	string `json:"instock"`
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func main() {
	defer timer("main")()
	c := colly.NewCollector(colly.Async(true))
	
	items := []Item{}

	c.OnHTML("div.side_categories li ul li", func(h *colly.HTMLElement) {
		link := h.ChildAttr("a", "href")
		c.Visit(h.Request.AbsoluteURL(link))
	})

	c.OnHTML("li.next a", func(h *colly.HTMLElement) {
		c.Visit(h.Request.AbsoluteURL(h.Attr("href")))
	})

	c.OnHTML("article.product_pod", func(h *colly.HTMLElement) {
		i := Item{
			Link:	h.ChildAttr("a", "href"),
			Name:	h.ChildAttr("h3 a", "title"),
			Price:	h.ChildText("p.price_color"),
			Instock:h.ChildText("p.instock"),
		}
		items = append(items, i)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting", r.URL)
	})

	// c.Visit("https://books.toscrape.com/")
	c.Visit("https://books.toscrape.com/catalogue/category/books/travel_2/index.html")
	c.Wait()

	data, err := json.MarshalIndent(items, " ", "")
	if err != nil{
		log.Fatal()
	}
	fmt.Println(string(data))
}