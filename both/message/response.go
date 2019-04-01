package message

const ExitCommand = "exit"
const ClearCommand = "clear"
const ReplaceCommand = "replace"

const RestoreOutputCommand = "restore-output"

const SavePrefixCommand = "save-prefix"
const RestorePrefixCommand = "restore-prefix"

const SetInputTypePasswordCommand = "input-type-password"
const SetInputTypeTextCommand = "input-type-text"

const ClearMessage = "^C"

const BufferLength = 5

func NewEmptyLine() *Response {
	return &Response{Message: "\n"}
}

type Response struct {
	Prefix    string
	Message   string
	Command   string
	NoHistory bool
}
