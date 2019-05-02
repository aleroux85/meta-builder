package builder_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	builder "github.com/aleroux85/meta-builder"
)

func TestFileWatch(t *testing.T) {
	testPath := "testdata/watching"
	srcPath := filepath.Join(testPath, "meta")
	metaFilename := "meta.json"
	metaFilePath := filepath.Join(testPath, metaFilename)

	t.Run("Test multiple files", func(t *testing.T) {
		construct(srcPath)
		defer destruct(testPath)

		c := builder.NewConfig(srcPath, testPath)
		p := builder.NewProject()
		c.Load(p, metaFilePath)
		c.BuildAll(false)
		c.Watch(10 * time.Millisecond)
		if c.Error() != nil {
			t.Errorf("%+v\n", c.Error())
			return
		}
		c.Finish()

		fmt.Println("start")

		testCases := []struct {
			desc             string
			changeFile       string
			changedFile      string
			preChangeContent string
			changeContent    string
			changedContent   string
		}{
			{
				desc:             "1A Project Directory aa File aaa.ext",
				changeFile:       "/aa/aaa.ext",
				changedFile:      "/aa/aaa.ext",
				preChangeContent: "Octagon plus",
				changeContent:    "new stuff",
				changedContent:   "new stuff plus",
			},
		}
		for _, tC := range testCases {
			t.Run(tC.desc, func(t *testing.T) {
				got, err := ioutil.ReadFile(testPath + tC.changedFile)
				if err != nil {
					t.Error(err)
					t.FailNow()
				}
				if tC.preChangeContent != string(got) {
					t.Errorf(`expected "%s", got "%s"`, tC.preChangeContent, got)
				}

				f1, err := os.OpenFile(srcPath+tC.changeFile, os.O_WRONLY|os.O_TRUNC, 0666)
				if err != nil {
					t.Error(err)
					t.FailNow()
				}
				f1.WriteString(tC.changeContent)
				f1.Close()

				time.Sleep(20 * time.Millisecond)

				got, err = ioutil.ReadFile(testPath + tC.changedFile)
				if err != nil {
					t.Error(err)
					t.FailNow()
				}
				if tC.changedContent != string(got) {
					t.Errorf(`expected "%s", got "%s"`, tC.changedContent, got)
				}
			})
		}

		c.StopWatching()
	})
}
