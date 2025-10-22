//nolint:forbidigo // enable printing function for prompt package
package prompt

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
)

func checkError(err error) {
	if err == nil {
		return
	}

	if errors.Is(err, promptui.ErrAbort) {
		return
	}

	if errors.Is(err, promptui.ErrInterrupt) ||
		errors.Is(err, promptui.ErrEOF) {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Prompt failed: %v\n", err)
	os.Exit(1)
}

// PromptPassword prompts the user to enter a password.
// If confirmation is true, the user will be asked to re-enter the password
// for confirmation.
func PromptPassword(label string, confirmation bool) string {
	prompt := promptui.Prompt{
		Label:   label,
		Mask:    '*',
		Pointer: promptui.PipeCursor,
	}
	password, err := prompt.Run()
	checkError(err)

	if confirmation {
		validate := func(input string) error {
			if input != password {
				return errors.New("passwords do not match")
			}

			return nil
		}

		confirmPrompt := promptui.Prompt{
			Label:    "Confirm password",
			Validate: validate,
			Mask:     '*',
			Pointer:  promptui.PipeCursor,
		}

		_, err := confirmPrompt.Run()
		checkError(err)
	}

	return password
}

// PromptConfirm prompts the user to confirm an operation.
// It returns true if the user confirms (Y/y), otherwise false.
// The program exits if the user aborts or if an unexpected error occurs.
func PromptConfirm(label string) bool {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
		Pointer:   promptui.PipeCursor,
	}
	result, err := prompt.Run()
	checkError(err)

	if result != "" && strings.ToUpper(result[:1]) == "Y" {
		return true
	}

	return false
}

// PromptInput prompts the user for a string input.
func PromptInput(label string) string {
	prompt := promptui.Prompt{
		Label:   label,
		Pointer: promptui.PipeCursor,
	}
	result, err := prompt.Run()
	checkError(err)

	return result
}

// PromptSelect displays a list of choices for the user to select from.
// It returns the index of the selected item.
func PromptSelect(label string, items []string) int {
	prompt := promptui.Select{
		Label:   label,
		Items:   items,
		Pointer: promptui.PipeCursor,
	}

	choice, _, err := prompt.Run()
	checkError(err)

	return choice
}

// PromptInputWithSuggestion prompts the user for a string input,
// showing a suggested default value.
func PromptInputWithSuggestion(label, suggestion string) string {
	prompt := promptui.Prompt{
		Label:   label,
		Default: suggestion,
		Pointer: promptui.PipeCursor,
	}
	result, err := prompt.Run()
	checkError(err)

	return result
}

// PromptInputWithRange prompts the user to enter an integer value
// within the specified range [min, max]. The default value is shown as a suggestion.
func PromptInputWithRange(label string, def, min, max int) int {
	prompt := promptui.Prompt{
		Label:   label,
		Default: fmt.Sprintf("%v", def),
		Pointer: promptui.PipeCursor,
		Validate: func(input string) error {
			num, err := strconv.Atoi(input)
			if err != nil {
				return fmt.Errorf("invalid number %q", input)
			}
			if num < min || num > max {
				return fmt.Errorf("enter a number between %v and %v", min, max)
			}

			return nil
		},
	}
	result, err := prompt.Run()
	checkError(err)

	num, _ := strconv.Atoi(result)

	return num
}
