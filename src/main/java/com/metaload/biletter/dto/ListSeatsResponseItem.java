package com.metaload.biletter.dto;

import com.metaload.biletter.model.Seat.SeatStatus;

public class ListSeatsResponseItem {

    private Long id;
    private Long row;
    private Long number;
    private SeatStatus status;

    public ListSeatsResponseItem() {
    }

    public ListSeatsResponseItem(Long id, Long row, Long number, SeatStatus status) {
        this.id = id;
        this.row = row;
        this.number = number;
        this.status = status;
    }

    // Getters and Setters
    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public Long getRow() {
        return row;
    }

    public void setRow(Long row) {
        this.row = row;
    }

    public Long getNumber() {
        return number;
    }

    public void setNumber(Long number) {
        this.number = number;
    }

    public SeatStatus getStatus() {
        return status;
    }

    public void setStatus(SeatStatus status) {
        this.status = status;
    }
}
