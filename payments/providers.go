package payments

type Providers interface {
	Send()
	Receive()
}
