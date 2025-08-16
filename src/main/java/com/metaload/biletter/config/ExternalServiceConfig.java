package com.metaload.biletter.config;

import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.stereotype.Component;

@Component
@ConfigurationProperties(prefix = "external")
public class ExternalServiceConfig {

    private HackloadConfig hackload = new HackloadConfig();

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

    public HackloadConfig getHackload() {
        return hackload;
    }

    public void setHackload(HackloadConfig hackload) {
        this.hackload = hackload;
    }

}
