package wallet

import (
	"github.com/google/uuid"
	"time"
)

type Wallet struct {
	ID         string
	CustomerID string
	Balance    float64
	Withdrawn  float64
	CreatedAt  time.Time
	events     []Event
	version    int
}

func NewWallet(customerID string) *Wallet {
	wallet := &Wallet{
		ID:         uuid.NewString(),
		CustomerID: customerID,
		Balance:    0,
		Withdrawn:  0,
		CreatedAt:  time.Now(),
		events:     make([]Event, 0),
	}

	return wallet
}

func (w *Wallet) HandleCommand(cmd Command) error {
	switch c := cmd.(type) {
	case *CreateCommand:
		w.addEvent(&Created{CustomerID: w.CustomerID, Timestamp: time.Now()})
		//w.addEvent(&Deposited{CustomerID: w.CustomerID, Amount: 100, Timestamp: time.Now()})
		return nil
	case *DepositCommand:
		w.addEvent(&Deposited{CustomerID: w.CustomerID, Amount: c.Amount, Timestamp: time.Now()})
		return nil
	case *WithdrawCommand:
		if w.Balance < c.Amount {
			return ErrInsufficientFunds
		}

		w.addEvent(&Withdrawn{CustomerID: w.CustomerID, Amount: c.Amount, Timestamp: time.Now(), OrderNumber: string(c.OrderNumber)})
		return nil
	}
	return nil
}

func (w *Wallet) ApplyEvent(event Event) {
	switch e := event.(type) {
	case *Deposited:
		w.Balance += e.Amount
	case *Withdrawn:
		w.Balance -= e.Amount
		w.Withdrawn += e.Amount
	}
	w.version++
}

func (w *Wallet) addEvent(event Event) {
	w.events = append(w.events, event)
	w.ApplyEvent(event)
}

func (w *Wallet) Reset() {
	w.events = make([]Event, 0)
}

func (w *Wallet) Events() []Event {
	return w.events
}

func (w *Wallet) Version() int {
	return w.version
}
