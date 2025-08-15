package com.metaload.biletter.dto.payment;

public class PaymentCancelRequest {

    private String teamSlug;
    private String token;
    private String paymentId;
    private String reason;
    private Receipt receipt;
    private Boolean force;
    private Object data;

    public PaymentCancelRequest() {
    }

    public static class Receipt {
        private String email;
        private String phone;
        private String taxation;

        // Getters and Setters
        public String getEmail() {
            return email;
        }

        public void setEmail(String email) {
            this.email = email;
        }

        public String getPhone() {
            return phone;
        }

        public void setPhone(String phone) {
            this.phone = phone;
        }

        public String getTaxation() {
            return taxation;
        }

        public void setTaxation(String taxation) {
            this.taxation = taxation;
        }
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

    public String getReason() {
        return reason;
    }

    public void setReason(String reason) {
        this.reason = reason;
    }

    public Receipt getReceipt() {
        return receipt;
    }

    public void setReceipt(Receipt receipt) {
        this.receipt = receipt;
    }

    public Boolean getForce() {
        return force;
    }

    public void setForce(Boolean force) {
        this.force = force;
    }

    public Object getData() {
        return data;
    }

    public void setData(Object data) {
        this.data = data;
    }
}
