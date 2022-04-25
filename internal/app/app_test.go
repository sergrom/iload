package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRow(t *testing.T) {
	tests := []struct {
		Row      string
		Url      string
		FileName string
	}{
		{Row: "", Url: "", FileName: ""},
		{Row: "http://something.org", Url: "http://something.org", FileName: ""},
		{Row: "http://something.org|filename", Url: "http://something.org", FileName: "filename"},
	}

	a := assert.New(t)

	for _, test := range tests {
		url, fileName := parseRow(test.Row)

		a.Equal(test.Url, url)
		a.Equal(test.FileName, fileName)
	}
}
