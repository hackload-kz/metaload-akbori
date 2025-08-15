package com.metaload.biletter.controller;

import com.metaload.biletter.dto.*;
import com.metaload.biletter.service.BookingService;
import jakarta.validation.Valid;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/bookings")
public class BookingController {

    private final BookingService bookingService;

    public BookingController(BookingService bookingService) {
        this.bookingService = bookingService;
    }

    @PostMapping
    public ResponseEntity<CreateBookingResponse> createBooking(@Valid @RequestBody CreateBookingRequest request) {
        var booking = bookingService.createBooking(request);
        return ResponseEntity.status(HttpStatus.CREATED)
                .body(new CreateBookingResponse(booking.getId()));
    }

    @GetMapping
    public ResponseEntity<List<ListBookingsResponseItem>> listBookings() {
        List<ListBookingsResponseItem> bookings = bookingService.getAllBookings();
        return ResponseEntity.ok(bookings);
    }

    @PatchMapping("/initiatePayment")
    public ResponseEntity<Void> initiatePayment(@Valid @RequestBody InitiatePaymentRequest request) {
        bookingService.initiatePayment(request.getBookingId());
        return ResponseEntity.ok().build();
    }

    @PatchMapping("/cancel")
    public ResponseEntity<Void> cancelBooking(@Valid @RequestBody CancelBookingRequest request) {
        bookingService.cancelBooking(request.getBookingId());
        return ResponseEntity.ok().build();
    }
}
