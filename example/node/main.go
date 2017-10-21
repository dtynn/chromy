package main

import (
	"context"
	"log"
	"time"

	"github.com/dtynn/chromy"
)

func main() {
	connector := chromy.Connect(
		chromy.ConnectTimeout(5*time.Second),
		chromy.ActionTimeout(1*time.Minute),
		chromy.TaskStepTimeout(10*time.Second),
	)

	t, err := connector.New(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	defer t.Close()

	tasks := chromy.Task{
		chromy.Navigate("https://www.google.com"),
		chromy.DocumentReady(),
		chromy.Click(`div.tsf-p > div.jsb > center > input:last-child`),
		chromy.Sleep(5 * time.Second),
	}

	err = t.Run(context.Background(), tasks)
	if err != nil {
		log.Println(err)
		return
	}
}
