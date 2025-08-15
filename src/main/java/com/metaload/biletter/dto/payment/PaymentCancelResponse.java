package com.metaload.biletter.dto.payment;

public class PaymentCancelResponse {

    private Boolean success;
    private String paymentId;
    private String orderId;
    private String status;
    private String cancellationType;
    private Long originalAmount;
    private Long cancelledAmount;
    private Long remainingAmount;
    private String currency;
    private String cancelledAt;
    private BankDetails bankDetails;
    private Refund refund;
    private Details details;

    public PaymentCancelResponse() {
    }

    public static class BankDetails {
        private String bankTransactionId;
        private String originalAuthorizationCode;
        private String cancellationAuthorizationCode;
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

        public String getOriginalAuthorizationCode() {
            return originalAuthorizationCode;
        }

        public void setOriginalAuthorizationCode(String originalAuthorizationCode) {
            this.originalAuthorizationCode = originalAuthorizationCode;
        }

        public String getCancellationAuthorizationCode() {
            return cancellationAuthorizationCode;
        }

        public void setCancellationAuthorizationCode(String cancellationAuthorizationCode) {
            this.cancellationAuthorizationCode = cancellationAuthorizationCode;
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

    public static class Refund {
        private String refundId;
        private String refundStatus;
        private String expectedProcessingTime;
        private String refundMethod;
        private CardInfo cardInfo;

        // Getters and Setters
        public String getRefundId() {
            return refundId;
        }

        public void setRefundId(String refundId) {
            this.refundId = refundId;
        }

        public String getRefundStatus() {
            return refundStatus;
        }

        public void setRefundStatus(String refundStatus) {
            this.refundStatus = refundStatus;
        }

        public String getExpectedProcessingTime() {
            return expectedProcessingTime;
        }

        public void setExpectedProcessingTime(String expectedProcessingTime) {
            this.expectedProcessingTime = expectedProcessingTime;
        }

        public String getRefundMethod() {
            return refundMethod;
        }

        public void setRefundMethod(String refundMethod) {
            this.refundMethod = refundMethod;
        }

        public CardInfo getCardInfo() {
            return cardInfo;
        }

        public void setCardInfo(CardInfo cardInfo) {
            this.cardInfo = cardInfo;
        }
    }

    public static class CardInfo {
        private String cardMask;
        private String cardType;
        private String issuingBank;

        // Getters and Setters
        public String getCardMask() {
            return cardMask;
        }

        public void setCardMask(String cardMask) {
            this.cardMask = cardMask;
        }

        public String getCardType() {
            return cardType;
        }

        public void setCardType(String cardType) {
            this.cardType = cardType;
        }

        public String getIssuingBank() {
            return issuingBank;
        }

        public void setIssuingBank(String issuingBank) {
            this.issuingBank = issuingBank;
        }
    }

    public static class Details {
        private String reason;
        private Boolean wasForced;
        private String processingDuration;

        // Getters and Setters
        public String getReason() {
            return reason;
        }

        public void setReason(String reason) {
            this.reason = reason;
        }

        public Boolean getWasForced() {
            return wasForced;
        }

        public void setWasForced(Boolean wasForced) {
            this.wasForced = wasForced;
        }

        public String getProcessingDuration() {
            return processingDuration;
        }

        public void setProcessingDuration(String processingDuration) {
            this.processingDuration = processingDuration;
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

    public String getCancellationType() {
        return cancellationType;
    }

    public void setCancellationType(String cancellationType) {
        this.cancellationType = cancellationType;
    }

    public Long getOriginalAmount() {
        return originalAmount;
    }

    public void setOriginalAmount(Long originalAmount) {
        this.originalAmount = originalAmount;
    }

    public Long getCancelledAmount() {
        return cancelledAmount;
    }

    public void setCancelledAmount(Long cancelledAmount) {
        this.cancelledAmount = cancelledAmount;
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

    public String getCancelledAt() {
        return cancelledAt;
    }

    public void setCancelledAt(String cancelledAt) {
        this.cancelledAt = cancelledAt;
    }

    public BankDetails getBankDetails() {
        return bankDetails;
    }

    public void setBankDetails(BankDetails bankDetails) {
        this.bankDetails = bankDetails;
    }

    public Refund getRefund() {
        return refund;
    }

    public void setRefund(Refund refund) {
        this.refund = refund;
    }

    public Details getDetails() {
        return details;
    }

    public void setDetails(Details details) {
        this.details = details;
    }
}
