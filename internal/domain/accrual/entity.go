package accrual

type State string

const (
	New        State = "NEW"
	Processing State = "PROCESSING"
	Invalid    State = "INVALID"
	Processed  State = "PROCESSED"
)

type Accrual struct {
	OrderNumber string
	State       State
	Amount      float64
}
