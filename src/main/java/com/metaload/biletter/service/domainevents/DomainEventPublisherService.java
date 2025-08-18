package com.metaload.biletter.service.domainevents;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.github.f4b6a3.uuid.UuidCreator;
import com.metaload.biletter.model.domainevents.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;

@Service
public class DomainEventPublisherService {

    private static final Logger logger = LoggerFactory.getLogger(DomainEventPublisherService.class);

    private final KafkaTemplate<String, String> kafkaTemplate;
    private final ObjectMapper objectMapper;

    @Value("${kafka.topics.booking-events}")
    private String bookingEventsTopic;

    public DomainEventPublisherService(KafkaTemplate<String, String> kafkaTemplate, ObjectMapper objectMapper) {
        this.kafkaTemplate = kafkaTemplate;
        this.objectMapper = objectMapper;
    }

    public void publishBookingCreated(BookingCreatedEvent data) {
        try {
            String eventData = objectMapper.writeValueAsString(data);

            DomainEvent event = new DomainEvent();
            event.setEventId(UuidCreator.getTimeOrderedEpoch().toString());
            event.setEventData(eventData);
            event.setEventType(BookingEvents.BOOKING_CREATED);
            event.setTimestamp(LocalDateTime.now());

            String eventStr = objectMapper.writeValueAsString(event);

            kafkaTemplate.send(bookingEventsTopic, data.getBookingId().toString(), eventStr)
                    .whenComplete((result, ex) -> {
                        if (ex == null) {
                            logger.info("Successfully published {} event for booking {}",
                                    BookingEvents.BOOKING_CREATED, data.getBookingId());
                        } else {
                            logger.error("Failed to publish {} event for booking {}",
                                    BookingEvents.BOOKING_CREATED, data.getBookingId(), ex);
                        }
                    });
        } catch (JsonProcessingException e) {
            logger.error("Failed to serialize {} event for booking {}",
                    BookingEvents.BOOKING_CREATED, data.getBookingId(), e);
        }
    }

    public void publishBookingCancelled(BookingCancelledEvent data) {
        try {
            String eventData = objectMapper.writeValueAsString(data);

            DomainEvent event = new DomainEvent();
            event.setEventId(UuidCreator.getTimeOrderedEpoch().toString());
            event.setEventData(eventData);
            event.setEventType(BookingEvents.BOOKING_CANCELLED);
            event.setTimestamp(LocalDateTime.now());
            String eventStr = objectMapper.writeValueAsString(event);
            kafkaTemplate.send(bookingEventsTopic, data.getBookingId().toString(), eventStr)
                    .whenComplete((result, ex) -> {
                        if (ex == null) {
                            logger.info("Successfully published {} event for booking {}",
                                    BookingEvents.BOOKING_CANCELLED, data.getBookingId());
                        } else {
                            logger.error("Failed to publish {} event for booking {}",
                                    BookingEvents.BOOKING_CANCELLED, data.getBookingId(), ex);
                        }
                    });
        } catch (JsonProcessingException e) {
            logger.error("Failed to serialize {} event for booking {}",
                    BookingEvents.BOOKING_CANCELLED, data.getBookingId(), e);
        }
    }

    public void publishSeatAddedToBooking(SeatAddedToBookingEvent data) {
        try {
            String eventData = objectMapper.writeValueAsString(data);

            DomainEvent event = new DomainEvent();
            event.setEventId(UuidCreator.getTimeOrderedEpoch().toString());
            event.setEventData(eventData);
            event.setEventType(BookingEvents.SEAT_ADDED_TO_BOOKING);
            event.setTimestamp(LocalDateTime.now());

            String eventStr = objectMapper.writeValueAsString(event);

            kafkaTemplate.send(bookingEventsTopic, data.getBookingId().toString(), eventStr)
                    .whenComplete((result, ex) -> {
                        if (ex == null) {
                            logger.info("Successfully published {} event for booking {} seat {}",
                                    BookingEvents.SEAT_ADDED_TO_BOOKING, data.getBookingId(), data.getSeatId());
                        } else {
                            logger.error("Failed to publish {} event for booking {} seat {}",
                                    BookingEvents.SEAT_ADDED_TO_BOOKING, data.getBookingId(), data.getSeatId(), ex);
                        }
                    });
        } catch (JsonProcessingException e) {
            logger.error("Failed to serialize {} event for booking {} seat {}",
                    BookingEvents.SEAT_ADDED_TO_BOOKING, data.getBookingId(), data.getSeatId(), e);
        }
    }

    public void publishSeatRemovedFromBooking(SeatRemovedFromBookingEvent data) {
        try {
            String eventData = objectMapper.writeValueAsString(data);

            DomainEvent event = new DomainEvent();
            event.setEventId(UuidCreator.getTimeOrderedEpoch().toString());
            event.setEventData(eventData);
            event.setEventType(BookingEvents.SEAT_REMOVED_FROM_BOOKING);
            event.setTimestamp(LocalDateTime.now());

            String eventStr = objectMapper.writeValueAsString(event);

            kafkaTemplate.send(bookingEventsTopic, data.getBookingId().toString(), eventStr)
                    .whenComplete((result, ex) -> {
                        if (ex == null) {
                            logger.info("Successfully published {} event for booking {} seat {}",
                                    BookingEvents.SEAT_REMOVED_FROM_BOOKING, data.getBookingId(), data.getSeatId());
                        } else {
                            logger.error("Failed to publish {} event for booking {} seat {}",
                                    BookingEvents.SEAT_REMOVED_FROM_BOOKING, data.getBookingId(), data.getSeatId(), ex);
                        }
                    });
        } catch (JsonProcessingException e) {
            logger.error("Failed to serialize {} event for booking {} seat {}",
                    BookingEvents.SEAT_REMOVED_FROM_BOOKING, data.getBookingId(), data.getSeatId(), e);
        }
    }
}
