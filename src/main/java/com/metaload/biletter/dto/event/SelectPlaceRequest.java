package com.metaload.biletter.dto.event;

public class SelectPlaceRequest {

    private String orderId;

    public SelectPlaceRequest() {
    }

    public SelectPlaceRequest(String orderId) {
        this.orderId = orderId;
    }

    public String getOrderId() {
        return orderId;
    }

    public void setOrderId(String orderId) {
        this.orderId = orderId;
    }
}
