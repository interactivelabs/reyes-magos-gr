package dtos

import "reyes-magos-gr/store/models"

type VolunteerListItem struct {
	models.Volunteer
	GivenCount       int64
	HeldCount        int64
	LastActivityDate *string
}

type VolunteerLocationGroup struct {
	Location   string
	Volunteers []VolunteerListItem
}
