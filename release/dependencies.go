package release

import "github.com/KowalskiPiotr98/gotabase"

var (
	getConnector = func() gotabase.Connector { return gotabase.GetConnection() }
)
