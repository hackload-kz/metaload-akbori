package com.metaload.biletter.dto;

public class ListBookingsResponseItemSeat {

    private Long id;

    public ListBookingsResponseItemSeat() {
    }

    public ListBookingsResponseItemSeat(Long id) {
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
