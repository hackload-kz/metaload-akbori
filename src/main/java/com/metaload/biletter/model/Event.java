package com.metaload.biletter.model;

import jakarta.persistence.*;

import java.time.LocalDateTime;

@Entity
@Table(name = "events")
public class Event {
    @Id
    private Long id;
    private String title;
    private String description;
    private String type;
    private LocalDateTime datetimeStart;
    private String provider;

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public String getType() {
        return type;
    }

    public void setType(String type) {
        this.type = type;
    }

    public LocalDateTime getDatetimeStart() {
        return datetimeStart;
    }

    public void setDatetimeStart(LocalDateTime datetimeStart) {
        this.datetimeStart = datetimeStart;
    }

    public String getProvider() {
        return provider;
    }

    public void setProvider(String provider) {
        this.provider = provider;
    }
}
