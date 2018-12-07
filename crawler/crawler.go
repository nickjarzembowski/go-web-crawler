package crawler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// extracts href value from html token
func extractHref(t html.Token) (href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
		}
	}
	return
}

// extracts paths from a tags
func extractAnchors(tokens *html.Tokenizer) map[string]int {
	links := make(map[string]int)
	for {
		tt := tokens.Next()
		switch tt {
		case html.ErrorToken:
			return links
		case html.StartTagToken:
			t := tokens.Token()
			if t.Data == "a" {
				href := extractHref(t)
				u, err := url.Parse(href)
				if err != nil {
					fmt.Println("Error")
				}
				if u.Host == "" || u.Host == "monzo.com" {
					links[u.Path] = 1
				}
			}
		}
	}
}

// fetches html and extracts the a paths from all a tags on the page
func fetchAnchors(URLString string) (map[string]int, error) {
	resp, err := http.Get(URLString)
	if err != nil {
		fmt.Println("\nNo more links...")
		return nil, err
	}
	b := resp.Body
	defer b.Close()
	anchors := extractAnchors(html.NewTokenizer(resp.Body))
	return anchors, err
}

// computes the group id for the given path
// all paths are grouped by the first segement of their uri e.g. /blog/2017 is
// grouped by blog
func calculateGroupID(path string, currentGroupID int, groups map[string]int) (int, int, map[string]int) {
	u, err := url.Parse(path)
	if err != nil {
		panic(err)
	}

	fistURISegment := "/"
	uriSegments := strings.Split(u.Path, "/")
	if len(uriSegments) > 1 {
		fistURISegment = uriSegments[1]
	}

	if _, groupExists := groups[fistURISegment]; !groupExists {
		groups[fistURISegment] = currentGroupID
		return currentGroupID + 1, currentGroupID, groups
	}

	return currentGroupID, groups[fistURISegment], groups
}

// processing on URIs
func FormatLink(link string, URL string) string {
	if strings.HasPrefix(link, "..") {
		link = link[2:]
	}
	if strings.HasPrefix(link, ".") {
		link = link[1:]
	}
	if strings.HasPrefix(link, URL) {
		link = link[len(URL):]
	}

	if link != URL && len(link) > 1 && link[len(link)-1:] == "/" {
		link = link[:len(link)-1]
	}
	if link == "" {
		link = "/"
	}
	if link != URL && link[0] != '/' {
		link = "/" + link
	}

	return strings.TrimSpace(link)
}

// extracts the nodes from the node map and creates a graph struct
func createGraph(nodes map[string]node, edges []edge) Graph {
	nodeList := make([]node, 0, len(nodes))
	for _, value := range nodes {
		nodeList = append(nodeList, value)
	}
	return Graph{nodeList, edges}
}

// ExportGraphJSON persists the provided graph in json format
func ExportGraphJSON(graph Graph) {
	jsonString, err := json.Marshal(graph)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile("sitemap.json", jsonString, 0644)
}

type node struct {
	ID    string `json:"id"`
	Group int    `json:"group"`
}
type edge struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

// Graph contains nodes and edges
type Graph struct {
	Nodes []node `json:"nodes"`
	Edges []edge `json:"edges"`
}

// G Graph built by the crawler
var G Graph

// Crawl crawls the given URL and extracts the site map
// does not crawl links outside of the URL
func Crawl(URL string) {

	then := time.Now()
	link := ""

	seenLinks := make(map[string]int)
	unseenLinks := []string{URL}

	nodes := make(map[string]node)
	edges := []edge{}

	groupCounter := 0
	groupID := 0
	groups := make(map[string]int)

	for {

		if len(unseenLinks) == 0 {
			break
		}

		prior := unseenLinks[0]
		link, unseenLinks = FormatLink(unseenLinks[0], URL), unseenLinks[1:]
		groupCounter, groupID, groups = calculateGroupID(link, groupCounter, groups)

		fmt.Printf("\nCrawling: %s %s %d", link, prior, len(unseenLinks))

		_, seen := seenLinks[link]
		if seen {
			fmt.Print("  <--- seen")
			continue
		}

		seenLinks[link] = 1
		nodes[link] = node{link, groupID}

		foundLinks, err := fetchAnchors(URL + link)
		if err != nil {
			continue
		}

		fmt.Printf("\nTotal seen links: %d, Total unseen links: %d, Total found links: %d", len(seenLinks), len(unseenLinks), len(foundLinks))

		for foundLink := range foundLinks {
			foundLink = FormatLink(foundLink, URL)
			groupCounter, groupID, groups = calculateGroupID(foundLink, groupCounter, groups)

			to := node{foundLink, groupID}
			if _, exists := nodes[to.ID]; !exists {
				nodes[foundLink] = to
			}
			edges = append(edges, edge{link, foundLink})

			_, ok := seenLinks[foundLink]
			if !ok {
				unseenLinks = append(unseenLinks, foundLink)
			}
		}

		G = createGraph(nodes, edges)
	}

	fmt.Printf("\n Total Duration %s", time.Since(then))

	ExportGraphJSON(createGraph(nodes, edges))

	fmt.Printf("\n Total pages: %d Total links: %d", len(nodes), len(edges))
}
