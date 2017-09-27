package headlessChrome

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/integrii/interactive"
)

// Debug enables debug output for this package to console
var Debug bool

// ChromePath is the command to execute chrome
var ChromePath = `/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`

// Args are the args that will be used to start chrome
var Args = []string{
	"--headless",
	"--disable-gpu",
	"--repl",
	// "--dump-dom",
	// "--window-size=1024,768",
	// "--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	// "--verbose",
}

const expectedFirstLine = `Type a Javascript expression to evaluate or "quit" to exit.`
const promptPrefix = `>>>`

// OutputSanitizer puts output coming from the consolw that
// does not begin with the input prompt into the session
// output channel
func (cs *ChromeSession) OutputSanitizer() {
	for text := range cs.Session.Output {
		if !strings.HasPrefix(text, promptPrefix) {
			cs.Output <- text
		}
	}
}

// ChromeSession is an interactive console Session with a Chrome
// instance.
type ChromeSession struct {
	Session *interactive.Session
	Output  chan string
	Input   chan string
}

// Exit exits the running command out by ossuing a 'quit'
// to the chrome console
func (cs *ChromeSession) Exit() {
	cs.Session.Write(`quit`)
	cs.Session.Exit()
}

// Write writes to the Session
func (cs *ChromeSession) Write(s string) {
	if Debug {
		fmt.Println("Writing to console:")
		fmt.Println(s)
	}
	cs.Session.Write(s)
}

// OutputPrinter prints all outputs from the output channel to the cli
func (cs *ChromeSession) OutputPrinter() {
	for l := range cs.Session.Output {
		fmt.Println(l)
	}
}

// forceClose issues a force kill to the command
func (cs *ChromeSession) forceClose() {
	cs.Session.ForceClose()
}

// ClickSelector calls a click() on the supplied selector
func (cs *ChromeSession) ClickSelector(s string) {
	cs.Write(`document.querySelector("` + s + `").click()`)
}

// ClickItemWithInnerHTML clicks an item that has the matching inner html
func (cs *ChromeSession) ClickItemWithInnerHTML(elementType string, s string, itemIndex int) {
	cs.Write(`var x = $("` + elementType + `").filter(function(idx) { return this.innerHTML.indexOf("` + s + `") == 0; });x[` + strconv.Itoa(itemIndex) + `].click()`)
}

// GetItemWithInnerHTML fetches the item with the specified innerHTML content
func (cs *ChromeSession) GetItemWithInnerHTML(elementType string, s string, itemIndex int) {
	cs.Write(`var x = $("` + elementType + `").filter(function(idx) { return this.innerHTML.indexOf("` + s + `") == 0; });x[` + strconv.Itoa(itemIndex) + `]`)
}

// GetContentOfItemWithClasses fetches the content of the element with the specified classes
func (cs *ChromeSession) GetContentOfItemWithClasses(classes string, itemIndex int) {
	cs.Write(`document.getElementsByClassName("` + classes + `")[` + strconv.Itoa(itemIndex) + `].innerHTML`)
}

// GetValueOfItemWithClasses returns the form value of the specified item
func (cs *ChromeSession) GetValueOfItemWithClasses(classes string, itemIndex int) {
	cs.Write(`document.getElementsByClassName("` + classes + `")[` + strconv.Itoa(itemIndex) + `].value`)
}

// GetContentOfItemWithSelector gets the content of an element with the specified selector
func (cs *ChromeSession) GetContentOfItemWithSelector(selector string) {
	cs.Write(`document.querySelector("` + selector + `").innerHTML()`)
}

// ClickItemWithClasses clicks on the first item it finds with the provided classes.
// Multiple classes are separated by spaces
func (cs *ChromeSession) ClickItemWithClasses(classes string, itemIndex int) {
	cs.Write(`document.getElementsByClassName("` + classes + `")[` + strconv.Itoa(itemIndex) + `].click()`)
}

// SetTextByID sets the text on the div with the specified id
func (cs *ChromeSession) SetTextByID(id string, text string) {
	cs.Write(`document.getElementById("` + id + `").innerHTML = "` + text + `"`)
}

// ClickItemWithID clicks an item with the specified id
func (cs *ChromeSession) ClickItemWithID(id string) {
	cs.Write(`document.getElementById("` + id + `").click()`)
}

// SetTextByClasses sets the text on the div with the specified id
func (cs *ChromeSession) SetTextByClasses(classes string, itemIndex int, text string) {
	cs.Write(`document.getElementsByClassName("` + classes + `")[` + strconv.Itoa(itemIndex) + `].innerHTML = "` + text + `"`)
}

// SetInputTextByClasses sets the input text for an input field
func (cs *ChromeSession) SetInputTextByClasses(classes string, itemIndex int, text string) {
	cs.Write(`document.getElementsByClassName("` + classes + `")[` + strconv.Itoa(itemIndex) + `].value = "` + text + `"`)
}

// NewBrowser starts a new chrome headless Session.
func NewBrowser(url string) (*ChromeSession, error) {
	var err error

	chromeSession := ChromeSession{}
	chromeSession.Output = make(chan string, 5000)

	// add url as last arg and create new Session
	args := append(Args, url)
	chromeSession.Session, err = interactive.NewSession(ChromePath, args)

	// map output and input channels for easy use
	chromeSession.Input = chromeSession.Session.Input

	go chromeSession.OutputSanitizer()

	firstLine := <-chromeSession.Output
	if !strings.Contains(firstLine, expectedFirstLine) {
		log.Println("ERROR: Unespected first line when initializing headless Chrome console:", firstLine)
	}

	return &chromeSession, err
}
