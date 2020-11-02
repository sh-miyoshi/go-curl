package file

import (
	"testing"
)

func TestRemoveChar(t *testing.T) {
	tt := []struct {
		data   []byte
		remove []byte
		expect []byte
	}{
		{
			[]byte("testvalue"),
			[]byte{'t'},
			[]byte("esvalue"),
		},
		{
			[]byte("testvalue"),
			[]byte{'x'},
			[]byte("testvalue"),
		},
	}

	for _, tc := range tt {
		n := removeChar(tc.data, tc.remove)
		if n != len(tc.expect) {
			t.Errorf("removeChar failed, expect length: %d, but got %d", len(tc.expect), n)
		}

		res := tc.data[:n]

		for i, d := range res {
			if d != tc.expect[i] {
				t.Errorf("removeChar failed, expect: %v, but got %v", tc.expect, res)
				return
			}
		}
		return
	}
}
