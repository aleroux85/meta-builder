package builder_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	builder "github.com/aleroux85/meta-builder"
)

func TestBuild(t *testing.T) {
	testPath := "testdata/building"
	srcPath := filepath.Join(testPath, "meta")
	metaFilename := "meta.json"
	metaFilePath := filepath.Join(testPath, metaFilename)

	construct(srcPath)
	defer destruct(testPath)

	t.Run("build all files", func(t *testing.T) {
		c := builder.NewConfig(srcPath, testPath)
		p := builder.NewProject()
		c.Load(p, metaFilePath)
		c.BuildAll(false)
		if c.Error() != nil {
			t.Errorf("%+v\n", c.Error())
			return
		}

		// normal file placement using keys for directory names
		if _, err := os.Stat(testPath + "/aa/aaa.ext"); err != nil {
			t.Error(err)
		}
		if _, err := os.Stat(testPath + "/aa/aab.ext"); err != nil {
			t.Error(err)
		}
		if _, err := os.Stat(testPath + "/aa/aaa/aaaa.ext"); err != nil {
			t.Error(err)
		}
		if _, err := os.Stat(testPath + "/aa/aaa/aaaa/aaaaa.ext"); err != nil {
			t.Error(err)
		}
		if _, err := os.Stat(testPath + "/aa/aab/jump/aaba.ext"); err != nil {
			t.Error(err)
		}

		// relative file placement
		if _, err := os.Stat(testPath + "/aba.ext"); err != nil {
			t.Error(err)
		}
		if _, err := os.Stat(testPath + "/abaa.ext"); err != nil {
			t.Error(err)
		}
		if _, err := os.Stat(testPath + "/abaaa.ext"); err != nil {
			t.Error(err)
		}
		if _, err := os.Stat(testPath + "/jump/abba.ext"); err != nil {
			t.Error(err)
		}
		if _, err := os.Stat(testPath + "/jump/here/abbaa.ext"); err != nil {
			t.Error(err)
		}

		// absolute file placement
		if _, err := os.Stat(testPath + "/acaa.ext"); err != nil {
			t.Error(err)
		}

		// files that should not exist
		if _, err := os.Stat(testPath + "/aaaax"); err == nil {
			t.Error("file aaaax should not exist")
		}

		testCases := []struct {
			desc string
			path string
			exp  string
		}{
			{
				desc: "1A Project Name",
				path: "/aa/aaa.ext",
				exp:  "Octagon",
			},
			{
				desc: "2A Plural method",
				path: "/aa/aac.ext",
				exp:  "Chickens",
			},
			{
				desc: "2B Title and Upper methods",
				path: "/aa/aad.ext",
				exp:  "Title UPPER",
			},
			{
				desc: "2C Clean and CleanUpper methods",
				path: "/aa/aae.ext",
				exp:  "jackson JACKSON",
			},
			{
				desc: "3A copy file",
				path: "/aa/aab.ext",
				exp:  "{{ .Prj.Name }}",
			},
			{
				desc: "3B copy all files",
				path: "/ad/ada/adaa.ext",
				exp:  "{{ .Prj.Name }}",
			},
			{
				desc: "4A compile file",
				path: "/all",
				exp:  "data",
			},
			{
				desc: "4B compile file",
				path: "/aaaaa",
				exp:  "data",
			},
			{
				desc: "5A test before force",
				path: "/aa/aac/aaca.ext",
				exp:  "abc",
			},
		}
		for _, tC := range testCases {
			t.Run(tC.desc, func(t *testing.T) {
				got, err := ioutil.ReadFile(testPath + tC.path)
				if err != nil {
					t.Error(err)
					t.FailNow()
				}
				if tC.exp != string(got) {
					t.Errorf(`expected "%s", got "%s"`, tC.exp, got)
				}
			})
		}
	})

	// test forcing rebuild after file changes
	f1, _ := os.OpenFile(srcPath+"/aa/aac/aaca.ext", os.O_WRONLY|os.O_TRUNC, 0666)
	f1.WriteString("def")
	f1.Close()

	t.Run("change a file and rebuild without forcing", func(t *testing.T) {
		c := builder.NewConfig(srcPath, testPath)
		p := builder.NewProject()
		c.Load(p, metaFilePath)
		c.BuildAll(false)
		if c.Error() != nil {
			t.Errorf("%+v\n", c.Error())
			return
		}

		got, err := ioutil.ReadFile(testPath + "/aa/aac/aaca.ext")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		expected := "abc"
		if expected != string(got) {
			t.Errorf(`expected "%s", got "%s"`, expected, got)
		}
	})

	t.Run("change a file and rebuild with forcing", func(t *testing.T) {
		c := builder.NewConfig(srcPath, testPath)
		p := builder.NewProject()
		c.Load(p, metaFilePath)
		c.BuildAll(true)
		if c.Error() != nil {
			t.Errorf("%+v\n", c.Error())
			return
		}

		got, err := ioutil.ReadFile(testPath + "/aa/aac/aaca.ext")
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		expected := "def"
		if expected != string(got) {
			t.Errorf(`expected "%s", got "%s"`, expected, got)
		}
	})
}

