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

func getHref(t html.Token) (href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
		}
	}
	return
}

func getAnchors(tokens *html.Tokenizer) map[string]int {
	links := make(map[string]int)
	for {
		tt := tokens.Next()
		switch tt {
		case html.ErrorToken:
			return links
		case html.StartTagToken:
			t := tokens.Token()
			if t.Data == "a" {
				href := getHref(t)
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

func get(URLString string) (map[string]int, error) {
	resp, err := http.Get(URLString)
	if err != nil {
		fmt.Println("\nNo more links...")
		return nil, err
	}
	b := resp.Body
	defer b.Close()
	tokens := html.NewTokenizer(b)
	anchors := getAnchors(tokens)
	resp.Body.Close()
	return anchors, err
}

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
	} else {
		return currentGroupID, groups[fistURISegment], groups
	}
}

func formatLink(link string, URL string) string {
	if !strings.HasPrefix(link, URL) {
		link = URL + link
	}
	if link[len(link)-1:] == "/" {
		link = link[:len(link)-1]
	}
	return link
}

func createGraph(nodes map[string]node, edges []edge) Graph {
	nodeList := make([]node, 0, len(nodes))

	for _, value := range nodes {
		nodeList = append(nodeList, value)
	}

	return Graph{nodeList, edges}
}

func ExportGraphJson(graph Graph) {
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
type Graph struct {
	Nodes []node `json:"nodes"`
	Edges []edge `json:"edges"`
}

func Crawl(URL string, channel chan Graph) {

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

		link, unseenLinks = formatLink(unseenLinks[0], URL), unseenLinks[1:]
		groupCounter, groupID, groups = calculateGroupID(link, groupCounter, groups)

		fmt.Printf("\nCrawling: %s %d", link, len(unseenLinks))

		_, seen := seenLinks[link]
		if seen {
			fmt.Print("  <--- seen")
			continue
		}

		seenLinks[link] = 1
		nodes[link] = node{link, groupID}

		foundLinks, err := get(link)
		if err != nil {
			continue
		}

		fmt.Printf("\nTotal seen links: %d, Total unseen links: %d, Total found links: %d", len(seenLinks), len(unseenLinks), len(foundLinks))

		for foundLink := range foundLinks {
			if !strings.HasPrefix(foundLink, URL) {
				foundLink = URL + foundLink
			}

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

		fmt.Printf("\n Total pages: %d Total links: %d", len(nodes), len(edges))

		nodeList := make([]node, 0, len(nodes))

		for _, value := range nodes {
			nodeList = append(nodeList, value)
		}
		channel <- Graph{nodeList, edges}
	}

	fmt.Printf("\n Total Duration %s", time.Since(then))

	ExportGraphJson(createGraph(nodes, edges))

	fmt.Printf("\n Total pages: %d Total links: %d", len(nodes), len(edges))
}
