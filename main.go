package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/EzraEllette/trie"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
)

func power(num int, degree int) int {
	if degree == 0 {
		return 1
	}
	product := num
	for i := 1; i < degree; i++ {
		product = product * num
	}

	return product
}

func stringToInt(digits string) int {
	var degree int = 0
	key := map[string]int{"0": 0, "1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9}
	var number int = 0
	for idx := range digits {
		number = number + (key[string(digits[idx])] * (power(10, degree)))
	}
	return number
}

func splitByTab(line string) (string, int) {
	var seenTab bool = false
	var word string = ""
	var number string = ""

	for idx := range line {
		char := string(line[idx])

		if char == "\t" {
			seenTab = true
			continue
		}

		if !seenTab {
			word = word + string(line[idx])
		} else {
			number = number + string(line[idx])
		}
	}
	return word, stringToInt(number)
}

func main() {

	file, err := os.Open("./data/1_1_all_alpha.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	dictionary := trie.New()
	reader := bufio.NewScanner(file)

	for reader.Scan() {
		line := reader.Text()
		line, value := splitByTab(line)
		dictionary.Insert(line, value)
	}

	if err := reader.Err(); err != nil {
		log.Fatal(err)
	}

	// Set logger
	l := log.New(log.Writer(), log.Prefix(), log.Flags())

	// Create astilectron
	a, err := astilectron.New(l, astilectron.Options{
		AppName:           "Autocomplete",
		AppIconDarwinPath: "public/icon.icns",
		BaseDirectoryPath: "./",
	})
	if err != nil {
		l.Fatal(fmt.Errorf("main: creating astilectron failed: %w", err))
	}
	defer a.Close()

	// Handle signals
	a.HandleSignals()

	// Start
	if err = a.Start(); err != nil {
		l.Fatal(fmt.Errorf("main: starting astilectron failed: %w", err))
	}

	// New window
	var w *astilectron.Window
	if w, err = a.NewWindow("public/index.html", &astilectron.WindowOptions{
		Center: astikit.BoolPtr(true),
		Height: astikit.IntPtr(700),
		Width:  astikit.IntPtr(700),
	}); err != nil {
		l.Fatal(fmt.Errorf("main: new window failed: %w", err))
	}

	// Create windows
	if err = w.Create(); err != nil {
		l.Fatal(fmt.Errorf("main: creating window failed: %w", err))
	}

	// This will listen to messages sent by Javascript
	w.OnMessage(func(m *astilectron.EventMessage) interface{} {
		// Unmarshal
		var s string
		m.Unmarshal(&s)

		// Process message
		term := dictionary.Autocomplete(s)

		return term
	})
	// Blocking pattern

	var m = a.NewMenu([]*astilectron.MenuItemOptions{
		{
			Label: astikit.StrPtr("Separator"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Label: astikit.StrPtr("Open Dev Tools"),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						w.OpenDevTools()
						return
					},
				},
				{
					Label: astikit.StrPtr("Close Dev Tools"),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						w.CloseDevTools()
						return
					},
				},
				{Type: astilectron.MenuItemTypeSeparator},
				{
					Label: astikit.StrPtr("Quit"),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						a.Close()
						return
					},
				},
			},
		},
	})

	m.Create()

	a.Wait()
}
