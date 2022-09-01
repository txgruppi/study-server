package models

import "encoding/gob"

func init() {
	gob.Register(&EventPostCreated{})
	gob.Register(&EventPostTitleUpdated{})
	gob.Register(&EventPostTextUpdated{})
	gob.Register(&EventPostTitleAndTextUpdated{})
}
