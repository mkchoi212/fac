package binding

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os/user"
	"sort"
	"strings"
	"time"

	"github.com/mkchoi212/fac/color"
	yaml "gopkg.in/yaml.v2"
)

var currentUser = user.Current

// Binding represents the user's key binding configuration
type Binding map[string]string

// Following constants represent all the actions available to the user
// The string literals are used to retrieve values from `Binding` and
// when writing/reading from .fac.yml
const (
	SelectLocal           = "select_local"
	SelectIncoming        = "select_incoming"
	ToggleViewOrientation = "toggle_view"
	ShowLinesUp           = "show_up"
	ShowLinesDown         = "show_down"
	ScrollUp              = "scroll_up"
	ScrollDown            = "scroll_down"
	EditCode              = "edit"
	NextConflict          = "next"
	PreviousConflict      = "previous"
	QuitApplication       = "quit"
	ShowHelp              = "help"
	ContinuousEvaluation  = "cont_eval"
)

// defaultBinding is used when the user has not specified any of the
// available actions via `.fac.yml`
var defaultBinding = Binding{
	SelectLocal:           "a",
	SelectIncoming:        "d",
	ToggleViewOrientation: "v",
	ShowLinesUp:           "w",
	ShowLinesDown:         "s",
	ScrollUp:              "j",
	ScrollDown:            "k",
	EditCode:              "e",
	NextConflict:          "n",
	PreviousConflict:      "p",
	QuitApplication:       "q",
	ShowHelp:              "h",
	ContinuousEvaluation:  "false",
}

// LoadSettings looks for a user specified key-binding settings file - `$HOME/.fac.yml`
// and returns a map representation of the file
// It also looks for errors, and ambiguities within the file and notifies the user of them
func LoadSettings() (b Binding, err error) {
	warnings := []string{}

	b, ok := parseSettings()
	if ok != nil {
		warnings = append(warnings, ok.Error())
	}

	verificationWarnings, fatals := b.verify()
	warnings = append(warnings, verificationWarnings...)

	if len(warnings) != 0 {
		fmt.Println(color.Yellow(color.Regular, "⚠️  %d infraction(s) detected in .fac.yml", len(warnings)))
		for _, msg := range warnings {
			fmt.Println(color.Yellow(color.Regular, "%s", msg))
		}
		fmt.Println()
	}

	if len(fatals) != 0 {
		var errorMsg bytes.Buffer
		errorMsg.WriteString(color.Red(color.Regular, "🚫  %d unrecoverable error(s) detected in .fac.yml", len(fatals)))
		for _, msg := range fatals {
			errorMsg.WriteString(color.Red(color.Regular, "\n%s", msg))
		}
		return nil, errors.New(errorMsg.String())
	}

	if len(warnings) != 0 {
		time.Sleep(time.Duration(2) * time.Second)
	}

	b.consolidate()
	return
}

// parseSettings looks for `$HOME/.fac.yml` and parses it into a `Binding` value
// If the file does not exist, it returns the `defaultBinding`
func parseSettings() (Binding, error) {
	userBinding := make(Binding)

	usr, err := currentUser()
	if err != nil {
		return defaultBinding, err
	}

	// Read config file
	f, err := ioutil.ReadFile(usr.HomeDir + "/.fac.yml")
	if err != nil {
		return defaultBinding, err
	}

	// Parse config file
	if err = yaml.Unmarshal(f, &userBinding); err != nil {
		return defaultBinding, err
	}

	return userBinding, nil
}

// consolidate takes the user's key-binding settings and fills the missings key-binds
// with the default key-binding values
func (b Binding) consolidate() {
	for key, defaultValue := range defaultBinding {
		userValue, ok := b[key]

		if !ok || userValue == "" {
			b[key] = defaultValue
		} else if len(userValue) > 1 && userValue != "false" && userValue != "true" {
			b[key] = string(userValue[0])
		}
	}
}

// verify looks through the user's key-binding settings and looks for any infractions such as..
// 1. Invalid/ignored key-binding keys
// 2. Multi-character key-mappings (except for `cont_eval`)
// 3. Duplicate key-mappings
func (b Binding) verify() (warnings []string, fatals []string) {
	bindTable := map[string][]string{}

	for k, v := range b {
		// Check for "1. Invalid/ignored key-binding keys"
		if _, ok := defaultBinding[k]; !ok {
			warnings = append(warnings, fmt.Sprintf("Invalid key: \"%s\" will be ignored", k))
			delete(b, k)
			continue
		}

		// Check for "2. Multi-character key-mappings"
		if k == ContinuousEvaluation && v != "false" && v != "true" {
			fatals = append(fatals, fmt.Sprintf("Invalid value: value for key '%s' must either be true or false", ContinuousEvaluation))
			continue
		} else if len(v) > 1 && k != ContinuousEvaluation {
			abbreviated := string(v[0])
			warnings = append(warnings, fmt.Sprintf("Illegal multi-character mapping: \"%s\" will be interpreted as '%s'", v, abbreviated))
			bindTable[abbreviated] = append(bindTable[abbreviated], k)
		} else {
			bindTable[v] = append(bindTable[v], k)
		}
	}

	// Check for "3. Duplicate key-mappings"
	for k, v := range bindTable {
		if len(v) > 1 {
			sort.Strings(v)
			duplicateValues := strings.Join(v, ", ")
			fatals = append(fatals, fmt.Sprintf("Duplicate key-mapping: \"%s\" are all represented by '%s'", duplicateValues, k))
		}
	}

	return
}

// Summary returns a short summary of the provided `Binding`
// and is used as the helpful string displayed by the user's input field
// e.g. "[w,a,s,d,e,?] >>"
func (b Binding) Summary() string {
	targetKeys := []string{
		b[ShowLinesUp],
		b[SelectLocal],
		b[ShowLinesDown],
		b[SelectIncoming],
		b[EditCode],
	}
	return "[" + strings.Join(targetKeys, ",") + ",?] >>"
}

// Help returns a help string that is displayed on the right panel of the UI
// It should provided an overall summary of all available key options
func (b Binding) Help() string {
	format := `
	%s - use local version
	%s - use incoming version
	%s - manually edit code

	%s - show more lines up
	%s - show more lines down
	%s - scroll up
	%s - scroll down
	
	%s - view orientation
	%s - next conflict
	%s - previous conflict
	
	%s | ? - help
	%s | Ctrl+C - quit
	`

	return fmt.Sprintf(format, b[SelectLocal], b[SelectIncoming], b[EditCode],
		b[ShowLinesUp], b[ShowLinesDown],
		b[ScrollUp], b[ScrollDown], b[ToggleViewOrientation], b[NextConflict], b[PreviousConflict],
		b[ShowHelp], b[QuitApplication])
}
