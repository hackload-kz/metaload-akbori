package com.metaload.biletter.controller;

import com.metaload.biletter.dto.CreateEventRequest;
import com.metaload.biletter.dto.CreateEventResponse;
import com.metaload.biletter.dto.ListEventsResponseItem;
import com.metaload.biletter.model.Event;
import com.metaload.biletter.service.EventService;
import jakarta.validation.Valid;
import org.springframework.format.annotation.DateTimeFormat;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.time.LocalDate;
import java.util.List;
import java.util.stream.Collectors;

@RestController
@RequestMapping("/events")
public class EventController {

    private final EventService eventService;

    public EventController(EventService eventService) {
        this.eventService = eventService;
    }

    @PostMapping
    public ResponseEntity<CreateEventResponse> createEvent(@Valid @RequestBody CreateEventRequest request) {
        Event event = eventService.createEvent(request);
        return ResponseEntity.status(HttpStatus.CREATED)
                .body(new CreateEventResponse(event.getId()));
    }

    @GetMapping
    public ResponseEntity<List<ListEventsResponseItem>> listEvents(
            @RequestParam(required = false) String query,
            @RequestParam(required = false) @DateTimeFormat(iso = DateTimeFormat.ISO.DATE) LocalDate date) {

        List<Event> events = eventService.findEvents(query, date);
        List<ListEventsResponseItem> response = events.stream()
                .map(event -> {
                    ListEventsResponseItem item = new ListEventsResponseItem();
                    item.setId(event.getId());
                    item.setTitle(event.getTitle());
                    return item;
                })
                .collect(Collectors.toList());

        return ResponseEntity.ok(response);
    }
}
