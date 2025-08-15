package com.metaload.biletter.dto;

import java.util.List;

public class ListBookingsResponseItem {

    private Long id;
    private Long eventId;
    private List<ListBookingsResponseItemSeat> seats;

    public ListBookingsResponseItem() {
    }

    public ListBookingsResponseItem(Long id, Long eventId) {
        this.id = id;
        this.eventId = eventId;
    }

    // Getters and Setters
    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public Long getEventId() {
        return eventId;
    }

    public void setEventId(Long eventId) {
        this.eventId = eventId;
    }

    public List<ListBookingsResponseItemSeat> getSeats() {
        return seats;
    }

    public void setSeats(List<ListBookingsResponseItemSeat> seats) {
        this.seats = seats;
    }
}
