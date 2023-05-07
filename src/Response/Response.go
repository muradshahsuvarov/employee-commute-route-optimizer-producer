package Response

type Response struct {
	Id        string
	Topic     string
	Partition int32
	Offset    int64
	Error     bool
}
