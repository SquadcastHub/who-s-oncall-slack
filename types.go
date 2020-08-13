package main

type SquadcastErrorResponse struct {
	Meta SquadcastErrorDetails `json:"meta"`
}

type SquadcastErrorDetails struct {
	Status       int    `json:"status"`
	ErrorMessage string `json:"error_message"`
}

type AccessTokenResponse struct {
	Data AccessTokenDetails `json:"data"`
}

type AccessTokenDetails struct {
	AccessToken string `json:"access_token"`
}

type SchedulesResponse struct {
	Data []SchedulesDetails `json:"data"`
}

type SchedulesDetails struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type OnCallResponse struct {
	Data OnCallDetails `json:"data"`
}

type OnCallDetails struct {
	ShiftType string        `json:"shift_type"`
	Users     []UserDetails `json:"users"`
}

type UserDetails struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type SlackWebhookRequest struct {
	Text        string                   `json:"text,omitempty"`
	Channel     string                   `json:"channel,omitempty"`
	Username    string                   `json:"username,omitempty"`
	IconURL     string                   `json:"icon_url,omitempty"`
	IconEmoji   string                   `json:"icon_emoji,omitempty"`
	Attachments []SlackWebhookAttachment `json:"attachments,omitempty"`
}

type SlackWebhookAttachment struct {
	Fallback string                         `json:"fallback,omitempty"`
	Pretext  string                         `json:"pretext,omitempty"`
	Color    string                         `json:"color,omitempty"`
	Fields   []SlackWebhookAttachmentFields `json:"fields,omitempty"`
}

type SlackWebhookAttachmentFields struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool   `json:"short,omitempty"`
}
