// This code was designed to be ran in main

package calendargetter_test

/*

	cg := calendargetter.NewCalendarProvider(calendarID)

	tz, err := time.LoadLocation("Australia/Perth")
	if err != nil {
		log.Error(err)
	}

	st := time.Date(
		2019, 8, 2, 18, 00, 0, 0, tz)
	ed := time.Date(
		2019, 8, 2, 18, 30, 0, 0, tz)

	ev1 := models.Event{
		Summary:     "Test item 1",
		Start:       st.Format(time.RFC3339),
		End:         ed.Format(time.RFC3339),
		Timezone:    tz.String(),
		Location:    "303 Pinjarra Road",
		Description: "Boop!",
	}

	// Create
	ev1, err = cg.Create(ev1)
	if err != nil {
		log.Error(err)
	}

	// List
	list, err := cg.List()
	if err != nil {
		log.Error(err)
	}

	for _, e := range list {
		if e.EventID == ev1.EventID {
			e.Link = ""
			if reflect.DeepEqual(e, ev1) {
				log.Infof("\nGOOD!\n")
				log.Infof("\n%s", pretty.Sprint(e))
			} else {
				log.Errorf("\nBAD!\n")
				log.Errorf("\n%s", pretty.Compare(e, ev1))
			}
		}
	}

	// Get
	ev3, err := cg.Get(ev1)
	if err != nil {
		log.Error(err)
	}
	ev3.Link = ""
	if reflect.DeepEqual(ev1, ev3) {
		log.Infof("\nGOOD!\n")
		log.Infof("\n%s", pretty.Sprint(ev1))
	} else {
		log.Errorf("\nBAD!\n")
		log.Errorf("\n%s", pretty.Compare(ev1, ev3))
	}

	// Update
	newend := time.Date(
		2019, 8, 2, 21, 30, 0, 0, tz)
	ev1.End = newend.Format(time.RFC3339)
	ev1.Summary = "Updated"

	uev1, err := cg.Update(ev1)
	if err != nil {
		log.Error(err)
	}

	ev4, err := cg.Get(uev1)
	if err != nil {
		log.Error(err)
	}

	if reflect.DeepEqual(uev1, ev4) {
		log.Infof("\nGOOD!\n")
		log.Infof("\n%s", pretty.Sprint(ev1))
	} else {
		log.Errorf("\nBAD!\n")
		log.Errorf("\n%s", pretty.Compare(ev1, ev3))
	}

	err = cg.Delete(ev4)
	if err != nil {
		log.Error(err)
	}
	// Clean up
	list = []models.Event{}
	list, err = cg.List()
	if err != nil {
		log.Error(err)
	}
	total := len(list)
	var count int
	for _, e := range list {
		err := cg.Delete(e)
		if err != nil {
			log.Errorf("%s\t%+v\n", pretty.Sprint(e), err)
		} else {
			count++
		}
	}
	log.Infof("Deleted %d of %d\n", count, total)


*/
