package telegoat

type config struct {
	telegramHost string
}

var cfg = config{
	telegramHost: "https://api.telegram.org/bot",
}
