// Package message
// @author Daud Valentino
package msg

import (
	"fmt"
	"strings"
	"sync"

	"sharefood/pkg/file"
)

var (
	once       sync.Once
	msg        map[string]*Message
	defaultMsg = msgLang{
		text:   "unknown",
		lang:   "en",
		status: 520,
	}
)

type Message struct {
	Name   string `yaml:"name"`
	Status int    `yaml:"status"`
	Langs  []Lang `yaml:"message"`
	lang   map[string]*msgLang
}

// doMap create content map from slice
func (m *Message) doMap() *Message {
	m.lang = make(map[string]*msgLang, 0)
	for _, c := range m.Langs {
		l := strings.ToLower(c.Lang)
		if _, ok := m.lang[l]; !ok {
			m.lang[l] = &msgLang{lang: c.Lang, status: m.Status, text: c.Text}
			continue
		}
	}

	return m
}

type Lang struct {
	Lang string `yaml:"lang"`
	Text string `yaml:"text"`
}

type msgLang struct {
	status int
	lang   string
	text   string
}

func (p msgLang) Text() string {
	return p.text
}

func (p msgLang) Status() int {
	return p.status
}

func (p msgLang) Lang() string {
	return p.lang
}

func Setup(fName string, paths ...string) error {
	var (
		err   error
		langs []Message
	)

	once.Do(func() {
		msg = make(map[string]*Message, 0)
		for _, p := range paths {
			f := fmt.Sprint(p, fName)
			err := file.ReadFromYAML(f, &langs)
			if err != nil {
				continue
			}
			err = nil
		}
	})

	if err != nil {
		return fmt.Errorf("unable to read config from files %v", err)
	}

	for _, m := range langs {
		if _, ok := msg[m.Name]; !ok {
			m := &Message{Name: m.Name, Status: m.Status, Langs: m.Langs}
			msg[m.Name] = m.doMap()
		}
	}

	return err
}

func cleanLangStr(s string) string {
	return strings.ToLower(strings.Trim(s, " "))
}

func Get(key string, lang string) msgLang {
	lang = cleanLangStr(lang)
	if m, ok := msg[key]; ok {
		if c, ok := m.lang[lang]; ok {
			return *c
		}

		return msgLang{status: m.Status}
	}

	return defaultMsg
}

// HaveLang func check language
func HaveLang(key string, lang string) bool {
	if m, ok := msg[key]; ok {
		if _, ok := m.lang[lang]; ok {
			return true
		}
	}

	return false
}
