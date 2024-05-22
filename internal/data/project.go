package data

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hisamcode/acis/internal/validator"
	"github.com/rs/xid"
)

// based on xid
type Emoji struct {
	ID        string
	Name      string
	Emoji     string
	CreatedAt time.Time
	PID       uint16
	Machine   string
}

// encode Emoji to id;emoji;name pattern
func (e *Emoji) Encode() string {
	return fmt.Sprintf("%s;%s;%s", e.ID, e.Emoji, e.Name)
}

func (e *Emoji) Decode(encodedString string) error {
	// TODO: error index out of range
	split := strings.Split(encodedString, ";")
	if len(split) > 3 {
		return errors.New("incorrect encoding")
	}

	if split[0] == "empty;empty;empty" {
		return nil
	}

	id, err := xid.FromString(split[0])
	if err != nil {
		return err
	}
	emoji := split[1]
	name := split[2]

	e.ID = split[0]
	e.Name = name
	e.Emoji = emoji
	e.CreatedAt = id.Time()
	e.Machine = string(id.Machine())
	e.PID = id.Pid()

	return nil
}

func CreateEmoji(name, emoji string) Emoji {
	emojiEncoded := fmt.Sprintf("%s;%s;%s", xid.New().String(), emoji, name)
	e := Emoji{}
	e.Decode(emojiEncoded)

	return e

}

func DecodeEmojis(encodedEmojis []string) []Emoji {
	emojis := []Emoji{}
	for _, v := range encodedEmojis {
		emoji := Emoji{}
		emoji.Decode(v)
		emojis = append(emojis, emoji)
	}

	return emojis
}

type Project struct {
	ID        int64
	CreatedAt time.Time
	Name      string
	Detail    string
	WTS       []string
	Emojis    []Emoji
	Version   int
	UserID    int64
}

func (p *Project) FindEmoji(ID string) (Emoji, error) {
	for _, v := range p.Emojis {
		if v.ID == ID {
			return v, nil
		}
	}
	return Emoji{}, errors.New("emoji not found")
}

func (p *Project) UpdateEmoji(emoji Emoji) error {
	for k, v := range p.Emojis {
		if v.ID == emoji.ID {
			p.Emojis[k] = emoji
			return nil
		}
	}
	return errors.New("emoji not found")
}

func (p *Project) FindIndexEmojiByID(id string) (int, error) {
	for k, v := range p.Emojis {
		if v.ID == id {
			return k, nil
		}
	}
	return -1, errors.New("emoji index not found")
}

func (p *Project) DeleteEmoji(emoji Emoji) error {
	index, err := p.FindIndexEmojiByID(emoji.ID)
	if err != nil {
		return err
	}

	p.Emojis = append(p.Emojis[:index], p.Emojis[index+1:]...)
	return nil
}

func ValidateEmoji(v *validator.Validator, emoji *Emoji) {
	v.Check(emoji.Name != "", "emoji_name", "name must be provided")
	v.Check(emoji.Emoji != "", "emoji", "emoji must be provided")
	v.Check(len(emoji.Name) <= 500, "emoji_name", "name must not be more than 500 bytes/character long")
	v.Check(len(emoji.Emoji) <= 500, "emoji", "emoji must not be more than 500 bytes/character long")
}

func ValidateProject(v *validator.Validator, project *Project) {
	v.Check(project.Name != "", "name", "Name must be provided")
	v.Check(len(project.Name) <= 500, "name", "Name must not be more than 500 bytes/character long")
	v.Check(len(project.Detail) <= 2000, "detail", "Detail must not be more than 2000 bytes/character long")
}
