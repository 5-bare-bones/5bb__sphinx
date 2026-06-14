package command_helper

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/5-bare-bones/5bb__sphinx/config"
	"github.com/5-bare-bones/5bb__sphinx/orderedmap"
	"github.com/5-bare-bones/5bb__sphinx/protobuf"
	"github.com/5-bare-bones/5bb__sphinx/vault/card"
	"github.com/5-bare-bones/5bb__sphinx/vault/entry"
	"github.com/5-bare-bones/5bb__sphinx/vault/file"
	"github.com/5-bare-bones/5bb__sphinx/vault/totp"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	bolt "go.etcd.io/bbolt"
)

func TestBuildBox(t *testing.T) {
	expected := `╭────── Box ─────╮
│ Jedi   │ Luke  │
│ Hobbit │ Frodo │
│        │ Sam   │
│ Wizard │ Harry │
╰────────────────╯`

	mp := orderedmap.New()
	mp.Set("Jedi", "Luke")
	mp.Set("Hobbit", `Frodo
Sam`)
	mp.Set("Wizard", "Harry")

	got := BuildBox("test/box", mp)
	assert.Equal(t, expected, got)
}

func TestErase(t *testing.T) {
	f, err := os.CreateTemp("", "")
	assert.NoError(t, err, "Failed creating temporary file")
	f.Close()

	err = Erase(f.Name())
	assert.NoError(t, err, "Failed erasing file")

	err = Erase(f.Name())
	assert.Error(t, err, "Expected the file to be erased")
}

func TestExistsTrue(t *testing.T) {
	vault := SetContext(t)

	name := "naboo/tatooine"
	createObjects(t, vault, name)

	cases := []struct {
		desc   string
		object object
	}{
		{
			desc:   "card",
			object: Card,
		},
		{
			desc:   "entry",
			object: Entry,
		},
		{
			desc:   "file",
			object: File,
		},
		{
			desc:   "totp",
			object: TOTP,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			err := Exists(vault, name, tc.object)
			assert.Error(t, err)

			err = Exists(vault, "naboo/tatooine/hoth", tc.object)
			assert.Error(t, err)

			err = Exists(vault, "naboo", tc.object)
			assert.Error(t, err)
		})
	}
}

func TestExistsFalse(t *testing.T) {
	vault := SetContext(t)

	cases := []struct {
		desc   string
		name   string
		object object
	}{
		{
			desc:   "card",
			name:   "test",
			object: Card,
		},
		{
			desc:   "entry",
			name:   "test",
			object: Entry,
		},
		{
			desc:   "file",
			name:   "testing/test",
			object: File,
		},
		{
			desc:   "totp",
			name:   "testing",
			object: TOTP,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			err := Exists(vault, tc.name, tc.object)
			assert.NoError(t, err)
		})
	}
}

func TestFmtExpires(t *testing.T) {
	cases := []struct {
		desc     string
		expires  string
		expected string
	}{
		{desc: "Never keyword", expires: "Never", expected: "Never"},
		{desc: "Never case-insensitive", expires: "never", expected: "Never"},
		{desc: "empty is Never", expires: "", expected: "Never"},
		{desc: "blank is Never", expires: "   ", expected: "Never"},
		{desc: "zero is Never", expires: "0", expected: "Never"},
		{desc: "full ISO date", expires: "2029-06-26", expected: "2029-06-26"},
		{desc: "full ISO date is trimmed", expires: " 2029-06-26 ", expected: "2029-06-26"},
		{desc: "year only -> end of year", expires: "2029", expected: "2029-12-31"},
		{desc: "month only -> end of 30-day month", expires: "2029-06", expected: "2029-06-30"},
		{desc: "month only -> end of 31-day month", expires: "2029-07", expected: "2029-07-31"},
		{desc: "February in a leap year", expires: "2024-02", expected: "2024-02-29"},
		{desc: "February in a common year", expires: "2025-02", expected: "2025-02-28"},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := FormatExpires(tc.expires)
			assert.NoError(t, err, "unexpected error formatting %q", tc.expires)
			assert.Equal(t, tc.expected, got)
		})
	}

	invalid := []struct {
		desc    string
		expires string
	}{
		{desc: "non-numeric", expires: "invalid format"},
		{desc: "legacy slash format", expires: "26/06/2029"},
		{desc: "month out of range", expires: "2029-13"},
		{desc: "day out of range for month", expires: "2029-06-31"},
		{desc: "impossible February day", expires: "2025-02-30"},
		{desc: "day zero", expires: "2029-06-00"},
		{desc: "too many components", expires: "2029-06-26-01"},
		{desc: "non-numeric month", expires: "2029-AB"},
	}

	for _, tc := range invalid {
		t.Run("invalid: "+tc.desc, func(t *testing.T) {
			_, err := FormatExpires(tc.expires)
			assert.Error(t, err, "expected %q to be rejected", tc.expires)
		})
	}
}

