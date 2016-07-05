package discord

type User struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
}

type Message struct {
	Id        string `json:"id"`
	ChannelId string `json:"channel_id"`
	Author    *User  `json:"author"`
	Content   string `json:"content"`
}