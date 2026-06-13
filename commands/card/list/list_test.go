package list

import (
	"testing"

	cmdutil "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/db/card"
	"github.com/5-bare-bones/5bb__sphinx/pb"

	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	db := cmdutil.SetContext(t)

	err := card.Create(db, &pb.Card{
		Name:   "test",
		Number: "1500135",
	})
	assert.NoError(t, err, "Failed creating the card")

	cases := []struct {
		desc   string
		name   string
		filter string
		hide   string
		qr     string
	}{
		{
			desc: "List one",
			name: "test",
		},
		{
			desc: "List one and show qr",
			name: "test",
			qr:   "true",
		},
		{
			desc:   "Filter by name",
			name:   "te*",
			filter: "true",
		},
		{
			desc: "List all",
			name: "",
		},
		{
			desc: "List one and hide",
			name: "test",
			hide: "true",
		},
	}

	cmd := NewCmd(db)
	f := cmd.Flags()

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			cmd.SetArgs([]string{tc.name})
			f.Set("filter", tc.filter)
			f.Set("hide", tc.hide)
			f.Set("qr", tc.qr)

			err := cmd.Execute()
			assert.NoError(t, err)
		})
	}
}

func TestListErrors(t *testing.T) {
	db := cmdutil.SetContext(t)

	err := card.Create(db, &pb.Card{Name: "test"})
	assert.NoError(t, err, "Failed creating the card")

	cases := []struct {
		desc   string
		name   string
		filter string
		qr     string
	}{
		{
			desc:   "Card does not exist",
			name:   "non-existent",
			filter: "false",
		},
		{
			desc:   "No cards found",
			name:   "non-existent",
			filter: "true",
		},
		{
			desc:   "Filter syntax error",
			name:   "[error",
			filter: "true",
		},
		{
			desc: "No data to encode",
			name: "test",
			qr:   "true",
		},
	}

	cmd := NewCmd(db)

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			cmd.SetArgs([]string{tc.name})
			f := cmd.Flags()
			f.Set("filter", tc.filter)
			f.Set("qr", tc.qr)

			err := cmd.Execute()
			assert.Error(t, err)
		})
	}
}

func TestPostRun(t *testing.T) {
	NewCmd(nil).PostRun(nil, nil)
}
