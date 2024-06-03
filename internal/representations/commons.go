package representations

import "time"

type ObjectID string
type User string

type Common struct {
	ID           ObjectID  `json:"id"`
	User         User      `json:"user"`
	CreationDate time.Time `json:"creation_date"`
	ModifiedDate time.Time `json:"modified_date"`
}
