package com.metaload.biletter.controller;

import com.metaload.biletter.dto.ListSeatsResponseItem;
import com.metaload.biletter.dto.ReleaseSeatRequest;
import com.metaload.biletter.dto.SelectSeatRequest;
import com.metaload.biletter.model.Seat.SeatStatus;
import com.metaload.biletter.service.BookingService;
import com.metaload.biletter.service.SeatService;
import jakarta.validation.Valid;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/seats")
public class SeatController {

    private final SeatService seatService;
    private final BookingService bookingService;

    public SeatController(SeatService seatService, BookingService bookingService) {
        this.seatService = seatService;
        this.bookingService = bookingService;
    }

    @GetMapping
    public ResponseEntity<List<ListSeatsResponseItem>> listSeats(
            @RequestParam Long event_id,
            @RequestParam(required = false) Integer page,
            @RequestParam(required = false) Integer pageSize,
            @RequestParam(required = false) Integer row,
            @RequestParam(required = false) SeatStatus status) {

        List<ListSeatsResponseItem> seats = seatService.getSeats(event_id, page, pageSize, row, status);
        return ResponseEntity.ok(seats);
    }

    @PatchMapping("/select")
    public ResponseEntity<Void> selectSeat(@Valid @RequestBody SelectSeatRequest request) {
        try {
            bookingService.selectSeat(request.getBookingId(), request.getSeatId());
            return ResponseEntity.ok().build();
        } catch (RuntimeException e) {
            return ResponseEntity.status(HttpStatus.INSUFFICIENT_SPACE_ON_RESOURCE).build();
        }
    }

    @PatchMapping("/release")
    public ResponseEntity<Void> releaseSeat(@Valid @RequestBody ReleaseSeatRequest request) {
        try {
            bookingService.releaseSeat(request.getSeatId());
            return ResponseEntity.ok().build();
        } catch (RuntimeException e) {
            return ResponseEntity.status(HttpStatus.INSUFFICIENT_SPACE_ON_RESOURCE).build();
        }
    }
}
