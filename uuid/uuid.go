package uuid

import (
	"github.com/rs/xid"
)

func New() string {
	guid := xid.New()
	return guid.String()
}

func IDFromString(id string) (string, error) {
	guid, err := xid.FromString(id)
	if err != nil {
		return "", err
	}
	return guid.String(), nil
}

func IsValidID(id string) bool {
	_, err := xid.FromString(id)
	return err == nil
}
