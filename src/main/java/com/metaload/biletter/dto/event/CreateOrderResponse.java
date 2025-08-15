package com.metaload.biletter.dto.event;

public class CreateOrderResponse {

    private String orderId;

    public CreateOrderResponse() {
    }

    public CreateOrderResponse(String orderId) {
        this.orderId = orderId;
    }

    public String getOrderId() {
        return orderId;
    }

    public void setOrderId(String orderId) {
        this.orderId = orderId;
    }
}
