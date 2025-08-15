package com.metaload.biletter.repository;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.PropertyNamingStrategy;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.metaload.biletter.model.Event;
import org.springframework.stereotype.Repository;

import java.io.File;
import java.io.IOException;
import java.time.LocalDate;
import java.util.ArrayList;
import java.util.List;

@Repository
public class EventRepository {
    private final List<Event> events;

    public EventRepository() {
        List<Event> events;
        try {
            ObjectMapper objectMapper = new ObjectMapper();
            objectMapper.setPropertyNamingStrategy(PropertyNamingStrategy.SNAKE_CASE);
            objectMapper.registerModule(new JavaTimeModule());
            String json = "requirements/events.json";
            events = objectMapper.readValue(new File(json), new TypeReference<>() {
            });
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
        this.events = events;
    }

    public Event findById(Long id) {
        return events.stream().filter(e -> e.getId().equals(id)).findFirst().orElse(null);
    }

    public List<Event> find(String query, LocalDate date) {
        if (query == null && date == null) {
            return events;
        }

        List<Event> filteredEvents = new ArrayList<>(events);
        if (query != null) {
            filteredEvents = filteredEvents.stream()
                    .filter(event -> event.getTitle().toLowerCase().contains(query.toLowerCase()))
                    .toList();
        }
        if (date != null) {
            filteredEvents = filteredEvents
                    .stream()
                    .filter(event -> event.getDatetimeStart().toLocalDate().isAfter(date))
                    .toList();
        }
        return filteredEvents;
    }
}
