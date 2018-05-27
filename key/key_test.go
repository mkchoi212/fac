package key

import (
	"errors"
	"os"
	"os/user"
	"sort"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/mkchoi212/fac/testhelper"
)

var tests = []struct {
	settings Binding
	expected Binding
	warnings []string
	fatals   []string
}{
	{
		settings: Binding{
			"foobar": "a",
			"hello":  "b",
		},
		expected: defaultBinding,
		warnings: []string{
			"Invalid key: \"foobar\" will be ignored",
			"Invalid key: \"hello\" will be ignored",
		},
		fatals: nil,
	},
	{
		settings: Binding{
			SelectLocal:    "i",
			SelectIncoming: "incoming",
		},
		expected: Binding{
			SelectLocal:    "l",
			SelectIncoming: "i",
		},
		warnings: []string{
			"Illegal multi-character mapping: \"incoming\" will be interpreted as 'i'",
		},
		fatals: []string{
			"Duplicate key-mapping: \"select_incoming, select_local\" are all represented by 'i'",
		},
	},
	{
		settings: Binding{
			ShowLinesDown: "d",
			ShowLinesUp:   "u",
		},
		expected: Binding{
			ShowLinesDown: "d",
			ShowLinesUp:   "u",
		},
		warnings: nil,
		fatals:   nil,
	},
}

func TestLoadSettings(t *testing.T) {
	currentUser = func() (*user.User, error) {
		return &user.User{HomeDir: "."}, nil
	}

	for _, test := range tests {
		if test.fatals != nil {
			continue
		}

		// Create dummy yml file with content
		f, err := os.Create(".fac.yml")
		testhelper.Ok(t, err)
		data, _ := yaml.Marshal(&test.settings)
		f.WriteString(string(data))
		f.Close()

		output, err := LoadSettings()
		testhelper.Ok(t, err)

		test.expected.consolidate()
		testhelper.Equals(t, test.expected, output)
	}
}

func TestParseSettings(t *testing.T) {
	// Test with invalid currentUser
	currentUser = func() (*user.User, error) {
		return nil, errors.New("Could not find current user")
	}
	binding, _ := parseSettings()
	testhelper.Equals(t, defaultBinding, binding)

	// Test with invalid directory
	currentUser = func() (*user.User, error) {
		return &user.User{HomeDir: "foobar"}, nil
	}
	binding, _ = parseSettings()
	testhelper.Equals(t, defaultBinding, binding)

	// Test with valid directory with empty file
	currentUser = func() (*user.User, error) {
		return &user.User{HomeDir: "."}, nil
	}
	f, err := os.Create(".fac.yml")
	testhelper.Ok(t, err)
	defer f.Close()

	binding, _ = parseSettings()
	testhelper.Equals(t, 0, len(binding))

	// Test valid directory with erroneous content
	f.WriteString("erroneous content")
	binding, _ = parseSettings()
	testhelper.Equals(t, defaultBinding, binding)
}

func TestVerify(t *testing.T) {
	for _, test := range tests {
		warnings, fatals := test.settings.verify()
		sort.Strings(warnings)
		sort.Strings(test.warnings)
		testhelper.Equals(t, test.warnings, warnings)
		testhelper.Equals(t, test.fatals, fatals)
	}
}

func TestSummary(t *testing.T) {
	summary := defaultBinding.Summary()
	testhelper.Equals(t, summary, "[w,a,s,d,e,?] >>")
}

func TestHelp(t *testing.T) {
	helpMsg := defaultBinding.Help()
	testhelper.Assert(t, len(helpMsg) != 0, "Help message should not be of length 0")
}

func TestConsolidate(t *testing.T) {
	for _, test := range tests {
		test.settings.consolidate()

		for k := range defaultBinding {
			_, ok := test.settings[k]
			testhelper.Equals(t, true, ok)
		}
	}
}
