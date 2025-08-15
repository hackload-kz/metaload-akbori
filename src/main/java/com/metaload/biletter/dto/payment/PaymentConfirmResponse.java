package com.metaload.biletter.dto.payment;

public class PaymentConfirmResponse {

    private Boolean success;
    private String paymentId;
    private String orderId;
    private String status;
    private Long authorizedAmount;
    private Long confirmedAmount;
    private Long remainingAmount;
    private String currency;
    private String confirmedAt;
    private BankDetails bankDetails;
    private Fees fees;
    private Settlement settlement;

    public PaymentConfirmResponse() {
    }

    public static class BankDetails {
        private String bankTransactionId;
        private String authorizationCode;
        private String rrn;
        private String responseCode;
        private String responseMessage;

        // Getters and Setters
        public String getBankTransactionId() {
            return bankTransactionId;
        }

        public void setBankTransactionId(String bankTransactionId) {
            this.bankTransactionId = bankTransactionId;
        }

        public String getAuthorizationCode() {
            return authorizationCode;
        }

        public void setAuthorizationCode(String authorizationCode) {
            this.authorizationCode = authorizationCode;
        }

        public String getRrn() {
            return rrn;
        }

        public void setRrn(String rrn) {
            this.rrn = rrn;
        }

        public String getResponseCode() {
            return responseCode;
        }

        public void setResponseCode(String responseCode) {
            this.responseCode = responseCode;
        }

        public String getResponseMessage() {
            return responseMessage;
        }

        public void setResponseMessage(String responseMessage) {
            this.responseMessage = responseMessage;
        }
    }

    public static class Fees {
        private Long processingFee;
        private Long totalFees;
        private String feeCurrency;

        // Getters and Setters
        public Long getProcessingFee() {
            return processingFee;
        }

        public void setProcessingFee(Long processingFee) {
            this.processingFee = processingFee;
        }

        public Long getTotalFees() {
            return totalFees;
        }

        public void setTotalFees(Long totalFees) {
            this.totalFees = totalFees;
        }

        public String getFeeCurrency() {
            return feeCurrency;
        }

        public void setFeeCurrency(String feeCurrency) {
            this.feeCurrency = feeCurrency;
        }
    }

    public static class Settlement {
        private String settlementDate;
        private Long settlementAmount;
        private String settlementCurrency;

        // Getters and Setters
        public String getSettlementDate() {
            return settlementDate;
        }

        public void setSettlementDate(String settlementDate) {
            this.settlementDate = settlementDate;
        }

        public Long getSettlementAmount() {
            return settlementAmount;
        }

        public void setSettlementAmount(Long settlementAmount) {
            this.settlementAmount = settlementAmount;
        }

        public String getSettlementCurrency() {
            return settlementCurrency;
        }

        public void setSettlementCurrency(String settlementCurrency) {
            this.settlementCurrency = settlementCurrency;
        }
    }

    // Getters and Setters
    public Boolean getSuccess() {
        return success;
    }

    public void setSuccess(Boolean success) {
        this.success = success;
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

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }

    public Long getAuthorizedAmount() {
        return authorizedAmount;
    }

    public void setAuthorizedAmount(Long authorizedAmount) {
        this.authorizedAmount = authorizedAmount;
    }

    public Long getConfirmedAmount() {
        return confirmedAmount;
    }

    public void setConfirmedAmount(Long confirmedAmount) {
        this.confirmedAmount = confirmedAmount;
    }

    public Long getRemainingAmount() {
        return remainingAmount;
    }

    public void setRemainingAmount(Long remainingAmount) {
        this.remainingAmount = remainingAmount;
    }

    public String getCurrency() {
        return currency;
    }

    public void setCurrency(String currency) {
        this.currency = currency;
    }

    public String getConfirmedAt() {
        return confirmedAt;
    }

    public void setConfirmedAt(String confirmedAt) {
        this.confirmedAt = confirmedAt;
    }

    public BankDetails getBankDetails() {
        return bankDetails;
    }

    public void setBankDetails(BankDetails bankDetails) {
        this.bankDetails = bankDetails;
    }

    public Fees getFees() {
        return fees;
    }

    public void setFees(Fees fees) {
        this.fees = fees;
    }

    public Settlement getSettlement() {
        return settlement;
    }

    public void setSettlement(Settlement settlement) {
        this.settlement = settlement;
    }
}
