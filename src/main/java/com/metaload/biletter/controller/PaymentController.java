package com.metaload.biletter.controller;

import com.metaload.biletter.dto.PaymentNotificationPayload;
import com.metaload.biletter.service.PaymentService;
import jakarta.validation.Valid;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/payments")
public class PaymentController {

    private static final Logger logger = LoggerFactory.getLogger(PaymentController.class);

    private final PaymentService paymentService;

    public PaymentController(PaymentService paymentService) {
        this.paymentService = paymentService;
    }

    /**
     * Эндпоинт для redirect после успешной оплаты
     * Вызывается платежным шлюзом для перенаправления пользователя
     */
    @GetMapping("/success")
    public ResponseEntity<String> notifyPaymentSuccess(@RequestParam String orderId) {
        logger.info("Payment success redirect received for order: {}", orderId);

        try {
            // Парсим orderId (может быть строкой или числом)
            Long orderIdLong = Long.parseLong(orderId);
            paymentService.notifyPaymentSuccess(orderIdLong);
            logger.info("Successfully processed payment success for order: {}", orderId);
            return ResponseEntity.ok("Payment successful! Your booking has been confirmed.");
        } catch (NumberFormatException e) {
            logger.error("Invalid orderId format: {}", orderId);
            return ResponseEntity.badRequest().body("Invalid order ID format");
        } catch (Exception e) {
            logger.error("Error processing payment success for order: {}", orderId, e);
            return ResponseEntity.internalServerError().body("Error processing payment success");
        }
    }

    /**
     * Эндпоинт для redirect после неудачной оплаты
     * Вызывается платежным шлюзом для перенаправления пользователя
     */
    @GetMapping("/fail")
    public ResponseEntity<String> notifyPaymentFailed(@RequestParam String orderId) {
        logger.info("Payment failure redirect received for order: {}", orderId);

        try {
            Long orderIdLong = Long.parseLong(orderId);
            paymentService.notifyPaymentFailure(orderIdLong);
            logger.info("Successfully processed payment failure for order: {}", orderId);
            return ResponseEntity.ok("Payment failed. Your booking has been cancelled.");
        } catch (NumberFormatException e) {
            logger.error("Invalid orderId format: {}", orderId);
            return ResponseEntity.badRequest().body("Invalid order ID format");
        } catch (Exception e) {
            logger.error("Error processing payment failure for order: {}", orderId, e);
            return ResponseEntity.internalServerError().body("Error processing payment failure");
        }
    }

    /**
     * Webhook эндпоинт для получения уведомлений от платежного шлюза
     * Вызывается автоматически платежным шлюзом при изменении статуса платежа
     */
    @PostMapping("/notifications")
    public ResponseEntity<String> onPaymentUpdates(@Valid @RequestBody PaymentNotificationPayload payload) {
        logger.info("Payment webhook notification received: paymentId={}, status={}, orderId={}",
                payload.getPaymentId(), payload.getStatus(),
                payload.getData() != null ? payload.getData().get("orderId") : "N/A");

        try {
            paymentService.processPaymentNotification(payload);
            logger.info("Successfully processed payment webhook for paymentId: {}", payload.getPaymentId());
            return ResponseEntity.ok("OK");
        } catch (Exception e) {
            logger.error("Error processing payment webhook for paymentId: {}", payload.getPaymentId(), e);
            return ResponseEntity.internalServerError().body("Error processing webhook");
        }
    }
}
