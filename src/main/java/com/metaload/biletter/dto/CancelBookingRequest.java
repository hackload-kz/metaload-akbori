package com.metaload.biletter.dto;

import jakarta.validation.constraints.NotNull;

public class CancelBookingRequest {

    @NotNull(message = "Booking ID is required")
    private Long bookingId;

    public CancelBookingRequest() {
    }

    public CancelBookingRequest(Long bookingId) {
        this.bookingId = bookingId;
    }

    // Getters and Setters
    public Long getBookingId() {
        return bookingId;
    }

    public void setBookingId(Long bookingId) {
        this.bookingId = bookingId;
    }
}
