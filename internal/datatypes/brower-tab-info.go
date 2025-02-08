package datatypes

import "fmt"

type BrowserTabInfo struct {
	URL   string
	Title string
}

func (b BrowserTabInfo) String() string {
	return fmt.Sprintf("URL: %s, Title: %s", b.URL, b.Title)
}
