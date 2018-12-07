### Go Site Map Generator

Go implementation of a web crawler that scrapes pages and links from a given site with D3 visualisation of the resulting node graph.

## How to

Start the application by running `go run main.go`

If on Windows give access to the network

Start the crawler by performing a GET request to `localhost:3000/crawl`

Get the status of the crawler by performing a GET request to `localhost:3000/status`

Get the data scraped by the crawler by performing a GET request to `localhost:3000/data`

Export the data scraped by the crawler to a JSON file by performing a GET request to `localhost:3000/export`

View the data by visiting the root `localhost:3000/`. Make sure the crawler is running before you do this - data is stored in memory and lost when the server is shut down. Refresh the page to see data scraped since last loaded. Hover over a node to see the uri.

## Notes

The crawler took approximately 7m28s to crawl the entire site