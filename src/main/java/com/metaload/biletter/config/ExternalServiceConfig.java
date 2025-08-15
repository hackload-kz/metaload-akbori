package com.metaload.biletter.config;

import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.stereotype.Component;

@Component
@ConfigurationProperties(prefix = "external")
public class ExternalServiceConfig {

    private HackloadConfig hackload = new HackloadConfig();
    private PaymentGatewayConfig paymentGateway = new PaymentGatewayConfig();

    public static class HackloadConfig {
        private String baseUrl = "http://localhost:8080";
        private String apiVersion = "v1";

        public String getBaseUrl() {
            return baseUrl;
        }

        public void setBaseUrl(String baseUrl) {
            this.baseUrl = baseUrl;
        }

        public String getApiVersion() {
            return apiVersion;
        }

        public void setApiVersion(String apiVersion) {
            this.apiVersion = apiVersion;
        }
    }

    public static class PaymentGatewayConfig {
        private String baseUrl = "https://gateway.hackload.com";
        private String apiVersion = "v1";
        private String teamSlug = "default-team";
        private String password = "default-password";

        public String getBaseUrl() {
            return baseUrl;
        }

        public void setBaseUrl(String baseUrl) {
            this.baseUrl = baseUrl;
        }

        public String getApiVersion() {
            return apiVersion;
        }

        public void setApiVersion(String apiVersion) {
            this.apiVersion = apiVersion;
        }

        public String getTeamSlug() {
            return teamSlug;
        }

        public void setTeamSlug(String teamSlug) {
            this.teamSlug = teamSlug;
        }

        public String getPassword() {
            return password;
        }

        public void setPassword(String password) {
            this.password = password;
        }
    }

    public HackloadConfig getHackload() {
        return hackload;
    }

    public void setHackload(HackloadConfig hackload) {
        this.hackload = hackload;
    }

    public PaymentGatewayConfig getPaymentGateway() {
        return paymentGateway;
    }

    public void setPaymentGateway(PaymentGatewayConfig paymentGateway) {
        this.paymentGateway = paymentGateway;
    }
}
