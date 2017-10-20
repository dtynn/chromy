package main

import (
	"context"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/dtynn/chromy"
	"github.com/dtynn/chromy/cdptype"
)

var commonJSpattern = regexp.MustCompile("http[s]?://.+/ReactJSstatic/js/entrys/common\\..+\\.bundle\\.js")

func main() {
	connector := chromy.Connect(
		chromy.ConnectTimeout(5*time.Second),
		chromy.ActionTimeout(1*time.Minute),
		chromy.TaskStepTimeout(5*time.Second),
	)

	t, err := connector.New(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	defer t.Close()

	onCategoryNodes := func(nodes ...*chromy.Node) error {
		for i, node := range nodes {
			var category, href string

			href = node.Attr("href")
			for _, child := range node.Children {
				if child.NodeType == cdptype.NodeTypeText {
					category = child.NodeValue
					break
				}
			}

			log.Printf("%d [%s]<%s>", i+1, category, href)
		}

		return nil
	}

	task := chromy.Task{
		chromy.Navigate("https://ezbuy.sg"),
		chromy.DocumentReady(),
		chromy.WaitResource(http.MethodGet, commonJSpattern),
		chromy.OnNodeAll(`div[id^="category-"] > div:first-child > a`, onCategoryNodes),
	}

	if err := t.Run(context.Background(), task); err != nil {
		log.Println("task error:", err)
	}
}
