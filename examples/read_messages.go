package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/sextech/chaturbate"
	"github.com/sextech/chaturbate/option"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <room>\n", os.Args[0])
		return
	}

	chat, err := chaturbate.NewChat(context.Background(), os.Args[1], option.WithoutAuthentication())

	if err != nil {
		log.Fatal(err)
	}

	chat.OnMessage = func(author, message string) {
		fmt.Printf("[>] %s: %s\n", author, message)
	}

	chat.OnTip = func(author string, amount int) {
		fmt.Printf("[!] %s tipped %d\n", author, amount)
	}

	chat.OnNotice = func(notices []string) {
		for _, notice := range notices {
			fmt.Printf("[*] %s\n", notice)
		}
	}

	log.Fatal(chat.Connect())
}
