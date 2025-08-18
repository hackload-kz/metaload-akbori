package com.metaload.biletter.model.domainevents;

import java.time.LocalDateTime;

public class BookingCancelledEvent {

    private Long bookingId;
    private String orderId;
    private Long eventId;
    private Integer userId;
    private String previousStatus;
    private LocalDateTime cancelledAt;

    public BookingCancelledEvent() {
    }

    public BookingCancelledEvent(Long bookingId,
                                String orderId,
                                Long eventId,
                                Integer userId,
                                String previousStatus,
                                LocalDateTime cancelledAt) {
        this.bookingId = bookingId;
        this.orderId = orderId;
        this.eventId = eventId;
        this.userId = userId;
        this.previousStatus = previousStatus;
        this.cancelledAt = cancelledAt;
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

    public String getPreviousStatus() {
        return previousStatus;
    }

    public void setPreviousStatus(String previousStatus) {
        this.previousStatus = previousStatus;
    }

    public LocalDateTime getCancelledAt() {
        return cancelledAt;
    }

    public void setCancelledAt(LocalDateTime cancelledAt) {
        this.cancelledAt = cancelledAt;
    }

    @Override
    public String toString() {
        return "BookingCancelledEvent{" +
                "bookingId=" + bookingId +
                ", orderId='" + orderId + '\'' +
                ", eventId=" + eventId +
                ", userId=" + userId +
                ", previousStatus='" + previousStatus + '\'' +
                ", cancelledAt=" + cancelledAt +
                '}';
    }
}