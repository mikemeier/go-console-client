package app

import (
	"console/both/message"
	"console/client/socket"
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"os"
	"strings"
)

var app *tview.Application
var history = []string{}
var historyPosition = 0
var savedOutput []string
var savedPrefix string
var hasInputMask bool

func Run(host string) {
	app = tview.NewApplication()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			if socket.Send(message.ExitCommand) {
				return nil
			}
		}
		return event
	})

	outputText := tview.NewTextView()
	outputText.SetScrollable(true)
	outputText.SetDynamicColors(true)
	outputText.SetBorder(false)

	inputField := tview.NewInputField()
	inputField.SetFieldBackgroundColor(tcell.ColorBlack)
	inputField.SetLabel(fmt.Sprintf("Connecting %s ...", host))
	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			text := inputField.GetText()
			if text != "" {
				socket.Send(text)
				inputField.SetText("")

				if !hasInputMask {
					history = append(history, text)
					historyPosition = 0
				}
			}
		}
	})
	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()

		if key == tcell.KeyUp || key == tcell.KeyDown {
			historyLen := len(history)
			if key == tcell.KeyUp {
				historyPosition--
				if historyPosition*-1 > historyLen {
					historyPosition = historyLen * -1
				}
			} else {
				historyPosition++
				if historyPosition > 0 {
					historyPosition = 0
				}
			}

			if historyPosition == 0 {
				inputField.SetText("")
			} else {
				inputField.SetText(history[historyLen+historyPosition])
			}
		}

		return event
	})

	flex := tview.NewFlex()
	flex.SetFullScreen(true)
	flex.SetDirection(tview.FlexRow)
	flex.SetBorder(false)
	flex.SetBackgroundColor(tcell.ColorBlack)

	flex.AddItem(outputText, 0, 100, false)
	flex.AddItem(inputField, 0, 1, true)

	go func() {
		for {
			select {
			case isConnected := <-socket.IsConnected:
				if !isConnected {
					outputText.SetText("")
					inputField.SetLabel(fmt.Sprintf("Connecting to %s... (#%d)", host, socket.Retries))
					app.Draw()
				}
			case msg := <-socket.Messages:
				switch msg.Command {
				case message.ExitCommand:
					socket.Disconnect()
					app.Stop()
					os.Exit(0)
				case message.ClearCommand:
					outputText.SetText("")
					savedOutput = []string{}
				case message.RestoreOutputCommand:
					outputText.SetText(strings.Join(savedOutput, "\n") + "\n")
				case message.SavePrefixCommand:
					savedPrefix = inputField.GetLabel()
				case message.RestorePrefixCommand:
					inputField.SetLabel(savedPrefix)
				case message.SetInputTypePasswordCommand:
					hasInputMask = true
					inputField.SetText("")
					inputField.SetMaskCharacter('*')
				case message.SetInputTypeTextCommand:
					hasInputMask = false
					inputField.SetText("")
					inputField.SetMaskCharacter(0)
				}

				if msg.Prefix != "" {
					inputField.SetLabel(msg.Prefix)
				}

				if msg.Message != "" {
					if !msg.NoHistory {
						savedOutput = append(savedOutput, msg.Message)
					}

					if msg.Command == message.ReplaceCommand {
						outputText.SetText(msg.Message)
					} else {
						fmt.Fprintln(outputText, msg.Message)
					}
				}

				app.Draw()
			}
		}
	}()

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}

	socket.Disconnect()
	os.Exit(0)
}
