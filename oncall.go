package main

import (
	"fmt"
	"log"
	"net/url"
)

const SCHEME = "https"
const SQUADCAST_HOST = "api.squadcast.com"
const ACCESS_TOKEN_PATH = "/v3/oauth/access-token"
const SCHEDULES_PATH = "/v3/schedules"
const ONCALL_PATH = "/v3/schedules/%s/on-call"

type OnCall struct {
	RefreshToken    string
	ScheduleName    string
	SlackWebhookURL string
	AccessToken     string
	ScheduleID      string
	OnCallShiftType string
	OnCallPeople    []string
}

func NewOnCall(params Params) *OnCall {
	return &OnCall{
		RefreshToken:    params.RefreshToken,
		ScheduleName:    params.ScheduleName,
		SlackWebhookURL: params.SlackWebhookURL,
	}
}

func (oc *OnCall) GetAccessToken() *OnCall {
	var resp AccessTokenResponse
	var sqerr SquadcastErrorResponse

	err := NewRequest().
		Get((&url.URL{
			Scheme: SCHEME,
			Host:   SQUADCAST_HOST,
			Path:   ACCESS_TOKEN_PATH,
		}).String()).
		SetHeader("X-Refresh-Token", oc.RefreshToken).
		With(&resp).
		WithFail(&sqerr).
		Do()

	if err != nil {
		log.Fatalf("%s : %s", err, sqerr.Meta.ErrorMessage)
	}

	oc.AccessToken = resp.Data.AccessToken
	return oc
}

func (oc *OnCall) GetScheduleID() *OnCall {
	var resp SchedulesResponse
	var sqerr SquadcastErrorResponse

	err := NewRequest().
		Get((&url.URL{
			Scheme: SCHEME,
			Host:   SQUADCAST_HOST,
			Path:   SCHEDULES_PATH,
		}).String()).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", oc.AccessToken)).
		With(&resp).
		WithFail(&sqerr).
		Do()
	if err != nil {
		log.Fatalf("%s : %s", err, sqerr.Meta.ErrorMessage)
	}

	for _, sch := range resp.Data {
		if sch.Name == oc.ScheduleName {
			oc.ScheduleID = sch.ID
			break
		}
	}
	if oc.ScheduleID == "" {
		log.Fatalf("Schedule of name: %s doesn't exist", oc.ScheduleName)
	}
	return oc
}

func (oc *OnCall) GetOnCallPeople() *OnCall {
	var resp OnCallResponse
	var sqerr SquadcastErrorResponse

	err := NewRequest().
		Get((&url.URL{
			Scheme: SCHEME,
			Host:   SQUADCAST_HOST,
			Path:   fmt.Sprintf(ONCALL_PATH, oc.ScheduleID),
		}).String()).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", oc.AccessToken)).
		With(&resp).
		WithFail(&sqerr).
		Do()
	if err != nil {
		log.Fatalf("%s : %s", err, sqerr.Meta.ErrorMessage)
	}

	oc.OnCallShiftType = resp.Data.ShiftType
	for _, user := range resp.Data.Users {
		oc.OnCallPeople = append(oc.OnCallPeople, fmt.Sprintf("%s %s", user.FirstName, user.LastName))
	}
	return oc
}

func (oc *OnCall) NotifySlack() {
	var resp string
	var slkerr string

	fields := make([]SlackWebhookAttachmentFields, 0)
	for _, user := range oc.OnCallPeople {
		fields = append(fields, SlackWebhookAttachmentFields{
			Title: user,
		})
	}

	body := SlackWebhookRequest{
		Attachments: []SlackWebhookAttachment{
			{
				Fallback: fmt.Sprintf("On-Call Update for Schedule: %s", oc.ScheduleName),
				Pretext:  fmt.Sprintf("People On-Call for Schedule: %s", oc.ScheduleName),
				Color:    "#00FF00",
				Fields:   fields,
			},
		},
	}

	err := NewRequest().
		Post(oc.SlackWebhookURL).
		Data(body).
		With(&resp).
		WithFail(&slkerr).
		Do()
	if err != nil {
		log.Fatalf("%s : %s", err, slkerr)
	}
}
