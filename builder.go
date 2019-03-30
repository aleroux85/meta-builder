package builder

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/jinzhu/inflection"
)

type BackRef interface {
	FileStructure() map[string]*FSDirectory
	CmdMatch() map[string]*Exec
	Up() BackRef
}

type DataBranch interface {
	Files() map[string]*FSDirectory
	SetFile(*FSFile)
	File() *FSFile
}

type Entity struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Directories map[string]*FSDirectory `json:"directories"`
	Execs       map[string]*Exec        `json:"execs"`
	Branch      DataBranch              `json:"-"`
	Parent      BackRef                 `json:"-"`
	Error       *error                  `json:"-"`
	changeDetector
}

func (m Entity) FileStructure() map[string]*FSDirectory {
	return m.Directories
}

func (m Entity) CmdMatch() map[string]*Exec {
	return m.Execs
}

func (m Entity) Up() BackRef {
	return m.Parent
}

type Exec struct {
	File string   `json:"file"`
	Exec []string `json:"exec"`
}

type changeDetector struct {
	hash   string `json:"-"`
	change uint8  `json:"-"`
}

func (cd changeDetector) Hash() string {
	return cd.hash
}

func (cd changeDetector) Change(change ...uint8) uint8 {
	if len(change) > 0 {
		cd.change = change[0]
	}
	return cd.change
}

func (cd *changeDetector) CalculateHash(m interface{}) error {
	json, err := json.Marshal(m)
	if err != nil {
		return err
	}

	h := sha1.New()
	_, err = h.Write(json)
	if err != nil {
		return err
	}
	cd.hash = fmt.Sprintf("%x", h.Sum(nil))

	return nil
}

type TemplateMethods struct {
}

func (TemplateMethods) Clean(s string) string {
	reg, _ := regexp.Compile("[^a-zA-Z]+")
	return reg.ReplaceAllString(s, "")
}

func (TemplateMethods) Upper(s string) string {
	return strings.ToUpper(s)
}

func (TemplateMethods) CleanUpper(s string) string {
	reg, _ := regexp.Compile("[^a-zA-Z]+")
	clean := reg.ReplaceAllString(s, "")
	return strings.ToUpper(clean)
}

func (TemplateMethods) Title(s string) string {
	return strings.Title(s)
}

func (TemplateMethods) Plural(s string) string {
	return inflection.Plural(s)
}
