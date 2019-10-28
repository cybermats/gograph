package searcher

import (
	"regexp"
)

// GSPath contains full path to an object in Google Cloud Storage
type GSPath struct {
	Bucket string
	Object string
}

type errorConst string

// PathError is returned when a path is not a valid Google Cloud Storage path
const PathError errorConst = "Not a valid gs path."

func (e errorConst) Error() string {
	return string(e)
}

var re = regexp.MustCompile(`^gs://([^/]+)/(.+)$`)

// ParseGCS parses a Google Cloud Storage path on the form of
// gs://<bucket>/<object> and returns a GSPath with this parsed info.
// If the parse doesn't match it will return a PathError.
func ParseGCS(path string) (*GSPath, error) {
	m := re.FindStringSubmatch(path)
	if m == nil {
		return nil, PathError
	}
	return &GSPath{m[1], m[2]}, nil
}
