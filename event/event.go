package event

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
