package teleport

// ReturnData builds a response in API handlers.
// OpAndToAndFrom[0] empty: use same operation as peer
// OpAndToAndFrom[1] empty: use peer as receiver
// OpAndToAndFrom[2] empty: use self as sender
func ReturnData(body interface{}, OpAndToAndFrom ...string) *NetData {
	data := &NetData{
		Status: SUCCESS,
		Body:   body,
	}
	if len(OpAndToAndFrom) > 0 {
		data.Operation = OpAndToAndFrom[0]
	}
	if len(OpAndToAndFrom) > 1 {
		data.To = OpAndToAndFrom[1]
	}
	if len(OpAndToAndFrom) > 2 {
		data.From = OpAndToAndFrom[2]
	}
	return data
}

// ReturnError returns an error response; receive should be the original *NetData.
func ReturnError(receive *NetData, status int, msg string, nodeuid ...string) *NetData {
	receive.Status = status
	receive.Body = msg
	receive.From = ""
	if len(nodeuid) > 0 {
		receive.To = nodeuid[0]
	} else {
		receive.To = ""
	}
	return receive
}
