package com.metaload.biletter.dto.payment;

public class PaymentCheckRequest {

    private String teamSlug;
    private String token;
    private String paymentId;
    private String orderId;
    private Boolean includeTransactions;
    private Boolean includeCardDetails;
    private Boolean includeCustomerInfo;
    private Boolean includeReceipt;
    private String language;

    public PaymentCheckRequest() {
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

    public String getPaymentId() {
        return paymentId;
    }

    public void setPaymentId(String paymentId) {
        this.paymentId = paymentId;
    }

    public String getOrderId() {
        return orderId;
    }

    public void setOrderId(String orderId) {
        this.orderId = orderId;
    }

    public Boolean getIncludeTransactions() {
        return includeTransactions;
    }

    public void setIncludeTransactions(Boolean includeTransactions) {
        this.includeTransactions = includeTransactions;
    }

    public Boolean getIncludeCardDetails() {
        return includeCardDetails;
    }

    public void setIncludeCardDetails(Boolean includeCardDetails) {
        this.includeCardDetails = includeCardDetails;
    }

    public Boolean getIncludeCustomerInfo() {
        return includeCustomerInfo;
    }

    public void setIncludeCustomerInfo(Boolean includeCustomerInfo) {
        this.includeCustomerInfo = includeCustomerInfo;
    }

    public Boolean getIncludeReceipt() {
        return includeReceipt;
    }

    public void setIncludeReceipt(Boolean includeReceipt) {
        this.includeReceipt = includeReceipt;
    }

    public String getLanguage() {
        return language;
    }

    public void setLanguage(String language) {
        this.language = language;
    }
}
