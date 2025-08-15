package com.metaload.biletter.service;

import com.metaload.biletter.dto.PaymentNotificationPayload;
import com.metaload.biletter.model.Booking;
import com.metaload.biletter.model.Booking.BookingStatus;
import com.metaload.biletter.repository.BookingRepository;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.Optional;

@Service
@Transactional
public class PaymentService {

    private static final Logger logger = LoggerFactory.getLogger(PaymentService.class);

    private final BookingRepository bookingRepository;

    public PaymentService(BookingRepository bookingRepository) {
        this.bookingRepository = bookingRepository;
    }

    public void notifyPaymentSuccess(Long orderId) {
        logger.info("Payment success notification received for order: {}", orderId);

        Optional<Booking> bookingOpt = bookingRepository.findById(orderId);
        if (bookingOpt.isPresent()) {
            Booking booking = bookingOpt.get();
            booking.setStatus(BookingStatus.CONFIRMED);
            bookingRepository.save(booking);

            logger.info("Booking {} confirmed after successful payment", orderId);
        } else {
            logger.warn("Booking not found for order ID: {}", orderId);
        }
    }

    public void notifyPaymentFailure(Long orderId) {
        logger.info("Payment failure notification received for order: {}", orderId);

        Optional<Booking> bookingOpt = bookingRepository.findById(orderId);
        if (bookingOpt.isPresent()) {
            Booking booking = bookingOpt.get();
            booking.setStatus(BookingStatus.CANCELLED);
            bookingRepository.save(booking);

            logger.info("Booking {} cancelled after payment failure", orderId);
        } else {
            logger.warn("Booking not found for order ID: {}", orderId);
        }
    }

    public void processPaymentNotification(PaymentNotificationPayload payload) {
        logger.info("Processing payment notification: paymentId={}, status={}",
                payload.getPaymentId(), payload.getStatus());

        // Находим бронирование по paymentId
        Optional<Booking> bookingOpt = bookingRepository.findByPaymentId(payload.getPaymentId());

        if (bookingOpt.isPresent()) {
            Booking booking = bookingOpt.get();

            switch (payload.getStatus().toUpperCase()) {
                case "CONFIRMED":
                case "COMPLETED":
                    booking.setStatus(BookingStatus.CONFIRMED);
                    logger.info("Payment confirmed for booking: {}", booking.getId());
                    break;

                case "FAILED":
                case "CANCELLED":
                case "REJECTED":
                    booking.setStatus(BookingStatus.CANCELLED);
                    logger.info("Payment failed for booking: {}", booking.getId());
                    break;

                default:
                    logger.info("Unknown payment status: {} for booking: {}",
                            payload.getStatus(), booking.getId());
                    break;
            }

            bookingRepository.save(booking);
        } else {
            logger.warn("No booking found for payment ID: {}", payload.getPaymentId());
        }
    }
}
