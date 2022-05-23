package fetchdata

type TaskAddFetchData struct {
	Auth     string          `json:"auth" bson:""auth`
	Settings FetchedSettings `json:"settings" bson:"settings"`
	Override FetchOverride   `json:"override" bson:"override"`
	Expire   int64           `json:"exp" bson:"exp"`
}
type Interval struct {
	Time int    `json:"time" bson:"time"`
	Text string `json:"text" bson:"text"`
}
type FetchedSettings struct {
	Interval  []Interval `json:"interval" bson:"interval"`
	PageRange [2]int     `json:"pageRange" bson:"pageRange"`
}
type FetchOverride struct {
	Owner int64 `json:"owner" bson:"owner"`
}
