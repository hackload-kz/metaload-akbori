package com.metaload.biletter.model.domainevents;

import java.time.LocalDateTime;

public class BookingCreatedEvent {

    private Long bookingId;
    private String orderId;
    private Long eventId;
    private Integer userId;
    private String status;
    private LocalDateTime createdAt;

    public BookingCreatedEvent() {
    }

    public BookingCreatedEvent(Long bookingId,
                               String orderId,
                               Long eventId,
                               Integer userId,
                               String status,
                               LocalDateTime createdAt) {
        this.bookingId = bookingId;
        this.orderId = orderId;
        this.eventId = eventId;
        this.userId = userId;
        this.status = status;
        this.createdAt = createdAt;
    }

    public Long getBookingId() {
        return bookingId;
    }

    public void setBookingId(Long bookingId) {
        this.bookingId = bookingId;
    }

    public String getOrderId() {
        return orderId;
    }

    public void setOrderId(String orderId) {
        this.orderId = orderId;
    }

    public Long getEventId() {
        return eventId;
    }

    public void setEventId(Long eventId) {
        this.eventId = eventId;
    }

    public Integer getUserId() {
        return userId;
    }

    public void setUserId(Integer userId) {
        this.userId = userId;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }

    public LocalDateTime getCreatedAt() {
        return createdAt;
    }

    public void setCreatedAt(LocalDateTime createdAt) {
        this.createdAt = createdAt;
    }

    @Override
    public String toString() {
        return "BookingCreatedEvent{" +
                "bookingId=" + bookingId +
                ", orderId='" + orderId + '\'' +
                ", eventId=" + eventId +
                ", userId=" + userId +
                ", status='" + status + '\'' +
                ", createdAt=" + createdAt +
                '}';
    }
}
