//go:generate goversioninfo -icon=clipboard.ico -skip-versioninfo=true
package main

import (
	"context"
	_ "embed"
	"github.com/getlantern/systray"
	"golang.design/x/clipboard"
	"log"
	"regexp"
)

//go:embed clipboard.ico
var icon []byte

// A struct to represent the state of the application.
type appState struct {
	// Whether the application is currently monitoring the clipboard.
	monitoring     bool
	substituteWith string
}

// A function to replace all the single occurrences of line breaks and the surrounding spaces to a single space,
// while reserving all the continuous line breaks.
func normalizeLineBreaks(text string, substitute string) string {
	r, err := regexp.Compile(`[\t\f\r ]*\n[\t\f ]*`)
	if err != nil {
		log.Println(err)
		return ""
	}
	if r.FindAllString(text, -1) == nil {
		return text
	}
	result := r.ReplaceAllString(text, substitute)
	return result
}

// A function to monitor the clipboard for changes.
func monitorClipboard(ctx context.Context, appState *appState) {
	clipboardChanged := clipboard.Watch(ctx, clipboard.FmtText)
	for {
		select {
		case <-ctx.Done():
			return
		case rawContent := <-clipboardChanged:
			// If the application is not monitoring the clipboard, skip the change.
			if !appState.monitoring {
				log.Println("Clipboard changed but not monitoring")
				continue
			}
			stringContent := string(rawContent)
			// If the clipboard content is empty, skip the change.
			if len(stringContent) == 0 {
				log.Println("Clipboard changed but empty")
				continue
			}
			log.Println("Clipboard changed: " + stringContent)
			// Normalize the line breaks in the clipboard content.
			normalizedContent := normalizeLineBreaks(stringContent, appState.substituteWith)
			log.Println("Normalized content: " + normalizedContent)
			// Set the new contents of the clipboard.
			clipboard.Write(clipboard.FmtText, []byte(normalizedContent))
			log.Println("Clipboard changed to: " + normalizedContent)
		}
	}
}

func listenToggle(ctx context.Context, toggle *systray.MenuItem, appState *appState) {
	for {
		select {
		case <-toggle.ClickedCh:
			if appState.monitoring {
				toggle.SetTitle("Start Monitoring")
			} else {
				toggle.SetTitle("Stop Monitoring")
			}
			appState.monitoring = !appState.monitoring
			continue
		case <-ctx.Done():
			return
		}
	}
}

func listenToggleSubstitute(ctx context.Context, toggle *systray.MenuItem, appState *appState) {
	for {
		select {
		case <-toggle.ClickedCh:
			if appState.substituteWith == "" {
				toggle.SetTitle("Substitute With Empty")
				appState.substituteWith = " "
			} else {
				toggle.SetTitle("Substitute With Space")
				appState.substituteWith = ""
			}
			continue
		case <-ctx.Done():
			return
		}
	}
}

func listenQuit(quit *systray.MenuItem, cancel context.CancelFunc) {
	for {
		<-quit.ClickedCh
		cancel()
		systray.Quit()
	}
}

func onSystrayReady() {
	appState := appState{monitoring: true, substituteWith: " "}
	if len(icon) > 0 {
		systray.SetIcon(icon)
	}
	systray.SetTitle("Paste Without Line Breaks")
	systray.SetTooltip("Paste Without Line Breaks")
	mToggle := systray.AddMenuItem("Stop Monitoring", "Stop Monitoring")
	mToggleSubstitute := systray.AddMenuItem("Substitute With Empty", "Substitute With Empty")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	ctx, cancel := context.WithCancel(context.Background())
	go listenQuit(mQuit, cancel)
	go listenToggle(ctx, mToggle, &appState)
	go listenToggleSubstitute(ctx, mToggleSubstitute, &appState)
	go monitorClipboard(ctx, &appState)
}

func onSystrayExit() {

}

func main() {
	err := clipboard.Init()
	if err != nil {
		log.Panic(err)
	}
	systray.Run(onSystrayReady, onSystrayExit)
}
