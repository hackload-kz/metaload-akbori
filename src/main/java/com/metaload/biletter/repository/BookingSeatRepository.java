package com.metaload.biletter.repository;

import com.metaload.biletter.model.BookingSeat;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;

@Repository
public interface BookingSeatRepository extends JpaRepository<BookingSeat, Long> {

    List<BookingSeat> findByBookingId(Long bookingId);

    List<BookingSeat> findByBookingIdIn(List<Long> bookingIds);

    @Query("SELECT bs FROM BookingSeat bs WHERE bs.booking.id = :bookingId AND bs.seat.id = :seatId")
    Optional<BookingSeat> findByBookingIdAndSeatId(@Param("bookingId") Long bookingId, @Param("seatId") Long seatId);

    @Query("SELECT bs FROM BookingSeat bs WHERE bs.seat.id = :seatId")
    List<BookingSeat> findBySeatId(@Param("seatId") Long seatId);

    @Query("SELECT COUNT(bs) FROM BookingSeat bs WHERE bs.booking.id = :bookingId")
    Long countByBookingId(@Param("bookingId") Long bookingId);
}
