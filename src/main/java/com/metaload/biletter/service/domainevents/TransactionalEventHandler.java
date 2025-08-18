package com.metaload.biletter.service.domainevents;

import com.metaload.biletter.model.domainevents.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;
import org.springframework.transaction.event.TransactionPhase;
import org.springframework.transaction.event.TransactionalEventListener;

@Component
public class TransactionalEventHandler {

    private static final Logger logger = LoggerFactory.getLogger(TransactionalEventHandler.class);

    private final DomainEventPublisherService domainEventPublisherService;

    public TransactionalEventHandler(DomainEventPublisherService domainEventPublisherService) {
        this.domainEventPublisherService = domainEventPublisherService;
    }

    @TransactionalEventListener(phase = TransactionPhase.AFTER_COMMIT)
    public void handleBookingCreatedEvent(BookingCreatedEvent event) {
        logger.debug("Publishing {} event to Kafka after transaction commit for booking {} event {} ",
                BookingEvents.BOOKING_CREATED, event.getBookingId(), event.getEventId());

        domainEventPublisherService.publishBookingCreated(event);
    }

    @TransactionalEventListener(phase = TransactionPhase.AFTER_COMMIT)
    public void handleBookingCancelledEvent(BookingCancelledEvent event) {
        logger.debug("Publishing {} event to Kafka after transaction commit for booking {} event {}",
                BookingEvents.BOOKING_CANCELLED, event.getBookingId(), event.getEventId());

        domainEventPublisherService.publishBookingCancelled(event);
    }

    @TransactionalEventListener(phase = TransactionPhase.AFTER_COMMIT)
    public void handleSeatAddedToBookingEvent(SeatAddedToBookingEvent event) {
        logger.info("Publishing {} event to Kafka after transaction commit for booking {} seat {}",
                BookingEvents.SEAT_ADDED_TO_BOOKING, event.getBookingId(), event.getSeatId());

        domainEventPublisherService.publishSeatAddedToBooking(event);
    }

    @TransactionalEventListener(phase = TransactionPhase.AFTER_COMMIT)
    public void handleSeatRemovedFromBookingEvent(SeatRemovedFromBookingEvent event) {
        logger.info("Publishing {} event to Kafka after transaction commit for booking {} seat {}",
                BookingEvents.SEAT_REMOVED_FROM_BOOKING, event.getBookingId(), event.getSeatId());

        domainEventPublisherService.publishSeatRemovedFromBooking(event);
    }
}
