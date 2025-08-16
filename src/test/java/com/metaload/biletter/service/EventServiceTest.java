package com.metaload.biletter.service;

import com.metaload.biletter.model.Event;
import com.metaload.biletter.repository.EventRepository;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.Mock;

import java.time.LocalDate;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class EventServiceTest {

    private EventRepository eventRepository = new EventRepository();

    @Test
    public void find() {
        EventService eventService = new EventService(eventRepository);
        List<Event> t = eventService.find("t", LocalDate.now(), 1, 20);
        assertNotNull(t);
    }
}