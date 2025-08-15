package com.metaload.biletter.dto;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;

public class CreateEventRequest {

    @NotBlank(message = "Title is required")
    private String title;

    @NotNull(message = "External flag is required")
    private Boolean external = false;

    // Getters and Setters
    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public Boolean getExternal() {
        return external;
    }

    public void setExternal(Boolean external) {
        this.external = external;
    }
}
