package builder_test

import (
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
}

func construct(srcFolder string) {
	os.Mkdir(srcFolder, os.ModePerm)

	os.Mkdir(srcFolder+"/aa", os.ModePerm)
	f1, _ := os.Create(srcFolder + "/aa/aaa.ext")
	f1.WriteString("{{ .Prj.Name }}")
	f1.Close()
	f1, _ = os.Create(srcFolder + "/aa/aab.ext.tmpl")
	f1.Close()
	f1, _ = os.Create(srcFolder + "/aa/aac.ext")
	f1.WriteString("test")
	f1.Close()
	f1, _ = os.Create(srcFolder + "/aa/aad.ext")
	f1.WriteString("{{ .Prj.Name }}")
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
	f1.WriteString("{{ .Prj.Name }}")
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
	os.RemoveAll(testFolder + "/passwords.json")
}
