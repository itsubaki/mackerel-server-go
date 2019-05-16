// +build appengine
package infrastructure

import "github.com/itsubaki/mackerel-api/pkg/interfaces/database"

func ShutdownHook(h database.SQLHandler) {
	// noop
}
