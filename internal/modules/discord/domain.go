package discord

// Embed représente un message embed Discord
type Embed struct {
	Title       string       `json:"title,omitempty"`
	Description string       `json:"description,omitempty"`
	Color       int          `json:"color,omitempty"`
	Fields      []EmbedField `json:"fields,omitempty"`
	Timestamp   string       `json:"timestamp,omitempty"`
	Footer      EmbedFooter  `json:"footer,omitempty"`
	Thumbnail   EmbedImage   `json:"thumbnail,omitempty"`
}

// EmbedField représente un champ dans un embed
type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// EmbedFooter représente le pied de page d'un embed
type EmbedFooter struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url,omitempty"`
}

// EmbedImage représente une image ou une icône dans un embed
type EmbedImage struct {
	URL string `json:"url,omitempty"`
}

// DiscordEmbedMessage représente un message Discord avec un embed
type DiscordEmbedMessage struct {
	Content string  `json:"content,omitempty"`
	Embeds  []Embed `json:"embeds"`
}
