package com.metaload.biletter.dto;

public class CreateEventResponse {

    private Long id;

    public CreateEventResponse() {
    }

    public CreateEventResponse(Long id) {
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
