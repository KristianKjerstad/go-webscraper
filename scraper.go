package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

func main() {

	scrapeEcommerceProducts()
}

// contains checks if a slice contains a specific element
func contains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

func scrapeEcommerceProducts() {

	websiteToScrape := "https://www.scrapingcourse.com/ecommerce/"

	var products []Product

	collector := colly.NewCollector()
	collector.OnRequest(requestHandler)
	collector.OnResponse(onResponseHandler)

	collector.OnHTML("li.product", func(e *colly.HTMLElement) {
		HTMLHandler(e, &products) // Passing the products slice as a reference
	})
	collector.OnScraped(func(r *colly.Response) {
		onScrapedHandler(r, &products)
	})

	collector.Visit(websiteToScrape)

	fmt.Println(collector)
	productsAsJSON, err := json.MarshalIndent(products, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling products: %v", err)
	}
	fmt.Println(string(productsAsJSON))
}

func requestHandler(request *colly.Request) {
	fmt.Println("Visiting: ", request.URL)

}

func onResponseHandler(response *colly.Response) {
	fmt.Println("Page visited", response.Request.URL)
}

// func HTMLPagesHandler(e *colly.HTMLElement, pagesToScrape *[]string, pagesDiscovered *[]string) {
// 	newPaginationLink := e.Attr("href")

// 	if !contains(*pagesToScrape, newPaginationLink) {
// 		if !contains(*pagesDiscovered, newPaginationLink) {
// 			*pagesToScrape = append(*pagesToScrape, newPaginationLink)
// 		}
// 		*pagesDiscovered = append(*pagesDiscovered, newPaginationLink)
// 	}
// }

func HTMLHandler(e *colly.HTMLElement, products *[]Product) {
	product := Product{}

	product.Url = e.ChildAttr("a", "href")
	product.Image = e.ChildAttr("img", "src")
	product.Name = e.ChildText("h2")
	product.Price = e.ChildText(".price")

	*products = append(*products, product)

}

func onScrapedHandler(response *colly.Response, products *[]Product) {

	file, err := os.Create("Products.csv")
	if err != nil {
		log.Fatalln("Failed to create CSV file", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	headers := []string{
		"Url",
		"Image",
		"Name",
		"Price",
	}
	writer.Write(headers)

	for _, product := range *products {
		record := []string{
			product.Url,
			product.Image,
			product.Name,
			product.Price,
		}

		writer.Write(record)
	}
	defer writer.Flush()

}
