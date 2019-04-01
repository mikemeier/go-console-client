package message

import (
	"reflect"
	"testing"
)

func TestParseCommand(t *testing.T) {
	arguments := parseCommand("echo gugusö")
	expected := []string{"echo", "gugusö"}
	isEqual(t, arguments, expected)

	arguments = parseCommand(" ")
	expected = nil
	isEqual(t, arguments, expected)

	arguments = parseCommand("echo \"gugusö hallo\"")
	expected = []string{"echo", "gugusö hallo"}
	isEqual(t, arguments, expected)

	arguments = parseCommand("echo \"gugusö hallo\" lala lölö")
	expected = []string{"echo", "gugusö hallo", "lala", "lölö"}
	isEqual(t, arguments, expected)

	arguments = parseCommand("echo \"gugusö hallo")
	expected = []string{"echo", "gugusö hallo"}
	isEqual(t, arguments, expected)
}

func isEqual(t *testing.T, result []string, expected []string) {
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("%#v != %#v", expected, result)
	}
}
