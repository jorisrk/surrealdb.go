package models

import (
	"fmt"
	"time"

	"github.com/surrealdb/surrealdb.go/pkg/constants"

	"github.com/fxamacker/cbor/v2"
	"github.com/gofrs/uuid"
)

type TableOrRecord interface {
	string | Table | RecordID | []Table | []RecordID
}

type Table string

// type UUID string

// type UUIDBin []byte
type UUID struct {
	uuid.UUID
}

type Decimal string

type CustomDateTime time.Time

func (d *CustomDateTime) MarshalCBOR() ([]byte, error) {
	enc := getCborEncoder()

	t := time.Time(*d)
	totalNS := t.UnixNano()

	s := totalNS / constants.OneSecondToNanoSecond
	ns := totalNS % constants.OneSecondToNanoSecond

	return enc.Marshal(cbor.Tag{
		Number:  uint64(DateTimeCompactString),
		Content: [2]int64{int64(s), int64(ns)},
	})
}

func (d *CustomDateTime) UnmarshalCBOR(data []byte) error {
	dec := getCborDecoder()

	var temp [2]interface{}
	err := dec.Unmarshal(data, &temp)
	if err != nil {
		return err
	}

	s, ok1 := temp[0].(uint64)
	ns, ok2 := temp[1].(uint64)

	if !ok1 || !ok2 {
		return fmt.Errorf("expected uint64 for seconds and nanoseconds, got %T and %T", temp[0], temp[1])
	}

	*d = CustomDateTime(time.Unix(int64(s), int64(ns)))

	return nil
}

type CustomNil struct {
}

func (c *CustomNil) MarshalCBOR() ([]byte, error) {
	enc := getCborEncoder()

	return enc.Marshal(cbor.Tag{
		Number:  uint64(NoneTag),
		Content: nil,
	})
}

func (c *CustomNil) UnMarshalCBOR(data []byte) error {
	*c = CustomNil{}
	return nil
}

var None = CustomNil{}
