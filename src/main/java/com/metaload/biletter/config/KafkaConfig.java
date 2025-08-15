package com.metaload.biletter.config;

import org.apache.kafka.clients.admin.NewTopic;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.kafka.config.TopicBuilder;

@Configuration
public class KafkaConfig {

    @Value("${kafka.topics.booking-events}")
    private String bookingEventsTopic;

    @Value("${kafka.topics.payment-events}")
    private String paymentEventsTopic;

    @Value("${kafka.topics.seat-election-events}")
    private String seatSelectionEventsTopic;

    @Bean
    public NewTopic bookingEventsTopic() {
        return TopicBuilder.name(bookingEventsTopic)
                .partitions(3)
                .replicas(1)
                .build();
    }

    @Bean
    public NewTopic paymentEventsTopic() {
        return TopicBuilder.name(paymentEventsTopic)
                .partitions(3)
                .replicas(1)
                .build();
    }

    @Bean
    public NewTopic seatSelectionEventsTopic() {
        return TopicBuilder.name(seatSelectionEventsTopic)
                .partitions(3)
                .replicas(1)
                .build();
    }
}
