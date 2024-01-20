package release

import "github.com/KowalskiPiotr98/gotabase"

var (
	getConnector   = func() gotabase.Connector { return gotabase.GetConnection() }
	getTransaction = func() (*gotabase.Transaction, error) { return gotabase.BeginTransaction() }
)
