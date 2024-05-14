package controllers

import (
	"bytes"
	"io"
	"net/http"

	"github.com/antchfx/htmlquery"
	"github.com/labstack/echo"
)

// company is a structure that contains the company's stock ticker from the client's HTTP request
type Company struct {
	Ticker string `json:"ticker" form:"ticker" query:"ticker"`
}

// GrabPrice - handler method for binding JSON body and scraping for stock price
func GrabPrice(c echo.Context) (err error) {
	// Read the Body content
	var bodyBytes []byte
	if c.Request().Body != nil {

		bodyBytes, _ = io.ReadAll(c.Request().Body)
	}
	// Restore the io.ReadCloser to its original state
	c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	company := new(Company)
	er := c.Bind(company) // bind the structure with the context body
	// on no panic!
	if er != nil {
		panic(er)
	}
	// company ticker
	ticker := company.Ticker
	// yahoo finance base URL
	baseURL := "https://finance.yahoo.com/quote/"
	// price XPath
	pricePath := "//*[@id=\"livePrice svelte-mgkamr\"]/@data-value | //*[@id=\"livePrice svelte-mgkamr\"]/span"
	// load HTML document by binding base url and passed in ticker
	doc, err := htmlquery.LoadURL(baseURL + ticker)
	// uh oh :( freak out!!
	if err != nil {
		panic(err)
	}
	// HTML Node
	context := htmlquery.FindOne(doc, pricePath)
	// from the Node get inner text
	price := string(htmlquery.InnerText(context))
	return c.JSON(http.StatusOK, price)
}
