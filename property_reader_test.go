package helpers

import (
	"os"
	"testing"
	//"io"
)

func createSampleFile() {
	f, _ := os.Create("test.properties")
	f.WriteString("key1=value1\n")
	f.WriteString("key2=value2\n")
	f.WriteString("\n")
	f.WriteString("key3=value3\n")
	f.WriteString("key4=\n")
	f.WriteString("key5\n")
	f.WriteString("key6=value6\n")
	f.Close()
}

func deleteSampleFile() {
	os.Remove("test.properties")
}

func TestMissingFile(t *testing.T) {

	success, err := LoadPropertyFile("config", "does.not.exist")
	if success == true || err == nil {
		t.Errorf("Expected ReadFile to fail.")
		t.Fail()
	}
}

func TestSimpleOpen(t *testing.T) {

	createSampleFile()

	success, err := LoadPropertyFile("config", "test.properties")
	if success == false || err != nil {
		t.Errorf("Encountered error: %s.", err)
		t.Fail()
	}

	deleteSampleFile()
}

func TestSimpleProperties(t *testing.T) {

	createSampleFile()

	success, err := LoadPropertyFile("config", "test.properties")
	if success == false || err != nil {
		t.Errorf("Encountered error: %s.", err)
		t.Fail()
	}

	v, f := Properties["config"]["key1"]
	if v != "value1" || f == false {
		t.Errorf("Did not find 'key1'.")
		t.Fail()
	}
	v, f = Properties["config"]["key2"]
	if v != "value2" || f == false {
		t.Errorf("Did not find 'key2'.")
		t.Fail()
	}
	v, f = Properties["config"]["key3"]
	if v != "value3" || f == false {
		t.Errorf("Did not find 'key3'.")
		t.Fail()
	}
	v, f = Properties["config"]["key4"]
	if v != "" || f == true {
		t.Errorf("Found 'key4' when it should be missing.")
		t.Fail()
	}
	v, f = Properties["config"]["key5"]
	if v != "" || f == true {
		t.Errorf("Found 'key5' when it should be missing.")
		t.Fail()
	}
	v, f = Properties["config"]["key6"]
	if v != "value6" || f == false {
		t.Errorf("Did not find 'key6'.")
		t.Fail()
	}
	v, f = Properties["config"]["key7"]
	if v != "" || f == true {
		t.Errorf("Found 'key7' when it should be missing.")
		t.Fail()
	}
	v, f = Properties["config"][""]
	if v != "" || f == true {
		t.Errorf("Found '' when it should be missing.")
		t.Fail()
	}

	deleteSampleFile()
}
