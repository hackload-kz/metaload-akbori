package com.metaload.biletter.service;

import com.metaload.biletter.config.ExternalServiceConfig;
import com.metaload.biletter.dto.hackload.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClient;
import reactor.core.publisher.Mono;

@Service
public class HackloadService {

    private static final Logger logger = LoggerFactory.getLogger(HackloadService.class);

    private final WebClient webClient;
    private final ExternalServiceConfig config;

    public HackloadService(WebClient.Builder webClientBuilder, ExternalServiceConfig config) {
        this.config = config;
        this.webClient = webClientBuilder
                .baseUrl(config.getHackload().getBaseUrl())
                .build();
    }

    public Mono<CreateOrderResponse> createOrder() {
        logger.info("Creating order in Hackload service");

        return webClient.post()
                .uri("/api/partners/{version}/orders", config.getHackload().getApiVersion())
                .retrieve()
                .bodyToMono(CreateOrderResponse.class)
                .doOnSuccess(response -> logger.info("Order created successfully: {}", response.getOrderId()))
                .doOnError(error -> logger.error("Failed to create order: {}", error.getMessage()));
    }

    public Mono<OrderDetails> getOrder(String orderId) {
        logger.info("Getting order details for: {}", orderId);

        return webClient.get()
                .uri("/api/partners/{version}/orders/{orderId}", config.getHackload().getApiVersion(), orderId)
                .retrieve()
                .bodyToMono(OrderDetails.class)
                .doOnSuccess(response -> logger.info("Order details retrieved: {}", response.getStatus()))
                .doOnError(error -> logger.error("Failed to get order details: {}", error.getMessage()));
    }

    public Mono<Void> submitOrder(String orderId) {
        logger.info("Submitting order: {}", orderId);

        return webClient.patch()
                .uri("/api/partners/{version}/orders/{orderId}/submit", config.getHackload().getApiVersion(), orderId)
                .retrieve()
                .bodyToMono(Void.class)
                .doOnSuccess(response -> logger.info("Order submitted successfully: {}", orderId))
                .doOnError(error -> logger.error("Failed to submit order: {}", error.getMessage()));
    }

    public Mono<Void> confirmOrder(String orderId) {
        logger.info("Confirming order: {}", orderId);

        return webClient.patch()
                .uri("/api/partners/{version}/orders/{orderId}/confirm", config.getHackload().getApiVersion(), orderId)
                .retrieve()
                .bodyToMono(Void.class)
                .doOnSuccess(response -> logger.info("Order confirmed successfully: {}", orderId))
                .doOnError(error -> logger.error("Failed to confirm order: {}", error.getMessage()));
    }

    public Mono<Void> cancelOrder(String orderId) {
        logger.info("Cancelling order: {}", orderId);

        return webClient.patch()
                .uri("/api/partners/{version}/orders/{orderId}/cancel", config.getHackload().getApiVersion(), orderId)
                .retrieve()
                .bodyToMono(Void.class)
                .doOnSuccess(response -> logger.info("Order cancelled successfully: {}", orderId))
                .doOnError(error -> logger.error("Failed to cancel order: {}", error.getMessage()));
    }

    public Mono<Place[]> getPlaces(Integer page, Integer pageSize) {
        logger.info("Getting places, page: {}, size: {}", page, pageSize);

        return webClient.get()
                .uri(uriBuilder -> uriBuilder
                        .path("/api/partners/{version}/places")
                        .queryParam("page", page != null ? page : 1)
                        .queryParam("pageSize", pageSize != null ? pageSize : 20)
                        .build(config.getHackload().getApiVersion()))
                .retrieve()
                .bodyToMono(Place[].class)
                .doOnSuccess(response -> logger.info("Places retrieved: {}", response.length))
                .doOnError(error -> logger.error("Failed to get places: {}", error.getMessage()));
    }

    public Mono<Place> getPlace(String placeId) {
        logger.info("Getting place details for: {}", placeId);

        return webClient.get()
                .uri("/api/partners/{version}/places/{placeId}", config.getHackload().getApiVersion(), placeId)
                .retrieve()
                .bodyToMono(Place.class)
                .doOnSuccess(response -> logger.info("Place details retrieved: {}", response.getId()))
                .doOnError(error -> logger.error("Failed to get place details: {}", error.getMessage()));
    }

    public Mono<Void> selectPlace(String placeId, String orderId) {
        logger.info("Selecting place: {} for order: {}", placeId, orderId);

        SelectPlaceRequest request = new SelectPlaceRequest();
        request.setOrderId(orderId);

        return webClient.patch()
                .uri("/api/partners/{version}/places/{placeId}/select", config.getHackload().getApiVersion(), placeId)
                .bodyValue(request)
                .retrieve()
                .bodyToMono(Void.class)
                .doOnSuccess(response -> logger.info("Place selected successfully: {} for order: {}", placeId, orderId))
                .doOnError(error -> logger.error("Failed to select place: {}", error.getMessage()));
    }

    public Mono<Void> releasePlace(String placeId) {
        logger.info("Releasing place: {}", placeId);

        return webClient.patch()
                .uri("/api/partners/{version}/places/{placeId}/release", config.getHackload().getApiVersion(), placeId)
                .retrieve()
                .bodyToMono(Void.class)
                .doOnSuccess(response -> logger.info("Place released successfully: {}", placeId))
                .doOnError(error -> logger.error("Failed to release place: {}", error.getMessage()));
    }
}
