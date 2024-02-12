package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type happyTestCase struct {
	name string
	have string
	want string
}

// Test that valid values are converted correctly.
func TestTimeConvertSuccess(t *testing.T) {
	tests := []happyTestCase{
		{
			"test 1",
			"12:01:00PM",
			"12:01:00",
		},
		{
			"test 2",
			"12:01:00AM",
			"00:01:00",
		},
		{
			"test 3",
			"10:01:00AM",
			"10:01:00",
		},
		{
			"test 4",
			"10:01:00AM",
			"10:01:00",
		},
	}

	for _, tt := range tests {
		func(tt happyTestCase) {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got, err := convertTime(tt.have)
				if err != nil {
					fmt.Println(err)
				}
				assert.Equal(t, tt.want, got)
			})
		}(tt)
	}
}

type errorTestCase struct {
	name string
	have string
}

// Check that our function returns error on incorrect time string conversion.
func TestTimeConvertError(t *testing.T) {
	tests := []errorTestCase{
		{
			"test 1",
			"13:01:00PM",
		},
		{
			"test 2",
			"13:01:00AM",
		},
		{
			"test 3",
			"10:01:00MM",
		},
		{
			"test 4",
			"random string",
		},
		{
			"test 4",
			"",
		},
	}

	for _, tt := range tests {
		func(tt errorTestCase) {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				_, err := convertTime(tt.have)
				assert.Error(t, err)
			})
		}(tt)
	}
}
