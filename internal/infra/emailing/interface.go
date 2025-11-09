package emailing

type Emailing interface {
	Send(msg Message) error

	SendAsyc(msg Message, onSendError func(error)) error
}
