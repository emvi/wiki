package util

import (
	"github.com/emvi/hide"
	"github.com/speps/go-hashids"
)

const (
	hashidSalt      = "emvi"
	hashidMinLength = 20
)

// GetPictureFilenameFromId generates a unique filename from ID.
func GetPictureFilenameFromId(id hide.ID) (string, error) {
	hashdata := hashids.NewData()
	hashdata.Salt = hashidSalt
	hashdata.MinLength = hashidMinLength
	hashid, err := hashids.NewWithData(hashdata)

	if err != nil {
		return "", err
	}

	hash, err := hashid.EncodeInt64([]int64{int64(id)})

	if err != nil {
		return "", err
	}

	return hash, nil
}
