package version

import (
	"testing"
)

func TestGet(t *testing.T)  {
	version = "alpha"
	gitCommit = "commitID"
	gitTreeState = "clean"
	buildDate = "tomorrow"

	v := Get()

	var versionTest = []struct{
		n string
		in string
		out string
	}{
		{
			"testing version",
			v.Version,
			version,
		},
		{
			"testing gitCommit",
			v.GitCommit,
			gitCommit,
		},
		{
			"testing gitTreeState",
			v.GitTreeState,
			gitTreeState,
		},
		{
			"testing buildDate",
			v.BuildDate,
			buildDate,
		},
	}

	for _, test := range versionTest {
		t.Run(test.n, func(t *testing.T) {
			if test.in != test.out {
				t.Errorf("%s != %s", test.in, test.out)
			}
		})
	}
}

func TestString(t *testing.T)  {
	version = "alpha"
	gitCommit = "commitID"
	gitTreeState = "clean"
	buildDate = "tomorrow"

	expected := "Version { version: alpha, gitCommit: commitID, gitTree: clean, buildDate: tomorrow }"

	v := Get()
	if expected != v.String() {
		t.Errorf("%s != %s", expected, v)
	}
}