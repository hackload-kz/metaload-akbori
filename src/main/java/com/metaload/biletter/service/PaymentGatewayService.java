package com.metaload.biletter.service;

import com.metaload.biletter.dto.payment.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClient;
import reactor.core.publisher.Mono;

import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.Map;
import java.util.TreeMap;

@Service
public class PaymentGatewayService {
    private static final Logger logger = LoggerFactory.getLogger(PaymentGatewayService.class);

    private final WebClient webClient;

    @Value("${app.url}")
    private String ticketServiceUrl;
    @Value("${app.paymentGateway.url}")
    private String paymentGatewayUrl;
    @Value("${app.paymentGateway.teamSlug}")
    private String teamSlug;
    @Value("${app.paymentGateway.password}")
    private String password;

    public PaymentGatewayService(WebClient.Builder webClientBuilder) {
        this.webClient = webClientBuilder
                .baseUrl(paymentGatewayUrl)
                .build();
    }

    public Mono<PaymentInitResponse> createPayment(PaymentInitRequest request) {
        logger.info("Creating payment for order: {}", request.getOrderId());

        // Generate authentication token
        String token = generateToken(request);
        request.setToken(token);

        return webClient.post()
                .uri("/PaymentInit/init")
                .bodyValue(request)
                .retrieve()
                .bodyToMono(PaymentInitResponse.class)
                .doOnSuccess(response -> logger.info("Payment created successfully: {}", response.getPaymentId()))
                .doOnError(error -> logger.error("Failed to create payment: {}", error.getMessage()));
    }

    /**
     * Создает PaymentInitRequest с базовыми параметрами
     */
    public PaymentInitRequest createPaymentRequest(String orderId, Long amount, String currency, String description,
            String email) {
        PaymentInitRequest request = new PaymentInitRequest();
        request.setTeamSlug(teamSlug);
        request.setOrderId(orderId);
        request.setAmount(amount);
        request.setCurrency(currency);
        request.setDescription(description);
        request.setEmail(email);
        request.setLanguage("ru");
        request.setPaymentExpiry(3600); // 1 час

        // Устанавливаем правильные URL для redirect и webhook
        request.setSuccessURL(ticketServiceUrl + "/payments/success?orderId=" + orderId);
        request.setFailURL(ticketServiceUrl + "/payments/fail?orderId=" + orderId);
        request.setNotificationURL(ticketServiceUrl + "/payments/notifications");

        return request;
    }

    public Mono<PaymentCheckResponse> checkPaymentStatus(String paymentId, String orderId) {
        logger.info("Checking payment status for paymentId: {} or orderId: {}", paymentId, orderId);

        PaymentCheckRequest request = new PaymentCheckRequest();
        request.setTeamSlug(teamSlug);

        if (paymentId != null) {
            request.setPaymentId(paymentId);
        } else if (orderId != null) {
            request.setOrderId(orderId);
        }

        // Generate token
        String token = generateToken(request);
        request.setToken(token);

        return webClient.post()
                .uri("/PaymentCheck/check")
                .bodyValue(request)
                .retrieve()
                .bodyToMono(PaymentCheckResponse.class)
                .doOnSuccess(response -> logger.info("Payment status checked successfully"))
                .doOnError(error -> logger.error("Failed to check payment status: {}", error.getMessage()));
    }

    public Mono<PaymentConfirmResponse> confirmPayment(String paymentId, Long amount) {
        logger.info("Confirming payment: {} with amount: {}", paymentId, amount);

        PaymentConfirmRequest request = new PaymentConfirmRequest();
        request.setTeamSlug(teamSlug);
        request.setPaymentId(paymentId);
        request.setAmount(amount);

        // Generate token
        String token = generateToken(request);
        request.setToken(token);

        return webClient.post()
                .uri("/PaymentConfirm/confirm")
                .bodyValue(request)
                .retrieve()
                .bodyToMono(PaymentConfirmResponse.class)
                .doOnSuccess(response -> logger.info("Payment confirmed successfully: {}", paymentId))
                .doOnError(error -> logger.error("Failed to confirm payment: {}", error.getMessage()));
    }

    public Mono<PaymentCancelResponse> cancelPayment(String paymentId, String reason) {
        logger.info("Cancelling payment: {} with reason: {}", paymentId, reason);

        PaymentCancelRequest request = new PaymentCancelRequest();
        request.setTeamSlug(teamSlug);
        request.setPaymentId(paymentId);
        request.setReason(reason);

        // Generate token
        String token = generateToken(request);
        request.setToken(token);

        return webClient.post()
                .uri("/PaymentCancel/cancel")
                .bodyValue(request)
                .retrieve()
                .bodyToMono(PaymentCancelResponse.class)
                .doOnSuccess(response -> logger.info("Payment cancelled successfully: {}", paymentId))
                .doOnError(error -> logger.error("Failed to cancel payment: {}", error.getMessage()));
    }

    private String generateToken(Object request) {
        try {
            // Extract the 5 required parameters in alphabetical order
            Map<String, String> params = new TreeMap<>();

            if (request instanceof PaymentInitRequest) {
                PaymentInitRequest req = (PaymentInitRequest) request;
                params.put("amount", String.valueOf(req.getAmount()));
                params.put("currency", req.getCurrency());
                params.put("orderId", req.getOrderId());
                params.put("teamSlug", req.getTeamSlug());
            } else if (request instanceof PaymentCheckRequest) {
                PaymentCheckRequest req = (PaymentCheckRequest) request;
                if (req.getPaymentId() != null) {
                    params.put("paymentId", req.getPaymentId());
                }
                if (req.getOrderId() != null) {
                    params.put("orderId", req.getOrderId());
                }
                params.put("teamSlug", req.getTeamSlug());
            } else if (request instanceof PaymentConfirmRequest) {
                PaymentConfirmRequest req = (PaymentConfirmRequest) request;
                params.put("amount", String.valueOf(req.getAmount()));
                params.put("paymentId", req.getPaymentId());
                params.put("teamSlug", req.getTeamSlug());
            } else if (request instanceof PaymentCancelRequest) {
                PaymentCancelRequest req = (PaymentCancelRequest) request;
                params.put("paymentId", req.getPaymentId());
                params.put("teamSlug", req.getTeamSlug());
            }

            // Concatenate values in alphabetical order
            StringBuilder tokenString = new StringBuilder();
            for (String value : params.values()) {
                tokenString.append(value);
            }

            // Add password
            tokenString.append(password);

            // Generate SHA-256 hash
            MessageDigest digest = MessageDigest.getInstance("SHA-256");
            byte[] hash = digest.digest(tokenString.toString().getBytes());

            StringBuilder hexString = new StringBuilder();
            for (byte b : hash) {
                String hex = Integer.toHexString(0xff & b);
                if (hex.length() == 1) {
                    hexString.append('0');
                }
                hexString.append(hex);
            }

            return hexString.toString();

        } catch (NoSuchAlgorithmException e) {
            logger.error("Failed to generate token: {}", e.getMessage());
            throw new RuntimeException("Failed to generate authentication token", e);
        }
    }
}
