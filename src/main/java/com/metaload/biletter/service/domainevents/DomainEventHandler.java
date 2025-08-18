package com.metaload.biletter.service.domainevents;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.metaload.biletter.model.Seat;
import com.metaload.biletter.model.domainevents.*;
import com.metaload.biletter.service.BookingService;
import com.metaload.biletter.service.EventProviderService;
import com.metaload.biletter.service.EventService;
import com.metaload.biletter.service.SeatService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.kafka.support.KafkaHeaders;
import org.springframework.messaging.handler.annotation.Header;
import org.springframework.messaging.handler.annotation.Payload;
import org.springframework.stereotype.Service;

@Service
public class DomainEventHandler {

    private static final Logger logger = LoggerFactory.getLogger(DomainEventHandler.class);

    private final ObjectMapper objectMapper;
    private final BookingService bookingService;
    private final EventProviderService eventProviderService;
    private final SeatService seatService;

    public DomainEventHandler(ObjectMapper objectMapper,
                              BookingService bookingService,
                              EventProviderService eventProviderService,
                              SeatService seatService) {
        this.objectMapper = objectMapper;
        this.bookingService = bookingService;
        this.eventProviderService = eventProviderService;
        this.seatService = seatService;
    }

    @KafkaListener(topics = "${kafka.topics.booking-events}", groupId = "domain-event-handlers")
    public void handleDomainEvent(
            @Payload String eventPayload,
            @Header(KafkaHeaders.RECEIVED_KEY) String eventKey,
            @Header(KafkaHeaders.RECEIVED_TOPIC) String topic,
            @Header(KafkaHeaders.RECEIVED_PARTITION) int partition) {

        try {
            logger.info("Received domain event from topic: {}, partition: {}, key: {}",
                    topic, partition, eventKey);

            DomainEvent domainEvent = objectMapper.readValue(eventPayload, DomainEvent.class);

            // Обрабатываем событие в зависимости от типа
            switch (domainEvent.getEventType()) {
                case BookingEvents.BOOKING_CREATED:
                    handleBookingCreatedEvent(domainEvent);
                    break;
                case BookingEvents.BOOKING_CANCELLED:
                    handleBookingCancelledEvent(domainEvent);
                    break;
                case BookingEvents.SEAT_ADDED_TO_BOOKING:
                    handleSeatAddedToBookingEvent(domainEvent);
                    break;
                case BookingEvents.SEAT_REMOVED_FROM_BOOKING:
                    handleSeatRemovedFromBookingEvent(domainEvent);
                    break;
                default:
                    logger.warn("Unknown event type: {}", domainEvent.getEventType());
            }

            logger.info("Successfully processed domain event: {} with id: {}",
                    domainEvent.getEventType(), domainEvent.getEventId());

        } catch (Exception e) {
            logger.error("Failed to process domain event with key: {}", eventKey, e);
            // В продакшене здесь можно добавить retry механизм или отправку в DLQ
        }
    }

    private void handleBookingCreatedEvent(DomainEvent domainEvent) {
        try {
            // Десериализуем данные события
            BookingCreatedEvent bookingEvent = objectMapper.readValue(
                    domainEvent.getEventData(), BookingCreatedEvent.class);

            logger.info("Processing {} event for booking: {}",
                    BookingEvents.BOOKING_CREATED, bookingEvent.getBookingId());

            // Выполняем бизнес-логику
            if (EventService.MAIN_EVENT.equals(bookingEvent.getEventId())) {
                bookingService.createOrderForBooking(bookingEvent.getBookingId());
            }

        } catch (Exception e) {
            logger.error("Failed to handle {} event: {}",
                    BookingEvents.BOOKING_CREATED, domainEvent.getEventId(), e);
            //throw e;
        }
    }

    private void handleBookingCancelledEvent(DomainEvent domainEvent) {
        try {
            // Десериализуем данные события
            BookingCancelledEvent cancelEvent = objectMapper.readValue(
                    domainEvent.getEventData(), BookingCancelledEvent.class);

            logger.info("Processing {} event for booking: {} with previous status: {}",
                    BookingEvents.BOOKING_CANCELLED, cancelEvent.getBookingId(), cancelEvent.getPreviousStatus());

            // Выполняем дополнительную бизнес-логику при отмене брони
            // Например, отмена платежа, уведомление пользователя, отмена во внешних системах
            if (EventService.MAIN_EVENT.equals(cancelEvent.getEventId())) {
                if (cancelEvent.getOrderId() != null) {
                    // Отмена заказа в ивент провайдере
                    //eventProviderService.cancelOrder(cancelEvent.getOrderId());
                    logger.info("Booking {} cancelled for main event, order: {}",
                            cancelEvent.getBookingId(), cancelEvent.getOrderId());
                } else {
                    logger.info("Booking {} cancelled for main event, the event has no order",
                            cancelEvent.getBookingId());
                }
            }

            // Логирование для аудита
            logger.info("Booking {} cancelled by user: {} at: {}",
                    cancelEvent.getBookingId(), cancelEvent.getUserId(), cancelEvent.getCancelledAt());

        } catch (Exception e) {
            logger.error("Failed to handle {}: {}",
                    BookingEvents.BOOKING_CANCELLED, domainEvent.getEventId(), e);
            //throw e;
        }
    }

    private void handleSeatAddedToBookingEvent(DomainEvent domainEvent) {
        try {
            // Десериализуем данные события
            SeatAddedToBookingEvent seatEvent = objectMapper.readValue(
                    domainEvent.getEventData(), SeatAddedToBookingEvent.class);

            logger.info("Processing {} event for booking: {} seat: {}",
                    BookingEvents.SEAT_ADDED_TO_BOOKING, seatEvent.getBookingId(), seatEvent.getSeatId());

            // Выбор места в ивент провайдере
            Seat seat = seatService.findById(seatEvent.getSeatId());
            //eventProviderService.selectPlace(seat.getProviderSeatId(), seatEvent.getOrderId());

        } catch (Exception e) {
            logger.error("Failed to handle {} event: {}",
                    BookingEvents.SEAT_ADDED_TO_BOOKING, domainEvent.getEventId(), e);
            //throw e;
        }
    }

    private void handleSeatRemovedFromBookingEvent(DomainEvent domainEvent) {
        try {
            // Десериализуем данные события
            SeatRemovedFromBookingEvent seatEvent = objectMapper.readValue(
                    domainEvent.getEventData(), SeatRemovedFromBookingEvent.class);

            logger.info("Processing {} event for booking: {} seat: {}",
                    BookingEvents.SEAT_REMOVED_FROM_BOOKING, seatEvent.getBookingId(), seatEvent.getSeatId());

            // Освобождение места в ивент провайдере
            Seat seat = seatService.findById(seatEvent.getSeatId());
            //eventProviderService.releasePlace(seat.getProviderSeatId());

        } catch (Exception e) {
            logger.error("Failed to handle {} event: {}",
                    BookingEvents.SEAT_REMOVED_FROM_BOOKING, domainEvent.getEventId(), e);
            //throw e;
        }
    }

}