func construct(srcFolder string) {
	os.Mkdir(srcFolder, os.ModePerm)

	os.Mkdir(srcFolder+"/aa", os.ModePerm)
	f1, _ := os.Create(srcFolder + "/aa/aaa.ext")
	f1.WriteString("{{ .Prj.Name }}")
	f1.Close()
	f1, _ = os.Create(srcFolder + "/aa/aab.ext.tmpl")
	f1.WriteString("{{ .Prj.Name }}")
	f1.Close()
	f1, _ = os.Create(srcFolder + "/aa/aac.ext")
	f1.WriteString(`{{ .Plural "Chicken" }}`)
	f1.Close()
	f1, _ = os.Create(srcFolder + "/aa/aad.ext")
	f1.WriteString(`{{ .Title "title" }} {{ .Upper "upper" }}`)
	f1.Close()
	f1, _ = os.Create(srcFolder + "/aa/aae.ext")
	f1.WriteString(`{{ .Clean "jack6&son" }} {{ .CleanUpper "jack6&son" }}`)
	f1.Close()
	os.Mkdir(srcFolder+"/aa/aaa", os.ModePerm)
	f1, _ = os.Create(srcFolder + "/aa/aaa/aaaa.ext")
	f1.WriteString("Chicken")
	f1.Close()
	os.Mkdir(srcFolder+"/aa/aaa/aaaa", os.ModePerm)
	f1, _ = os.Create(srcFolder + "/aa/aaa/aaaa/aaaaa.ext")
	f1.Close()
	os.Mkdir(srcFolder+"/aa/aab", os.ModePerm)
	f1, _ = os.Create(srcFolder + "/aa/aab/aaba.ext")
	f1.Close()
	os.Mkdir(srcFolder+"/aa/aac", os.ModePerm)
	f1, _ = os.Create(srcFolder + "/aa/aac/aaca.ext")
	f1.WriteString("abc")
	f1.Close()

	os.Mkdir(srcFolder+"/ab", os.ModePerm)
	f1, _ = os.Create(srcFolder + "/ab/aba.ext")
	f1.Close()
	f1, _ = os.Create(srcFolder + "/ab/abx.ext")
	f1.WriteString("Peanuts")
	f1.Close()
	os.Mkdir(srcFolder+"/ab/aba", os.ModePerm)
	f1, _ = os.Create(srcFolder + "/ab/aba/abaa.ext")
	f1.Close()
	os.Mkdir(srcFolder+"/ab/abb", os.ModePerm)
	f1, _ = os.Create(srcFolder + "/ab/abb/abba.ext")
	f1.Close()
	os.Mkdir(srcFolder+"/ab/aba/abaa", os.ModePerm)
	f1, _ = os.Create(srcFolder + "/ab/aba/abaa/abaaa.ext")
	f1.Close()
	os.Mkdir(srcFolder+"/ab/abb/abba", os.ModePerm)
	f1, _ = os.Create(srcFolder + "/ab/abb/abba/abbaa.ext")
	f1.Close()

	os.Mkdir(srcFolder+"/ac", os.ModePerm)
	os.Mkdir(srcFolder+"/ac/aca", os.ModePerm)
	f1, _ = os.Create(srcFolder + "/ac/aca/acaa.ext")
	f1.Close()

	os.Mkdir(srcFolder+"/ad", os.ModePerm)
	os.Mkdir(srcFolder+"/ad/ada", os.ModePerm)
	f1, _ = os.Create(srcFolder + "/ad/ada/adaa.ext")
	f1.WriteString("{{ .Prj.Name }}")
	f1.Close()
}

func destruct(testFolder string) {
	os.RemoveAll(testFolder + "/meta")
	os.RemoveAll(testFolder + "/aa")
	os.RemoveAll(testFolder + "/ab")
	os.RemoveAll(testFolder + "/aba.ext")
	os.RemoveAll(testFolder + "/abaa.ext")
	os.RemoveAll(testFolder + "/abaaa.ext")
	os.RemoveAll(testFolder + "/abba.ext")
	os.RemoveAll(testFolder + "/jump")
	os.RemoveAll(testFolder + "/acaa.ext")
	os.RemoveAll(testFolder + "/ad")
	os.RemoveAll(testFolder + "/passwords.json")
	os.RemoveAll(testFolder + "/all")
	os.RemoveAll(testFolder + "/aaaaa")
}
