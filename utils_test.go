package main

import (
	"testing"
)

func TestStingsEqual(t *testing.T) {
	// test 1
	slice1 := []string{"this", "is", "a", "test."}
	slice2 := []string{"this", "a", "is", "test."}
	if StringsEqual(slice1, slice2) {
		t.Errorf("Received True when False was expected for slices with differently ordered elements")
	}
	// test 2
	slice1 = []string{"this", "is", "a", "test."}
	slice2 = []string{"this", "is", "a", "test."}
	if !StringsEqual(slice1, slice2) {
		t.Errorf("Received False when True was expected for slices with same elements in the same oder")
	}
	// test 3
	slice1 = []string{"this", "is", "a", "test."}
	slice2 = []string{"this", "is", "a", "test.", ""}
	if StringsEqual(slice1, slice2) {
		t.Errorf("Recevied True when False was expected for slices with different amount of elements")
	}
}

func TestStringsIndexOf(t *testing.T) {
	// test 1
	slice := []string{"this", "is", "a", "test."}
	index := StringsIndexOf(slice, "is") // expect 2
	if index != 1 {
		t.Errorf("Index Received: %d did not match expected: %d", index, 1)
	}

	// test 2
	slice = []string{"A", "more", "difficult", "", "test", "has", " whitespaces ", " ", ""}
	index1 := StringsIndexOf(slice, " ")
	if index1 != 7 {
		t.Errorf("Index Received: %d did not match expected: %d", index1, 7)
	}
	index2 := StringsIndexOf(slice, "")
	if index2 != 3 {
		t.Errorf("Index Received: %d did not match expected: %d", index1, 3)
	}
}

func TestStringsRemoveElements(t *testing.T) {
	// test 1
	slice := []string{"this", "is", "a", "test."}
	expectedSlice := []string{"this", "is", "test."}
	resSlice := StringsRemoveElements(slice, "a")
	if !StringsEqual(expectedSlice, resSlice) {
		t.Errorf("Expected did not match received with first slice")
	}

	// test 2
	slice = []string{"A", "more", "difficult", "", "test", "has", " whitespaces ", " ", ""}
	expectedSlice = []string{"A", "more", "difficult", "test", "has", " whitespaces ", " "}
	resSlice = StringsRemoveElements(slice, "")
	if !StringsEqual(expectedSlice, resSlice) {
		t.Errorf("Expected did not match received with second slice")
	}

	// test 3
	slice = []string{"A", "more", "difficult", "", "test", "has", " whitespaces ", " ", ""}
	expectedSlice = []string{"A", "difficult", "", "test", " whitespaces ", " ", ""}
	resSlice = StringsRemoveElements(slice, "has", "more")
	if !StringsEqual(expectedSlice, resSlice) {
		t.Errorf("Expected did not match received with third slice")
	}

	// test 4
	slice = []string{"system.js", "system_events"}
	filterSlice := []string{"system.js", "jobs.queue", "admin_logs", "partial_subscriptions", "recipes_reviews_requests",
		"stats_subscriptions"}
	expectedSlice = []string{"system_events"}
	resSlice = StringsRemoveElements(slice, filterSlice...)
	if !StringsEqual(expectedSlice, resSlice) {
		t.Errorf("Expected did not match received with fourth slice")
	}

}
