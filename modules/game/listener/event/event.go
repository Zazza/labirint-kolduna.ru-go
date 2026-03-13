package event

// Event - общий интерфейс для всех событий
type Event interface {
	GetName() string
}
