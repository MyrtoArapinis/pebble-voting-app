package voting

import (
	"github.com/giry-dev/pebble-voting-app/pebble-core/anoncred"
	"github.com/giry-dev/pebble-voting-app/pebble-core/util"
)

type Tallier struct {
	election  *Election
	creds     anoncred.CredentialSet
	serialNos util.BytesSet
}
