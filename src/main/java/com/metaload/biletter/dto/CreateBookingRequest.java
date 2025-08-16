package com.metaload.biletter.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import jakarta.validation.constraints.NotNull;

public class CreateBookingRequest {

    @NotNull(message = "Event ID is required")
    @JsonProperty("event_id")
    private Long eventId;

    public CreateBookingRequest() {
    }

    public CreateBookingRequest(Long eventId) {
        this.eventId = eventId;
    }

    // Getters and Setters
    public Long getEventId() {
        return eventId;
    }

    public void setEventId(Long eventId) {
        this.eventId = eventId;
    }
}
