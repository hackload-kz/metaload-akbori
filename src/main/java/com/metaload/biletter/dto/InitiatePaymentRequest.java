package com.metaload.biletter.dto;

import jakarta.validation.constraints.NotNull;

public class InitiatePaymentRequest {

    @NotNull(message = "Booking ID is required")
    private Long bookingId;

    public InitiatePaymentRequest() {
    }

    public InitiatePaymentRequest(Long bookingId) {
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