func TestMustExist(t *testing.T) {
	vault := SetContext(t)

	name := "test/testing"
	createObjects(t, vault, name)

	t.Run("Success", func(t *testing.T) {
		objects := []object{Card, Entry, File, TOTP}
		for _, obj := range objects {
			err := MustExist(vault, obj)(nil, []string{name})
			assert.NoError(t, err)
		}
	})

	t.Run("Fail", func(t *testing.T) {
		cases := []struct {
			desc       string
			name       string
			errMessage string
		}{
			{
				desc:       "Record does not exist",
				name:       "test",
				errMessage: "\"test\" does not exist. Did you mean \"test/testing\"?",
			},
			{
				desc:       "Record does not exist 2",
				name:       "test/tstng",
				errMessage: "\"test/tstng\" does not exist. Did you mean \"test/testing\"?",
			},
			{
				desc:       "Empty name",
				name:       "",
				errMessage: "invalid name",
			},
			{
				desc:       "Invalid name",
				name:       "test//testing",
				errMessage: "invalid name",
			},
		}

		for _, tc := range cases {
			t.Run(tc.desc, func(t *testing.T) {
				err := MustExist(vault, Card)(nil, []string{tc.name})
				assert.Error(t, err)

				assert.Equal(t, tc.errMessage, err.Error())
			})
		}
	})

	t.Run("Empty args", func(t *testing.T) {
		err := MustExist(vault, Card)(nil, []string{})
		assert.Error(t, err)
	})

	t.Run("Directories", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			err := MustExist(vault, Card, true)(nil, []string{"test/"})
			assert.NoError(t, err)
		})

		t.Run("Not exists", func(t *testing.T) {
			err := MustExist(vault, Card, true)(nil, []string{"unexistent/"})
			assert.Error(t, err)
		})
	})
}

