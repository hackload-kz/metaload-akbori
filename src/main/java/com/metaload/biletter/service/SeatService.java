package com.metaload.biletter.service;

import com.metaload.biletter.dto.ListSeatsResponseItem;
import com.metaload.biletter.dto.event.Place;
import com.metaload.biletter.model.Event;
import com.metaload.biletter.model.Seat;
import com.metaload.biletter.model.Seat.SeatStatus;
import com.metaload.biletter.repository.SeatRepository;
import org.apache.commons.lang3.BooleanUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.scheduling.annotation.Async;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import reactor.core.publisher.Mono;

import java.math.BigDecimal;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.stream.Collectors;

@Service
@Transactional
public class SeatService {

    private static final Logger log = LoggerFactory.getLogger(SeatService.class);
    private final SeatRepository seatRepository;
    private final EventService eventService;
    private final EventProviderService eventProviderService;

    public SeatService(SeatRepository seatRepository, EventService eventService, EventProviderService eventProviderService) {
        this.seatRepository = seatRepository;
        this.eventService = eventService;
        this.eventProviderService = eventProviderService;
    }

    public List<ListSeatsResponseItem> getSeats(Long eventId, Integer page, Integer pageSize,
                                                Integer row, SeatStatus status) {
        // Проверяем, что событие существует
        eventService.findById(eventId);

        // Создаем пагинацию
        Pageable pageable = PageRequest.of(
                page != null ? page - 1 : 0,
                pageSize != null ? pageSize : 20);

        // Получаем места с фильтрами
        Page<Seat> seatsPage = seatRepository.findByEventIdAndFilters(eventId, row, status, pageable);

        return seatsPage.getContent().stream()
                .map(this::mapToResponseItem)
                .collect(Collectors.toList());
    }

    public Seat findById(Long id) {
        return seatRepository.findById(id)
                .orElseThrow(() -> new RuntimeException("Seat not found with id: " + id));
    }

    public List<Seat> getFreeSeatsByEvent(Long eventId) {
        return seatRepository.findByEventIdAndStatus(eventId, SeatStatus.FREE);
    }

    public void createSeatsForEvent(Long eventId, int rows, int seatsPerRow, BigDecimal defaultPrice) {
        // Проверяем, что событие существует
        eventService.findById(eventId);

        for (int row = 1; row <= rows; row++) {
            for (int seatNumber = 1; seatNumber <= seatsPerRow; seatNumber++) {
                Seat seat = new Seat();
                Event event = eventService.findById(eventId);
                seat.setEvent(event);
                seat.setRowNumber(row);
                seat.setSeatNumber(seatNumber);
                seat.setStatus(SeatStatus.FREE);
                seat.setPrice(defaultPrice);

                seatRepository.save(seat);
            }
        }
    }

    public void generateSeatsForSingleEvent(Long eventId) {
        List<Seat> buffer = new ArrayList<>(1000);
        int seatIndex = 0;

        Event event = new Event();
        event.setId(eventId);

        for (int row = 1; row <= 100; row++) {
            for (int seatNum = 1; seatNum <= 1000; seatNum++) {
                seatIndex++;
                Seat seat = new Seat();
                seat.setEvent(event);
                seat.setRowNumber(row);
                seat.setSeatNumber(seatNum);
                seat.setStatus(SeatStatus.FREE);
                seat.setPrice(getPrice(seatIndex));

                buffer.add(seat);

                if (buffer.size() >= 1000) {
                    seatRepository.saveAll(buffer);
                    seatRepository.flush();
                    buffer.clear();
                }
            }
        }
        if (!buffer.isEmpty()) {
            seatRepository.saveAll(buffer);
            seatRepository.flush();
        }
    }

    private BigDecimal getPrice(int index) {
        if (index <= 10_000) return BigDecimal.valueOf(40_000);
        else if (index <= 25_000) return BigDecimal.valueOf(80_000);
        else if (index <= 45_000) return BigDecimal.valueOf(120_000);
        else if (index <= 70_000) return BigDecimal.valueOf(160_000);
        else return BigDecimal.valueOf(200_000);
    }

    private ListSeatsResponseItem mapToResponseItem(Seat seat) {
        ListSeatsResponseItem item = new ListSeatsResponseItem();
        item.setId(seat.getId());
        item.setRow(seat.getRowNumber().longValue());
        item.setNumber(seat.getSeatNumber().longValue());
        item.setStatus(seat.getStatus());
        item.setPrice(seat.getPrice());
        return item;
    }
}
