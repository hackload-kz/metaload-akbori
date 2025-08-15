package com.metaload.biletter;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.scheduling.annotation.EnableAsync;

@SpringBootApplication
@EnableAsync
public class BiletterApplication {

    public static void main(String[] args) {
        SpringApplication.run(BiletterApplication.class, args);
    }
}
