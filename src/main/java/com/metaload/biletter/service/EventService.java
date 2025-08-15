package com.metaload.biletter.service;

import com.metaload.biletter.dto.CreateEventRequest;
import com.metaload.biletter.model.Event;
import com.metaload.biletter.repository.EventRepository;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDate;
import java.util.List;

@Service
@Transactional
public class EventService {

    private final EventRepository eventRepository;

    public EventService(EventRepository eventRepository) {
        this.eventRepository = eventRepository;
    }

    public Event createEvent(CreateEventRequest request) {
        Event event = new Event();
        event.setTitle(request.getTitle());
        event.setExternal(request.getExternal());

        return eventRepository.save(event);
    }

    public List<Event> findEvents(String query, LocalDate date) {
        if (query == null && date == null) {
            return eventRepository.findAll();
        }

        return eventRepository.findEventsByQueryAndDate(query, date);
    }

    public Event findById(Long id) {
        return eventRepository.findById(id)
                .orElseThrow(() -> new RuntimeException("Event not found with id: " + id));
    }

    public List<Event> findByExternal(Boolean external) {
        return eventRepository.findByExternal(external);
    }
}
