package reflect

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIterateArray(t *testing.T) {
	testCases := []struct {
		name       string
		entity     any
		wantValues []any
		wantErr    error
	}{
		{
			name:       "ints",
			entity:     [3]int{1, 2, 3},
			wantValues: []any{1, 2, 3},
		},
		{
			name:       "slice",
			entity:     []int{1, 2, 3},
			wantValues: []any{1, 2, 3},
		},
	}

	for _, tc := range testCases {
		array, err := IterateArray(tc.entity)
		assert.Equal(t, tc.wantErr, err)
		if err != nil {
			return
		}
		assert.Equal(t, tc.wantValues, array)
	}
}

/*func TestIterateMap(t *testing.T) {
	testCases := []struct {
		name       string
		entity     any
		wantKeys   []any
		wantValues []any
		wantErr    error
	}{
		{
			name: "ints",
			entity: map[string]string{
				"1": "a",
				"2": "b",
			},
			wantKeys:   []any{"1", "2"},
			wantValues: []any{"a"},
		},
	}

	for _, tc := range testCases {
		keys, values, err := IterateMap(tc.entity)
		assert.Equal(t, tc.wantErr, err)
		if err != nil {
			return
		}
		assert.Equal(t, len(tc.wantKeys), len(keys))
		assert.Equal(t, len(tc.wantValues), len(values))
		assert.EqualValues(t, tc.wantKeys, keys)
		assert.EqualValues(t, tc.wantValues, values)
	}
}
*/
