package message

import "strings"

func NewRequest(input string) Request {
	return Request{input: input, arguments: parseCommand(input)}
}

type Request struct {
	input     string
	arguments []string
}

func (r Request) GetCommandName() string {
	return strings.ToLower(r.GetCommandNameRaw())
}

func (r Request) GetCommandNameRaw() string {
	if len(r.arguments) <= 0 {
		return ""
	}

	return r.arguments[0]
}

func (r Request) GetArguments() []string {
	if len(r.arguments) <= 1 {
		return []string{}
	}

	return r.arguments[1:]
}

func (r Request) GetArgumentsString() string {
	return strings.Join(r.GetArguments(), " ")
}

func (r Request) HasExactlyArgumentsN(n int) bool {
	return len(r.GetArguments()) == n
}

func (r Request) HasAtLeastArgumentsN(n int) bool {
	return len(r.GetArguments()) >= n
}

func (r Request) GetArgumentN(n int) string {
	n--
	arguments := r.GetArguments()
	if len(arguments)-1 < n {
		return ""
	}

	return arguments[n]
}

func (r Request) GetArgumentsFromN(n int) []string {
	n--
	arguments := r.GetArguments()
	if len(arguments)-1 < n {
		return []string{}
	}

	return arguments[n:]
}

func (r Request) GetArgumentsFromNAsString(n int) string {
	return strings.Join(r.GetArgumentsFromN(n), " ")
}