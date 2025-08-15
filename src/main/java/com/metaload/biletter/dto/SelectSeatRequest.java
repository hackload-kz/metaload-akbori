package com.metaload.biletter.dto;

import jakarta.validation.constraints.NotNull;

public class SelectSeatRequest {

    @NotNull(message = "Booking ID is required")
    private Long bookingId;

    @NotNull(message = "Seat ID is required")
    private Long seatId;

    public SelectSeatRequest() {
    }

    public SelectSeatRequest(Long bookingId, Long seatId) {
        this.bookingId = bookingId;
        this.seatId = seatId;
    }

    // Getters and Setters
    public Long getBookingId() {
        return bookingId;
    }

    public void setBookingId(Long bookingId) {
        this.bookingId = bookingId;
    }

    public Long getSeatId() {
        return seatId;
    }

    public void setSeatId(Long seatId) {
        this.seatId = seatId;
    }
}
