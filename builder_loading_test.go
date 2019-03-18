package builder_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	builder "github.com/aleroux85/meta-builder"
	"github.com/aleroux85/meta-builder/refmap"
	"github.com/pkg/errors"
)

func TestLoad(t *testing.T) {
	source := "testdata/loading"
	metaFilename := "meta.json"
	metaFilepath := filepath.Join(source, metaFilename)
	incorrectFilename := "incorrect.json"
	incorrectFilepath := filepath.Join(source, incorrectFilename)

	t.Run("load project with error", func(t *testing.T) {
		project := builder.NewProject()
		*project.Error = fmt.Errorf("pre-existing error")
		project.Load(metaFilepath)
		if *project.Error == nil {
			t.Errorf("expected project error")
		}
		if (*project.Error).Error() != "pre-existing error" {
			t.Errorf(`expected "pre-existing error" error, got "%s"`, (*project.Error).Error())
		}
	})

	t.Run("load project with non-existing file", func(t *testing.T) {
		project := builder.NewProject()
		project.Load("testdata/loading/non-existing.json")
		if *project.Error == nil {
			t.Errorf("expected project error")
		}
		errorString := "open testdata/loading/non-existing.json: no such file or directory"
		if (*project.Error).Error() != errorString {
			t.Errorf(`expected "%s" error, got "%s"`, errorString, (*project.Error).Error())
		}
	})

	t.Run("load project with incorrect file", func(t *testing.T) {
		project := builder.NewProject()
		project.Load(incorrectFilepath)
		if *project.Error == nil {
			t.Errorf("expected project error")
		}
		errorString := "unexpected end of JSON input"
		if (*project.Error).Error() != errorString {
			t.Errorf(`expected "%s" error, got "%s"`, errorString, (*project.Error).Error())
		}
	})

	t.Run("load project", func(t *testing.T) {
		project := builder.NewProject()
		project.Load(metaFilepath)
		if *project.Error != nil {
			t.Errorf("got project error")
		}
		name := "a"
		if project.Name != name {
			t.Errorf(`expected "%s", got "%s"`, name, project.Name)
		}
	})
}

