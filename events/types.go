package events
// все что находится в пакете events сможет работать с 
// любым мессенджером
type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(e Event) error
}

type Type int

const (
	// iota - присваивает первой константе значение 0, а далее увеличивает на 1
	Unknown Type = iota
	Message // автоматом будет 1 из-за iota
)

type Event struct {
	Type Type
	Text string
	// в meta таким образом можно положить что угодно
	Meta interface{}
}