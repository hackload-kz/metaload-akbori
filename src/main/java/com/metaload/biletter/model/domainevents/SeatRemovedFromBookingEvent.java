package com.metaload.biletter.model.domainevents;

import java.time.LocalDateTime;

public class SeatRemovedFromBookingEvent {
    private Long bookingId;
    private Long seatId;
    private Long eventId;
    private String orderId;
    private Integer userId;
    private Integer rowNumber;
    private Integer seatNumber;
    private LocalDateTime timestamp;

    public SeatRemovedFromBookingEvent() {
    }

    public SeatRemovedFromBookingEvent(Long bookingId,
                                       Long seatId,
                                       Long eventId,
                                       String orderId,
                                       Integer userId,
                                       Integer rowNumber,
                                       Integer seatNumber,
                                       LocalDateTime timestamp) {
        this.bookingId = bookingId;
        this.seatId = seatId;
        this.eventId = eventId;
        this.orderId = orderId;
        this.userId = userId;
        this.rowNumber = rowNumber;
        this.seatNumber = seatNumber;
        this.timestamp = timestamp;
    }

    public Long getBookingId() {
        return bookingId;
    }

    public void setBookingId(Long bookingId) {
        this.bookingId = bookingId;
    }

    public Long getSeatId() {
        return seatId;
    }

    public void setSeatId(Long seatId) {
        this.seatId = seatId;
    }

    public Long getEventId() {
        return eventId;
    }

    public void setEventId(Long eventId) {
        this.eventId = eventId;
    }

    public String getOrderId() {
        return orderId;
    }

    public void setOrderId(String orderId) {
        this.orderId = orderId;
    }

    public Integer getUserId() {
        return userId;
    }

    public void setUserId(Integer userId) {
        this.userId = userId;
    }

    public Integer getRowNumber() {
        return rowNumber;
    }

    public void setRowNumber(Integer rowNumber) {
        this.rowNumber = rowNumber;
    }

    public Integer getSeatNumber() {
        return seatNumber;
    }

    public void setSeatNumber(Integer seatNumber) {
        this.seatNumber = seatNumber;
    }

    public LocalDateTime getTimestamp() {
        return timestamp;
    }

    public void setTimestamp(LocalDateTime timestamp) {
        this.timestamp = timestamp;
    }
}
