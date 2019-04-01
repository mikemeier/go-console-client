package message

func parseCommand(input string) []string {
	var arguments []string
	var isInside bool
	var currentArgument string

	appendCurrentArgument := func() {
		if currentArgument != "" {
			arguments = append(arguments, currentArgument)
		}
		currentArgument = ""
	}

	for _, in := range input {
		char := string(in)

		if char == " " {
			if isInside {
				currentArgument += char
				continue
			}

			appendCurrentArgument()
			continue
		}

		if char == "\"" {
			if isInside {
				appendCurrentArgument()
				isInside = false
				continue
			}

			isInside = true
			continue
		}

		currentArgument += char
	}

	if currentArgument != "" {
		arguments = append(arguments, currentArgument)
	}

	return arguments
}