func TestMustExistList(t *testing.T) {
	vault := SetContext(t)
	cmd := &cobra.Command{}
	cmd.Flags().Bool("filter", false, "")
	objects := []object{Card, Entry, File, TOTP}

	name := "test"
	createObjects(t, vault, name)

	cases := []struct {
		desc   string
		name   string
		filter bool
	}{
		{
			desc: "Found name",
			name: name,
		},
		{
			desc: "Empty name",
			name: "",
		},
		{
			desc:   "Filtering",
			name:   "t",
			filter: true,
		},
	}

	t.Run("Success", func(t *testing.T) {
		for _, tc := range cases {
			t.Run(tc.desc, func(t *testing.T) {
				for _, obj := range objects {
					cmd.Args = MustExistList(vault, obj)
					cmd.Flags().Set("filter", strconv.FormatBool(tc.filter))

					err := cmd.Args(cmd, []string{tc.name})
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("Fail", func(t *testing.T) {
		cmd.Args = MustExistList(vault, Entry)
		cmd.Flag("filter").Changed = false

		err := cmd.Args(cmd, []string{"non-existent"})
		assert.Error(t, err)
	})
}

func TestMustNotExist(t *testing.T) {
	vault := SetContext(t)
	cmd := &cobra.Command{}
	objects := []object{Card, Entry, File, TOTP}

	t.Run("Success", func(t *testing.T) {
		for _, obj := range objects {
			cmd.Args = MustNotExist(vault, obj)
			err := cmd.Args(cmd, []string{"test"})
			assert.NoError(t, err)
		}
	})

	t.Run("Fail", func(t *testing.T) {
		err := entry.Create(vault, &protobuf.Entry{Name: "test"})
		assert.NoError(t, err)
		err = entry.Create(vault, &protobuf.Entry{Name: "dir/"})
		assert.NoError(t, err)

		cases := []struct {
			desc     string
			name     string
			allowDir []bool
		}{
			{
				desc: "Exists",
				name: "test",
			},
			{
				desc:     "Directory exists",
				name:     "dir/",
				allowDir: []bool{true},
			},
			{
				desc: "Empty name",
				name: "",
			},
			{
				desc: "Invalid name",
				name: "testing//test",
			},
		}

		for _, tc := range cases {
			t.Run(tc.desc, func(t *testing.T) {
				cmd.Args = MustNotExist(vault, Entry, tc.allowDir...)
				err := cmd.Args(cmd, []string{tc.name})
				assert.Error(t, err)
			})
		}
	})

	t.Run("No arguments", func(t *testing.T) {
		cmd.Args = MustNotExist(vault, Entry, false)
		err := cmd.Args(cmd, []string{})
		assert.Error(t, err)
	})
}

func TestNormalizeName(t *testing.T) {
	cases := []struct {
		desc     string
		name     string
		expected string
		allowDir []bool
	}{
		{
			desc:     "Normalize",
			name:     " / Go/Forum / ",
			expected: "go/forum",
		},
		{
			desc:     "Empty",
			name:     "",
			expected: "",
		},
		{
			desc:     "Allow dir",
			name:     "testing/",
			expected: "testing/",
			allowDir: []bool{true},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := NormalizeName(tc.name, tc.allowDir...)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestSelectEditor(t *testing.T) {
	t.Run("Default editor", func(t *testing.T) {
		expected := "nano"
		config.Set("editor", expected)
		defer config.Reset()

		got := SelectEditor()
		assert.Equal(t, expected, got)
	})

	t.Run("EDITOR", func(t *testing.T) {
		expected := "editor"
		os.Setenv("EDITOR", expected)
		defer os.Unsetenv("EDITOR")

		got := SelectEditor()
		assert.Equal(t, expected, got)
	})

	t.Run("VISUAL", func(t *testing.T) {
		expected := "visual"
		os.Setenv("VISUAL", expected)
		defer os.Unsetenv("VISUAL")

		got := SelectEditor()
		assert.Equal(t, expected, got)
	})

	t.Run("Default", func(t *testing.T) {
		got := SelectEditor()
		assert.Equal(t, "micro", got)
	})
}

func TestSupportedManagers(t *testing.T) {
	t.Run("Supported", func(t *testing.T) {
		list := []string{"1password", "bitwarden", "keepass", "keepassx", "keepassxc", "lastpass"}
		for _, name := range list {
			err := SupportedManagers()(nil, []string{name})
			assert.NoError(t, err)
		}
	})

	t.Run("Unsupported", func(t *testing.T) {
		list := []string{"", "unsupported"}
		for _, name := range list {
			err := SupportedManagers()(nil, []string{name})
			assert.Error(t, err)
		}
	})
}

func TestWatchFile(t *testing.T) {
	f, err := os.CreateTemp("", "*")
	assert.NoError(t, err)
	defer f.Close()

	_, err = f.Write([]byte("test"))
	assert.NoError(t, err)

	done := make(chan struct{}, 1)
	errCh := make(chan error, 1)
	go WatchFile(f.Name(), done, errCh)

	// Sleep to write after the file is being watched
	time.Sleep(50 * time.Millisecond)
	_, err = f.Write([]byte("test-watch-file"))
	assert.NoError(t, err)

	select {
	case <-done:

	case <-errCh:
		t.Errorf("Watching file failed: %v", err)
	}
}

func TestWatchFileErrors(t *testing.T) {
	cases := []struct {
		desc     string
		filename string
		initial  bool
	}{
		{
			desc:     "Initial stat error",
			filename: "test_error.json",
			initial:  true,
		},
		{
			desc:     "For loop stat error",
			filename: "test_error.json",
			initial:  false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			done := make(chan struct{}, 1)
			errCh := make(chan error, 1)

			if !tc.initial {
				err := os.WriteFile(tc.filename, []byte("test error"), 0o644)
				assert.NoError(t, err, "Failed creating file")
			}

			go WatchFile(tc.filename, done, errCh)

			// Sleep to wait until the file is created and fail once inside the for loop
			if !tc.initial {
				time.Sleep(10 * time.Millisecond)
				err := os.Remove(tc.filename)
				assert.NoError(t, err, "Failed removing the file")
			}

			select {
			case <-done:
				t.Error("Expected an error and it succeeded")

			case <-errCh:
			}
		})
	}
}

func TestWriteClipboard(t *testing.T) {
	if clipboard.Unsupported {
		t.Skip("No clipboard utilities available")
	}

	cmd := &cobra.Command{}

	t.Run("Default timeout", func(t *testing.T) {
		config.Set("clipboard.timeout", time.Nanosecond)
		defer config.Reset()

		err := WriteClipboard(cmd, 0, "", "test")
		assert.NoError(t, err)

		got, err := clipboard.ReadAll()
		assert.NoError(t, err)

		assert.Empty(t, got, "Expected the clipboard to be empty")
	})

	t.Run("t > 0", func(t *testing.T) {
		err := WriteClipboard(cmd, time.Nanosecond, "", "test")
		assert.NoError(t, err)

		got, err := clipboard.ReadAll()
		assert.NoError(t, err)

		assert.Empty(t, got, "Expected the clipboard to be empty")
	})

	t.Run("t = 0", func(t *testing.T) {
		clip := "test"
		err := WriteClipboard(cmd, 0, "", clip)
		assert.NoError(t, err)

		got, err := clipboard.ReadAll()
		assert.NoError(t, err)

		assert.Equal(t, clip, got)
	})
}

func TestFormatSuggestions(t *testing.T) {
	cases := []struct {
		desc           string
		expectedResult string
		suggestions    []string
	}{
		{
			desc: "One suggestion",
			suggestions: []string{
				"car",
			},
			expectedResult: "\"car\"",
		},
		{
			desc: "Two suggestions",
			suggestions: []string{
				"car",
				"cur",
			},
			expectedResult: "\"car\" or \"cur\"",
		},
		{
			desc: "Three suggestions",
			suggestions: []string{
				"car",
				"cur",
				"core",
			},
			expectedResult: "\"car\", \"cur\" or \"core\"",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			suggestion := formatSuggestions(tc.suggestions)

			assert.Equal(t, tc.expectedResult, suggestion)
		})
	}
}

func TestGetNameSuggestions(t *testing.T) {
	cases := []struct {
		desc                string
		names               []string
		name                string
		expectedSuggestions []string
	}{
		{
			desc: "By distance",
			names: []string{
				"show",
				"bat",
				"category",
				"rat",
				"car",
				"pop",
			},
			name: "hay",
			expectedSuggestions: []string{
				"show",
				"bat",
				"rat",
				"car",
			},
		},
		{
			desc: "By prefix",
			names: []string{
				"category",
				"careful",
				"career",
			},
			name: "car",
			expectedSuggestions: []string{
				"careful",
				"career",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			suggestions := getNameSuggestions(tc.name, tc.names)

			assert.Equal(t, tc.expectedSuggestions, suggestions)
		})
	}
}

func TestLevenshteinDistance(t *testing.T) {
	cases := []struct {
		nameA            string
		nameB            string
		expectedDistance int
	}{
		{
			nameA:            "car",
			nameB:            "show",
			expectedDistance: 1,
		},
		{
			nameA:            "car",
			nameB:            "hat",
			expectedDistance: 2,
		},
		{
			nameA:            "car",
			nameB:            "carry",
			expectedDistance: 2,
		},
		{
			nameA:            "car",
			nameB:            "born",
			expectedDistance: 3,
		},
		{
			nameA:            "car",
			nameB:            "karting",
			expectedDistance: 5,
		},
	}

	for _, tc := range cases {
		t.Run(strconv.Itoa(tc.expectedDistance), func(t *testing.T) {
			distance := levenshteinDistance(tc.nameA, tc.nameB)

			assert.Equal(t, tc.expectedDistance, distance)
		})
	}
}

func createObjects(t *testing.T, vault *bolt.DB, name string) {
	t.Helper()
	err := entry.Create(vault, &protobuf.Entry{Name: name})
	assert.NoError(t, err)
	err = card.Create(vault, &protobuf.Card{Name: name})
	assert.NoError(t, err)
	err = file.Create(vault, &protobuf.File{Name: name})
	assert.NoError(t, err)
	err = totp.Create(vault, &protobuf.TOTP{Name: name})
	assert.NoError(t, err)
}
