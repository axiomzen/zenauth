package helpers

import (
	"net/http"

	newrelic "github.com/newrelic/go-agent"
)

type Monitor struct {
	Enabled     bool
	NewRelicApp newrelic.Application
}

func (m *Monitor) NoticeErr(txn newrelic.Transaction, err error) error {
	if m.Enabled && txn != nil {
		txn.NoticeError(err)
	}
	return err
}

func (m *Monitor) StartTransaction(name string, rw http.ResponseWriter, req *http.Request) newrelic.Transaction {
	if m.Enabled {
		return m.NewRelicApp.StartTransaction(name, rw, req)
	}
	return nil
}

func (m *Monitor) EndTransaction(txn newrelic.Transaction) error {
	if m.Enabled && txn != nil {
		return txn.End()
	}
	return nil
}

func (m *Monitor) AddAttribute(txn newrelic.Transaction, name string, value interface{}) error {
	if m.Enabled && txn != nil {
		return txn.AddAttribute(name, value)
	}
	return nil
}
