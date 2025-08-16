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

    public List<Event> find(String query, LocalDate date, Integer page, Integer pageSize) {
        // Ограничиваем pageSize максимум 20 согласно API спецификации
        int maxPageSize = Math.min(pageSize != null ? pageSize : 20, 20);
        int pageNumber = page != null ? Math.max(page, 1) : 1;

        return eventRepository.find(query, date);
    }

    public Event findById(Long id) {
        return eventRepository.findById(id);
    }

}
