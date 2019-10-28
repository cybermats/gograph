package searcher

import "testing"

func TestGSPath(t *testing.T) {
	type testCase struct {
		path         string
		expectedPath *GSPath
		expectedErr  bool
	}

	var testCases = []testCase{
		{"gs://bucket/object", &GSPath{"bucket", "object"}, false},
		{"gs://bucket/object.txt", &GSPath{"bucket", "object.txt"}, false},
		{"gs://bucket/dir/object.txt", &GSPath{"bucket", "dir/object.txt"}, false},
		{"gs:/bucket/dir/object.txt", nil, true},
		{"gs//bucket/object", nil, true},
		{"gs:///bucket/object", nil, true},
		{"g://bucket/object", nil, true},
		{"//bucket/object", nil, true},
		{"gs://bucket/", nil, true},
		{"gs://", nil, true},
		{"gs:///", nil, true},
		{"foo", nil, true},
	}

	for i, c := range testCases {
		gspath, err := ParseGCS(c.path)
		if c.expectedPath != nil {
			if gspath == nil ||
				gspath.Bucket != c.expectedPath.Bucket ||
				gspath.Object != c.expectedPath.Object {
				t.Errorf("#%d failed. Expected: %v, actual: %v",
					i, c.expectedPath, gspath)
			}
		} else {
			if gspath != nil {
				t.Errorf("#%d failed. Expected: %v, actual: %v",
					i, c.expectedPath, gspath)
			}
		}
		if ((err != nil) && !c.expectedErr) ||
			((err == nil) && c.expectedErr) {
			t.Errorf("#%d failed. Expected: %v, actual: %v",
				i, c.expectedErr, err)
		}

	}

}
