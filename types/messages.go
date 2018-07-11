package types

import "github.com/golang-collections/go-datastructures/queue"

type Request struct {
	From	string
	Clock   LamportTime
	Body 	map[string]interface{} `json:"body,omitempty"`
}

type Response struct {
	From	string					`json:"from"`
	Success bool                   	`json:"success"`
	Body    map[string]interface{} 	`json:"body,omitempty"`
}

func (req Request) Compare(other queue.Item) int {
	otherReq := other.(Request)

	if req.Clock > otherReq.Clock {
		return 1
	} else if req.Clock == otherReq.Clock {
		return 0
	} else {
		return -1
	}
}

func NewRequest() Request {
	return Request{
		Body:make(map[string]interface{}),
	}
}

func NewResponse() Response {
	return Response{
		Body:make(map[string]interface{}),
	}
}