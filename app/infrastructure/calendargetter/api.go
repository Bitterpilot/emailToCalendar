package calendargetter

import (
	"github.com/bitterpilot/emailToCalendar/models"
	"google.golang.org/api/calendar/v3"
)

// Create
func (p *CalendarProvider) Create(event models.Event) (models.Event, error) {
	e := convertToCalendar(event)

	id, err := p.service.Events.Insert(p.calID, e).Do()
	if err != nil {
		return models.Event{}, err
	}
	event.EventID = id.Id
	return event, nil
}

// List
func (p *CalendarProvider) List() ([]models.Event, error) {
	list, err := p.service.Events.List(p.calID).Do()
	if err != nil {
		return nil, err
	}

	var ret []models.Event
	for _, e := range list.Items {
		ret = append(ret, convertToModel(e))
	}
	return ret, nil
}

// Get
func (p *CalendarProvider) Get(e models.Event) (models.Event, error) {
	rsp, err := p.service.Events.Get(p.calID, e.EventID).Do()
	if err != nil {
		return models.Event{}, err
	}
	return convertToModel(rsp), nil
}

// Update
func (p *CalendarProvider) Update(new models.Event) (models.Event, error) {
	current, err := p.service.Events.Get(p.calID, new.EventID).Do()
	if err != nil {
		return models.Event{}, err
	}

	current.Summary = new.Summary
	current.Location = new.Location
	current.Description = new.Description
	current.Start.DateTime = new.Start
	current.End.DateTime = new.End

	rsp, err := p.service.Events.Update(p.calID, current.Id, current).Do()
	if err != nil {
		return models.Event{}, err
	}
	e := convertToModel(rsp)
	return e, nil
}

// Delete
func (p *CalendarProvider) Delete(e models.Event) error {
	return p.service.Events.Delete(p.calID, e.EventID).Do()
}

func convertToCalendar(new models.Event) *calendar.Event {
	ge := &calendar.Event{
		Summary:     new.Summary,
		Location:    new.Location,
		Description: new.Description,
		Start: &calendar.EventDateTime{
			DateTime: new.Start,
			TimeZone: new.Timezone,
		},
		End: &calendar.EventDateTime{
			DateTime: new.End,
			TimeZone: new.Timezone,
		},
	}
	return ge
}

func convertToModel(ge *calendar.Event) models.Event {
	new := models.Event{
		EventID:     ge.Id,
		Summary:     ge.Summary,
		Start:       ge.Start.DateTime,
		End:         ge.End.DateTime,
		Timezone:    ge.Start.TimeZone,
		Location:    ge.Location,
		Description: ge.Description,
		Link:        ge.HtmlLink,
	}
	return new
}
