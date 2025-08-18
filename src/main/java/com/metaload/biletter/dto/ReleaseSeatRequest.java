package com.metaload.biletter.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import jakarta.validation.constraints.NotNull;

public class ReleaseSeatRequest {

    @NotNull(message = "Seat ID is required")
    @JsonProperty("seat_id")
    private Long seatId;

    public ReleaseSeatRequest() {
    }

    public ReleaseSeatRequest(Long seatId) {
        this.seatId = seatId;
    }

    // Getters and Setters
    public Long getSeatId() {
        return seatId;
    }

    public void setSeatId(Long seatId) {
        this.seatId = seatId;
    }
}
