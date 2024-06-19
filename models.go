package main

type Product struct {
	Url   string `json:url`
	Image string `json:image`
	Name  string `json:name`
	Price string `json:price`
}
