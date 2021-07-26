package main

import (
	"learn-go-with-tests/mocking"
	"os"
)

const spanish = "Spanish"
const french = "French"
const spanishHelloPrefix = "Hola, "
const frenchHelloPrefix = "Bonjour, "
const englishHelloPrefix = "Hello, "

// Hello ...
func Hello(name string, language string) string {
	if name == "" {
		name = "world"
	}

	return greetingPrefix(language) + name
}

func greetingPrefix(language string) (prefix string) {
	switch language {
	case french:
		prefix = frenchHelloPrefix
	case spanish:
		prefix = spanishHelloPrefix
	default:
		prefix = englishHelloPrefix
	}
	return
}

func main() {
	// log.Fatal(http.ListenAndServe(":5000", http.HandlerFunc(di.MyGreetHandler)))
	sleeper := &mocking.DefaultSleeper{}
	mocking.Countdown(os.Stdout, sleeper)
}
