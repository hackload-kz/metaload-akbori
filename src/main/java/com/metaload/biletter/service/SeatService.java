package com.metaload.biletter.service;

import com.metaload.biletter.dto.ListSeatsResponseItem;
import com.metaload.biletter.model.Seat;
import com.metaload.biletter.model.Seat.SeatStatus;
import com.metaload.biletter.repository.SeatRepository;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.stream.Collectors;

@Service
@Transactional
public class SeatService {

    private final SeatRepository seatRepository;
    private final EventService eventService;

    public SeatService(SeatRepository seatRepository, EventService eventService) {
        this.seatRepository = seatRepository;
        this.eventService = eventService;
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

    public void createSeatsForEvent(Long eventId, int rows, int seatsPerRow) {
        // Проверяем, что событие существует
        eventService.findById(eventId);

        for (int row = 1; row <= rows; row++) {
            for (int seatNumber = 1; seatNumber <= seatsPerRow; seatNumber++) {
                Seat seat = new Seat();
                seat.setEvent(eventService.findById(eventId));
                seat.setRowNumber(row);
                seat.setSeatNumber(seatNumber);
                seat.setStatus(SeatStatus.FREE);

                seatRepository.save(seat);
            }
        }
    }

    private ListSeatsResponseItem mapToResponseItem(Seat seat) {
        ListSeatsResponseItem item = new ListSeatsResponseItem();
        item.setId(seat.getId());
        item.setRow(seat.getRowNumber().longValue());
        item.setNumber(seat.getSeatNumber().longValue());
        item.setStatus(seat.getStatus());
        return item;
    }
}
