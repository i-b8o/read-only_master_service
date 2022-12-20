package usecase_doc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIDs(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		input string
		regID string
		chID  string
		pID   string
	}{
		{
			input: "372952/4e92c731969781306ebd1095867d2385f83ac7af/335104",
			regID: "372952",
			chID:  "4e92c731969781306ebd1095867d2385f83ac7af",
			pID:   "335104",
		},
		{
			input: "/document/cons_doc_LAW_2875/",
			regID: "2875",
			chID:  "",
			pID:   "",
		},
	}

	for _, test := range tests {
		regID, chID, pID := getIDs(test.input)
		assert.Equal(test.regID, regID)
		assert.Equal(test.chID, chID)
		assert.Equal(test.pID, pID)
	}
}
