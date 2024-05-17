package data

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rs/xid"
)

// based on xid
type Emoji struct {
	Name      string
	Emoji     string
	CreatedAt time.Time
	PID       uint16
	Machine   string
}

func (e *Emoji) Encode() string {
	return fmt.Sprintf("%s;%s;%s", xid.New().String(), e.Emoji, e.Name)
}

func (e *Emoji) Decode(encodedString string) error {
	split := strings.Split(encodedString, ";")
	if len(split) > 3 {
		return errors.New("incorrect encoding")
	}

	id, err := xid.FromString(split[0])
	if err != nil {
		return err
	}
	emoji := split[1]
	name := split[2]

	e.Name = name
	e.Emoji = emoji
	e.CreatedAt = id.Time()
	e.Machine = string(id.Machine())
	e.PID = id.Pid()

	return nil
}

type Transaction struct {
	ID        int64
	CreatedAt time.Time
	Nominal   float64
	Detail    string
	WTSID     int8
	EmojiID   Emoji
	Version   int
	ProjectID int64
	CreatedBy int64
}
