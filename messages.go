package chaturbate

type connectMessage struct {
	User         string `json:"user"`
	Password     string `json:"password"`
	Room         string `json:"room"`
	RoomPassword string `json:"room_password"`
}

type joinRoomMessage struct {
	Room string `json:"room"`
}

type notify struct {
	Type               string   `json:"type"`
	Username           string   `json:"username"`
	IsMod              bool     `json:"is_mod"`
	SendTo             string   `json:"send_to"`
	ToUsername         string   `json:"to_username"`
	DontSendTo         string   `json:"dont_send_to"`
	FromUsername       string   `json:"from_username"`
	Amount             int      `json:"amount"`
	InFanclub          bool     `json:"in_fanclub"`
	HasTokens          bool     `json:"has_tokens"`
	TippedRecently     bool     `json:"tipped_recently"`
	TippedAlotRecently bool     `json:"tipped_alot_recently"`
	TippedTonsRecently bool     `json:"tipped_tons_recently"`
	Message            string   `json:"message"`
	History            bool     `json:"history"`
	Msg                []string `json:"msg"`
	Background         string   `json:"background"`
	Foreground         string   `json:"foreground"`
	Weight             string   `json:"weight"`
}

type roomMessage struct {
	User               string `json:"user"`
	Color              string `json:"c"`
	XSuccessful        bool   `json:"X-Successful"`
	InFanclub          bool   `json:"in_fanclub"`
	F                  string `json:"f"` // UNKNOWN
	Gender             string `json:"gender"`
	HasTokens          bool   `json:"has_tokens"`
	Message            string `json:"m"`
	TippedAlotRecently bool   `json:"tipped_alot_recently"`
	TippedTonsRecently bool   `json:"tipped_tons_recently"`
	TippedRecently     bool   `json:"tipped_recently"`
	IsMod              bool   `json:"is_mod"`
}
