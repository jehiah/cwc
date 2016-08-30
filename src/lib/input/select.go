package input

import (
	"bytes"
	"fmt"
	"strconv"
)

// Select asks the user to select a item from the given list by the number.
// It shows the given query and list to user. The response is returned as string
// from the list. By default, it checks the input is the number and is not
// out of range of the list and if not returns error. If Loop is true, it continue to
// ask until it receives valid input.
//
// If the user sends SIGINT (Ctrl+C) while reading input, it catches
// it and return it as a error.
func (i *UI) Select(query string, dflt interface{}, list ...interface{}) (interface{}, error) {

	// Find default index which opts.Default indicates
	defaultIndex := -1

	if dflt != nil {
		var defaultVal string
		switch dflt := dflt.(type) {
		case string:
			defaultVal = dflt
		case fmt.Stringer:
			defaultVal = dflt.String()
		default:
			return "", fmt.Errorf("default does not support fmt.Stringer")
		}

		for i, item := range list {
			switch item := item.(type) {
			case string:
				if item == defaultVal {
					defaultIndex = i
				}
			case fmt.Stringer:
				if item.String() == defaultVal {
					defaultIndex = i
				}
			default:
				return "", fmt.Errorf("list[%d] %T does not support fmt.Stringer", i, list[i])
			}
		}

		// DefaultVal is set but doesn't exist in list
		if defaultIndex == -1 {
			// This error message is not for user
			// Should be found while development
			return "", fmt.Errorf("default is specified but item does not exist in list")
		}
	}

	// Construct the query & display it to user
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s\n\n", query))
	for i, item := range list {
		buf.WriteString(fmt.Sprintf("%d. %s\n", i+1, item))
	}

	buf.WriteString("\n")
	fmt.Fprintf(i.Writer, buf.String())

	// resultStr and resultErr are return val of this function
	var result interface{}
	var resultErr error
	for {

		// Construct the asking line to input
		var buf bytes.Buffer
		buf.WriteString("Enter a number")

		// Add default val if provided
		if defaultIndex >= 0 {
			buf.WriteString(fmt.Sprintf(" (Default is %d)", defaultIndex+1))
		}

		buf.WriteString(": ")
		fmt.Fprintf(i.Writer, buf.String())

		// Read user input from reader.
		line, err := i.readline()
		if err != nil {
			resultErr = err
			break
		}

		// line is empty but default is provided returns it
		if line == "" && defaultIndex >= 0 {
			result = list[defaultIndex]
			break
		}

		//  && opts.Required
		if line == "" {
			// if !opts.Loop {
			// 	resultErr = ErrEmpty
			// 	break
			// }

			fmt.Fprintf(i.Writer, "Input must not be empty. Answer by a number.\n\n")
			continue
		}

		// Convert user input string to int val
		n, err := strconv.Atoi(line)
		if err != nil {
			// if !opts.Loop {
			// 	resultErr = ErrNotNumber
			// 	break
			// }

			fmt.Fprintf(i.Writer,
				"%q is not a valid input. Answer by a number.\n\n", line)
			continue
		}

		// Check answer is in range of list
		if n < 1 || len(list) < n {
			// if !opts.Loop {
			// 	resultErr = ErrOutOfRange
			// 	break
			// }

			fmt.Fprintf(i.Writer,
				"%q is not a valid choice. Choose a number from 1 to %d.\n\n",
				line, len(list))
			continue
		}

		// // validate input by custom function
		// if v != nil {
		// 	if err := v(line)
		// }
		// validate := opts.validateFunc()
		// if err := validate(line); err != nil {
		// 	if !opts.Loop {
		// 		resultErr = err
		// 		break
		// 	}
		//
		// 	fmt.Fprintf(i.Writer, "Failed to validate input string: %s\n\n", err)
		// 	continue
		// }

		// Reach here means it gets ideal input.
		result = list[n-1]
		break
	}

	// Insert the new line for next output
	fmt.Fprintf(i.Writer, "\n")

	return result, resultErr
}

func SelectString(query string, dflt string, list ...string) (string, error) {
	newList := make([]interface{}, len(list))
	for i, v := range list {
		newList[i] = v
	}
	var d interface{} = dflt
	s, err := Select(query, d, newList...)
	if err != nil {
		return "", err
	}
	return s.(string), nil
}

func Select(query string, dflt interface{}, list ...interface{}) (interface{}, error) {
	return Default.Select(query, dflt, list...)
}
