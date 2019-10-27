package helper

import "testing"

func TestSortMapByValueEmpty(t *testing.T) {
	data := make(map[string]int)

	actual := SortMapByValue(data)

	if len(actual) > 0 {
		t.Fail()
	}
}

type testCase struct {
	data     map[string]int
	expected PairList
}

var testCases = []testCase{
	{map[string]int{"a": 1, "b": 2, "c": 3}, PairList{{"c", 3}, {"b", 2}, {"a", 1}}},
	{map[string]int{"a": 3, "b": 2, "c": 1}, PairList{{"a", 3}, {"b", 2}, {"c", 1}}},
	{map[string]int{"a": 3, "b": 2, "c": 2}, PairList{{"a", 3}, {"c", 2}, {"b", 2}}},
	{map[string]int{"a": 3, "c": 2, "b": 2}, PairList{{"a", 3}, {"c", 2}, {"b", 2}}},
}

func testArrays(a PairList, b PairList) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i].Key != b[i].Key ||
			a[i].Value != b[i].Value {
			return false
		}
	}
	return true
}

func TestSortMapByValue(t *testing.T) {
	for i, test := range testCases {
		actual := SortMapByValue(test.data)
		if !testArrays(actual, test.expected) {
			t.Errorf("Test #%d failed. Expected: %v, Actual %v", i, test.expected, actual)
		}
	}
}
