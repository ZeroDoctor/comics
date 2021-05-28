package channel

var (
	Shutdown = make(chan bool, 1)

	InTreeChan  = make(chan Data, 4)
	OutTreeChan = make(chan Data)

	InHeaderChan  = make(chan Data, 4)
	OutHeaderChan = make(chan Data)

	InScreenChan  = make(chan Data, 4)
	OutScreenChan = make(chan Data)

	InCommandChan  = make(chan Data, 4)
	OutCommandChan = make(chan Data)
)

type Data struct {
	Type    string
	Integer int
	Boolean bool
	String  string
	Object  interface{}
	Command string
}
