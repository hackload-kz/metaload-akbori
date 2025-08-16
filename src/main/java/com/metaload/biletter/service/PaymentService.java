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

    /**
     * Обработка redirect после успешной оплаты
     * Вызывается платежным шлюзом при перенаправлении пользователя
     */
    public void notifyPaymentSuccess(Long orderId) {
        logger.info("Payment success notification received for order: {}", orderId);

        Optional<Booking> bookingOpt = findBookingByOrderId(orderId);
        if (bookingOpt.isPresent()) {
            Booking booking = bookingOpt.get();
            booking.setStatus(BookingStatus.CONFIRMED);
            bookingRepository.save(booking);

            logger.info("Booking {} confirmed after successful payment", orderId);
        } else {
            logger.warn("Booking not found for order ID: {}", orderId);
        }
    }

    /**
     * Обработка redirect после неудачной оплаты
     * Вызывается платежным шлюзом при перенаправлении пользователя
     */
    public void notifyPaymentFailure(Long orderId) {
        logger.info("Payment failure notification received for order: {}", orderId);

        Optional<Booking> bookingOpt = findBookingByOrderId(orderId);
        if (bookingOpt.isPresent()) {
            Booking booking = bookingOpt.get();
            booking.setStatus(BookingStatus.CANCELLED);
            bookingRepository.save(booking);

            logger.info("Booking {} cancelled after payment failure", orderId);
        } else {
            logger.warn("Booking not found for order ID: {}", orderId);
        }
    }

    /**
     * Обработка webhook уведомлений от платежного шлюза
     * Вызывается автоматически при изменении статуса платежа
     */
    public void processPaymentNotification(PaymentNotificationPayload payload) {
        logger.info("Processing payment notification: paymentId={}, status={}",
                payload.getPaymentId(), payload.getStatus());

        // Сначала пытаемся найти бронирование по paymentId
        Optional<Booking> bookingOpt = bookingRepository.findByPaymentId(payload.getPaymentId());

        // Если не найдено по paymentId, пытаемся найти по orderId из data
        if (bookingOpt.isEmpty() && payload.getData() != null && payload.getData().get("orderId") != null) {
            String orderIdStr = payload.getData().get("orderId").toString();

            // Сначала пытаемся найти по orderId как строке
            Optional<Booking> byOrderId = findBookingByOrderId(orderIdStr);
            if (byOrderId.isPresent()) {
                bookingOpt = byOrderId;
            } else {
                // Если не найдено, пытаемся как число
                try {
                    Long orderId = Long.parseLong(orderIdStr);
                    bookingOpt = findBookingByOrderId(orderId);
                } catch (NumberFormatException e) {
                    logger.warn("Invalid orderId format in webhook data: {}", orderIdStr);
                }
            }
        }

        if (bookingOpt.isPresent()) {
            Booking booking = bookingOpt.get();

            // Обновляем paymentId если его еще нет
            if (booking.getPaymentId() == null) {
                booking.setPaymentId(payload.getPaymentId());
            }

            switch (payload.getStatus().toUpperCase()) {
                case "CONFIRMED":
                case "COMPLETED":
                    booking.setStatus(BookingStatus.CONFIRMED);
                    logger.info("Payment confirmed for booking: {}", booking.getId());
                    break;

                case "FAILED":
                case "CANCELLED":
                case "REJECTED":
                case "EXPIRED":
                    booking.setStatus(BookingStatus.CANCELLED);
                    logger.info("Payment failed for booking: {}", booking.getId());
                    break;

                case "AUTHORIZED":
                    // Платеж авторизован, но еще не подтвержден
                    logger.info("Payment authorized for booking: {}", booking.getId());
                    break;

                default:
                    logger.info("Unknown payment status: {} for booking: {}",
                            payload.getStatus(), booking.getId());
                    break;
            }

            bookingRepository.save(booking);
        } else {
            logger.warn("No booking found for payment ID: {} or orderId from data", payload.getPaymentId());
        }
    }

    /**
     * Поиск бронирования по orderId
     * Сначала пытаемся найти по ID, затем по orderId
     */
    private Optional<Booking> findBookingByOrderId(Long orderId) {
        // Сначала пытаемся найти по ID (если orderId = booking.id)
        Optional<Booking> byId = bookingRepository.findById(orderId);
        if (byId.isPresent()) {
            return byId;
        }

        // Затем пытаемся найти по orderId
        Optional<Booking> byOrderId = bookingRepository.findByOrderId(orderId.toString());
        if (byOrderId.isPresent()) {
            return byOrderId;
        }

        return Optional.empty();
    }

    /**
     * Поиск бронирования по orderId (строковый)
     */
    private Optional<Booking> findBookingByOrderId(String orderId) {
        return bookingRepository.findByOrderId(orderId);
    }
}