func TestLoadProcess(t *testing.T) {
	var err error
	source := "testdata/loading"
	metaFilename := "meta.json"
	metaFilepath := filepath.Join(source, metaFilename)

	project := builder.NewProject(&err)

	project.Load(metaFilepath)
	if err != nil {
		t.Error(err)
	}

	refMap := refmap.NewRefMap()
	refMap.Start()
	refMap.Set("location", source)

	project.Process(refMap)
	if err != nil {
		err = errors.Wrap(err, "processing configuration file")
		return
	}

	testCases := []struct {
		desc, exp, got string
	}{
		{
			desc: "1A Project Name",
			exp:  project.Name,
			got:  "a",
		},
		{
			desc: "2A Project Directory aa Name",
			exp:  project.Directories["aa"].Name,
			got:  "aa",
		},
		{
			desc: "2B Project Directory aa File aaa.ext Name",
			exp:  project.Directories["aa"].Files["aaa.ext"].Name,
			got:  "aaa.ext",
		},
		{
			desc: "2C Project Directory aa Directory aaa Name",
			exp:  project.Directories["aa"].Directories["aaa"].Name,
			got:  "aaa",
		},
		{
			desc: "2D Project Directory aa Directory aaa File aaaa.ext Name",
			exp:  project.Directories["aa"].Directories["aaa"].Files["aaaa.ext"].Name,
			got:  "aaaa.ext",
		},
		{
			desc: "2E Project Directory aa Directory aaa Directory aaaa Name",
			exp:  project.Directories["aa"].Directories["aaa"].Directories["aaaa"].Name,
			got:  "aaaa",
		},
		{
			desc: "2D Project Directory aa Directory aaa Directory aaaa File aaaaa.ext Name",
			exp:  project.Directories["aa"].Directories["aaa"].Directories["aaaa"].Files["aaaaa.ext"].Name,
			got:  "aaaaa.ext",
		},
		{
			desc: "3A Project Directory aa Directory aab",
			exp:  project.Directories["aa"].Directories["aab"].Name,
			got:  "aab",
		},
		{
			desc: "3AA Project Directory aa Directory aab File aaba.ext Name",
			exp:  project.Directories["aa"].Directories["aab"].Files["aaba.ext"].Name,
			got:  "aaba.ext",
		},
		{
			desc: "3AB Project Directory aa Directory aab DestinationPath",
			exp:  project.Directories["aa"].Directories["aab"].DestinationPath,
			got:  "aa/aab/jump",
		},
		{
			desc: "4A Project Directory ab Name",
			exp:  project.Directories["ab"].Name,
			got:  "ab",
		},
		{
			desc: "4AA Project Directory ab File aba.ext Name",
			exp:  project.Directories["ab"].Files["aba.ext"].Name,
			got:  "aba.ext",
		},
		{
			desc: "4AB Project Directory ab DestinationPath",
			exp:  project.Directories["ab"].DestinationPath,
			got:  ".",
		},
		{
			desc: "4B Project Directory ab Directory aba Name",
			exp:  project.Directories["ab"].Directories["aba"].Name,
			got:  "aba",
		},
		{
			desc: "4BA Project Directory ab Directory aba File abaa.ext Name",
			exp:  project.Directories["ab"].Directories["aba"].Files["abaa.ext"].Name,
			got:  "abaa.ext",
		},
		{
			desc: "4BB Project Directory ab Directory aba DestinationPath",
			exp:  project.Directories["ab"].Directories["aba"].DestinationPath,
			got:  ".",
		},
		{
			desc: "4C Project Directory ab Directory aba Directory abaa Name",
			exp:  project.Directories["ab"].Directories["aba"].Directories["abaa"].Name,
			got:  "abaa",
		},
		{
			desc: "4CA Project Directory ab Directory aba Directory abaa File abaaa.ext Name",
			exp:  project.Directories["ab"].Directories["aba"].Directories["abaa"].Files["abaaa.ext"].Name,
			got:  "abaaa.ext",
		},
		{
			desc: "4CB Project Directory ab Directory aba Directory abaa DestinationPath",
			exp:  project.Directories["ab"].Directories["aba"].Directories["abaa"].DestinationPath,
			got:  ".",
		},
		{
			desc: "5A Project Directory ab Directory abb Name",
			exp:  project.Directories["ab"].Directories["abb"].Name,
			got:  "abb",
		},
		{
			desc: "5AA Project Directory ab Directory abb File abba.ext Name",
			exp:  project.Directories["ab"].Directories["abb"].Files["abba.ext"].Name,
			got:  "abba.ext",
		},
		{
			desc: "5AB Project Directory ab Directory abb DestinationPath",
			exp:  project.Directories["ab"].Directories["abb"].DestinationPath,
			got:  "jump",
		},
		{
			desc: "5B Project Directory ab Directory abb Directory abba Name",
			exp:  project.Directories["ab"].Directories["abb"].Directories["abba"].Name,
			got:  "abba",
		},
		{
			desc: "5BA Project Directory ab Directory abb Directory abba File abbaa.ext Name",
			exp:  project.Directories["ab"].Directories["abb"].Directories["abba"].Files["abbaa.ext"].Name,
			got:  "abbaa.ext",
		},
		{
			desc: "5BB Project Directory ab Directory abb Directory abba DestinationPath",
			exp:  project.Directories["ab"].Directories["abb"].Directories["abba"].DestinationPath,
			got:  "jump/here",
		},
		{
			desc: "6A Project Directory ac Directory aca Name",
			exp:  project.Directories["ac"].Directories["aca"].Name,
			got:  "aca",
		},
		{
			desc: "6AA Project Directory ac Directory aca File acaa.ext Name",
			exp:  project.Directories["ac"].Directories["aca"].Files["acaa.ext"].Name,
			got:  "acaa.ext",
		},
		{
			desc: "6AB Project Directory ac Directory aca DestinationPath",
			exp:  project.Directories["ac"].Directories["aca"].DestinationPath,
			got:  "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if tC.exp != tC.got {
				t.Errorf(`expected "%s", got "%s"`, tC.exp, tC.got)
			}
		})
	}
}

func TestLoadPasswords(t *testing.T) {
	var err error
	source := "testdata/loading"
	metaFilename := "passwords.json"
	metaFilepath := filepath.Join(source, metaFilename)

	p := builder.NewProject(&err)
	p.LoadSecrets(metaFilepath)
	if err != nil {
		t.Errorf("%+v", err)
		return
	}

	if len(p.Secrets) == 0 {
		t.Error("passwords are missing")
	}

	testPassword := p.Secrets[0]
	p.Secrets = []string{}
	p.LoadSecrets(metaFilepath)
	if len(p.Secrets) == 0 {
		t.Error("passwords are missing")
	}
	if testPassword != p.Secrets[0] {
		t.Error("first password not matching")
	}

	err = os.Remove(metaFilepath)
	if err != nil {
		panic(err)
	}
}
