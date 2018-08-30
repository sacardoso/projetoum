package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

func setup() string {
	parent, err := ioutil.TempDir("", "example")
	if err != nil {
		log.Fatal(err)
	}
	dir := fmt.Sprintf("%s/a/1", parent)
	os.MkdirAll(dir, 0777)
	dir = fmt.Sprintf("%s/a/x", parent)
	ioutil.WriteFile(dir, []byte("oi"), 0600)

	os.MkdirAll(fmt.Sprintf("%s/b/1/2/3/4", parent), 0777)
	os.MkdirAll(fmt.Sprintf("%s/b/2", parent), 0777)
	return parent
}

func teardown(dir string) {
	os.RemoveAll(dir)
	currentLevels = 0
	maxLevels = -1
}

func TestByteCountSI(t *testing.T) {
	type testStruct struct {
		num      int64
		expected string
	}

	tests := []testStruct{
		testStruct{num: int64(100), expected: "100 B"},
		testStruct{num: int64(2490), expected: "2.5 kB"},
		testStruct{num: int64(7200000), expected: "7.2 MB"},
		testStruct{num: int64(12800000000), expected: "12.8 GB"},
	}

	for _, test := range tests {
		result := ByteCountSI(test.num)
		if result != test.expected {
			t.Errorf("%s != %s", result, test.expected)
		}
	}
}

func TestTree(t *testing.T) {
	parent := setup()
	defer teardown(parent)

	result, err := tree(fmt.Sprintf("%s/a", parent), "")
	if err != nil {
		t.Errorf("err == %v\n", err)
	}
	expected := `a [ 136 B ]
├──1 [ 68 B ]
└──x [ 2 B ]
`
	if result != expected {
		t.Errorf("'%s' != '%s'", result, expected)
	}
}

func TestTreeFourLevels(t *testing.T) {
	parent := setup()
	defer teardown(parent)

	result, err := tree(fmt.Sprintf("%s/b", parent), "")
	if err != nil {
		t.Errorf("err == %v\n", err)
	}
	expected := `b [ 136 B ]
├──1 [ 102 B ]
│  └──2 [ 102 B ]
│     └──3 [ 102 B ]
│        └──4 [ 68 B ]
└──2 [ 68 B ]
`
	if result != expected {
		t.Errorf("'%s' != '%s'", result, expected)
	}
}

func TestTreeMaxLevels(t *testing.T) {
	parent := setup()
	defer teardown(parent)
	maxLevels = 2

	result, err := tree(fmt.Sprintf("%s/b", parent), "")
	if err != nil {
		t.Errorf("err == %v\n", err)
	}
	expected := `b [ 136 B ]
├──1 [ 102 B ]
│  └──2 [ 102 B ]
└──2 [ 68 B ]
`
	if result != expected {
		t.Errorf("'%s' != '%s'", result, expected)
	}
}

func TestTreeMaxLevelsZero(t *testing.T) {
	parent := setup()
	defer teardown(parent)
	maxLevels = 0

	result, err := tree(fmt.Sprintf("%s/b", parent), "")
	if err != nil {
		t.Errorf("err == %v\n", err)
	}
	expected := `b [ 136 B ]
`
	if result != expected {
		t.Errorf("'%s' != '%s'", result, expected)
	}
}

func TestTreeMoreDirs(t *testing.T) {
	parent := setup()
	defer teardown(parent)
	maxLevels = 4 // alterei apara o 1 e 2 aparecerem

	result, err := tree(fmt.Sprintf("%s", parent), "")
	if err != nil {
		t.Errorf("err == %v\n", err)
	}
	expected := fmt.Sprintf(`%s [ 136 B ]
├──a [ 136 B ]
│  ├──1 [ 68 B ]
│  └──x [ 2 B ]
└──b [ 136 B ]
   ├──1 [ 102 B ]
   └──2 [ 68 B ]
`, strings.TrimPrefix(parent, "/var/folders/jc/xh3g03sj3c7_56l1trmyqydr0000gp/T/"))
	if result != expected {
		t.Errorf("'%s' != '%s'", result, expected)
	}
}
