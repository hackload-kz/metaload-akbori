package com.metaload.biletter.dto.payment;

public class PaymentInitRequest {

    private String teamSlug;
    private String token;
    private Long amount;
    private String orderId;
    private String currency;
    private String description;
    private String successURL;
    private String failURL;
    private String notificationURL;
    private Integer paymentExpiry;
    private String email;
    private String language;

    public PaymentInitRequest() {
    }

    // Getters and Setters
    public String getTeamSlug() {
        return teamSlug;
    }

    public void setTeamSlug(String teamSlug) {
        this.teamSlug = teamSlug;
    }

    public String getToken() {
        return token;
    }

    public void setToken(String token) {
        this.token = token;
    }

    public Long getAmount() {
        return amount;
    }

    public void setAmount(Long amount) {
        this.amount = amount;
    }

    public String getOrderId() {
        return orderId;
    }

    public void setOrderId(String orderId) {
        this.orderId = orderId;
    }

    public String getCurrency() {
        return currency;
    }

    public void setCurrency(String currency) {
        this.currency = currency;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public String getSuccessURL() {
        return successURL;
    }

    public void setSuccessURL(String successURL) {
        this.successURL = successURL;
    }

    public String getFailURL() {
        return failURL;
    }

    public void setFailURL(String failURL) {
        this.failURL = failURL;
    }

    public String getNotificationURL() {
        return notificationURL;
    }

    public void setNotificationURL(String notificationURL) {
        this.notificationURL = notificationURL;
    }

    public Integer getPaymentExpiry() {
        return paymentExpiry;
    }

    public void setPaymentExpiry(Integer paymentExpiry) {
        this.paymentExpiry = paymentExpiry;
    }

    public String getEmail() {
        return email;
    }

    public void setEmail(String email) {
        this.email = email;
    }

    public String getLanguage() {
        return language;
    }

    public void setLanguage(String language) {
        this.language = language;
    }
}
