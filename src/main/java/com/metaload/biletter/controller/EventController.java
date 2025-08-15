package com.metaload.biletter.controller;

import com.metaload.biletter.dto.ListEventsResponseItem;
import com.metaload.biletter.model.Event;
import com.metaload.biletter.service.EventService;
import org.springframework.format.annotation.DateTimeFormat;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.time.LocalDate;
import java.util.List;

@RestController
@RequestMapping("/events")
public class EventController {

    private final EventService eventService;

    public EventController(EventService eventService) {
        this.eventService = eventService;
    }

    @GetMapping
    public ResponseEntity<List<ListEventsResponseItem>> listEvents(
            @RequestParam(required = false) String query,
            @RequestParam(required = false) @DateTimeFormat(iso = DateTimeFormat.ISO.DATE) LocalDate date) {

        List<Event> events = eventService.find(query, date);
        List<ListEventsResponseItem> response = events.stream()
                .map(event -> new ListEventsResponseItem(event.getId(), event.getTitle()))
                .toList();

        return ResponseEntity.ok(response);
    }
}
