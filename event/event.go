package event

import (
	"fmt"
	"time"
)

type EventTimes struct {
	Name      string
	StartTime string
	EndTime   string
}

type ResourceType struct {
	Id       string
	Kind     string
	Name     string
	PathName string
}

type TagsType struct {
	Color string
	Name  string
	Id    string
}

type EventType struct {
	InstanceId string
	StartTime  string
	EndTime    string
	Name       string
	Location   string
	Times      []EventTimes
	Resources  []ResourceType
	Tags       []TagsType
}

func FetchEvents() {
	year, month, day := time.Now().Add(-72 * time.Hour).Date()
	fmt.Println(year, month, day)
	// resp, err := http.Get("https://jsonplaceholder.typicode.com/posts/1")
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Println(resp)
}
