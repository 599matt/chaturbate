# Chaturbate

[![GoDoc](https://godoc.org/github.com/sextech/chaturbate?status.svg)](https://godoc.org/github.com/sextech/chaturbate)

Non-official Chaturbate API **UNDER DEVELOPMENT**

### Chat

You can connect to room chat by using `NewChat` method.

Example :

```go
chat, err := chaturbate.NewChat(context.Background(), "ROOM_NAME", option.WithoutAuthentication())

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
```

Available callbacks :

- `OnMessage func(author, message string)`
- `OnTip func(author string, amount int)`
- `OnNotice func(notices []string)`
- `OnUserEntry func(user string)`
- `OnUserLeave func(user string)`
- `OnMute func()`