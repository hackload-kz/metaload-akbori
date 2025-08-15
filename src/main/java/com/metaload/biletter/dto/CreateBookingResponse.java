package com.metaload.biletter.dto;

public class CreateBookingResponse {

    private Long id;

    public CreateBookingResponse() {
    }

    public CreateBookingResponse(Long id) {
        this.id = id;
    }

    // Getters and Setters
    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }
}
