package order

type State string

const (
	New        State = "NEW"
	Processing State = "PROCESSING"
	Invalid    State = "INVALID"
	Processed  State = "PROCESSED"
)
