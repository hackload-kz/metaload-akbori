package com.metaload.biletter.service;

import com.metaload.biletter.dto.ListEventsResponseItem;
import com.metaload.biletter.model.Event;
import com.metaload.biletter.repository.EventRepository;
import org.springframework.data.domain.PageRequest;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.util.List;

@Service
public class EventService {

    public static final Long MAIN_EVENT = 1L;

    private final EventRepository eventRepository;

    public EventService(EventRepository eventRepository) {
        this.eventRepository = eventRepository;
    }

    @Transactional(readOnly = true)
    public List<ListEventsResponseItem> find(String query, LocalDate date, Integer page, Integer pageSize) {
        // Ограничиваем pageSize максимум 20 согласно API спецификации
        int maxPageSize = Math.min(pageSize != null ? pageSize : 20, 20);
        int pageNumber = page != null ? Math.max(page, 1) : 1;

        LocalDateTime dt = date == null ? null : date.atStartOfDay();

        if (query == null && dt == null) {
            return eventRepository.find(PageRequest.of(pageNumber - 1, maxPageSize));
        } else if (query == null) {
            return eventRepository.find(dt, PageRequest.of(pageNumber - 1, maxPageSize));
        } else if (dt == null) {
            return eventRepository.find(query, PageRequest.of(pageNumber - 1, maxPageSize));
        }

        return eventRepository.find(query, dt, PageRequest.of(pageNumber - 1, maxPageSize));
    }

    @Transactional(readOnly = true)
    public Event findById(Long id) {
        return eventRepository.findById(id).orElseThrow();
    }

}
