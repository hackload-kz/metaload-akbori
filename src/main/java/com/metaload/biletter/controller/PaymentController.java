package com.metaload.biletter.controller;

import com.metaload.biletter.dto.PaymentNotificationPayload;
import com.metaload.biletter.service.PaymentService;
import jakarta.validation.Valid;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/payments")
public class PaymentController {

    private final PaymentService paymentService;

    public PaymentController(PaymentService paymentService) {
        this.paymentService = paymentService;
    }

    @GetMapping("/success")
    public ResponseEntity<String> notifyPaymentSuccess(@RequestParam Long orderId) {
        paymentService.notifyPaymentSuccess(orderId);
        return ResponseEntity.ok("OK");
    }

    @GetMapping("/fail")
    public ResponseEntity<String> notifyPaymentFailed(@RequestParam Long orderId) {
        paymentService.notifyPaymentFailure(orderId);
        return ResponseEntity.ok("OK");
    }

    @PostMapping("/notifications")
    public ResponseEntity<String> onPaymentUpdates(@Valid @RequestBody PaymentNotificationPayload payload) {
        paymentService.processPaymentNotification(payload);
        return ResponseEntity.ok("OK");
    }
}
