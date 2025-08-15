package com.metaload.biletter.repository;

import com.metaload.biletter.model.Booking;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;

@Repository
public interface BookingRepository extends JpaRepository<Booking, Long> {

    List<Booking> findByEventId(Long eventId);

    @Query("SELECT b FROM Booking b WHERE b.status = 'PENDING' OR b.status = 'PAYMENT_PENDING'")
    List<Booking> findActiveBookings();

    @Query("SELECT b FROM Booking b WHERE b.paymentId = :paymentId")
    Optional<Booking> findByPaymentId(@Param("paymentId") String paymentId);

    @Query("SELECT b FROM Booking b WHERE b.event.id = :eventId AND b.status IN ('PENDING', 'PAYMENT_PENDING')")
    List<Booking> findActiveBookingsByEventId(@Param("eventId") Long eventId);
}
