package com.metaload.biletter.repository;

import com.metaload.biletter.model.Seat;
import com.metaload.biletter.model.Seat.SeatStatus;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Lock;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import jakarta.persistence.LockModeType;

import java.util.List;
import java.util.Optional;

@Repository
public interface SeatRepository extends JpaRepository<Seat, Long> {

    @Lock(LockModeType.PESSIMISTIC_WRITE)
    @Query("SELECT s FROM Seat s WHERE s.id = :seatId")
    Optional<Seat> findByIdForUpdate(@Param("seatId") Long seatId);

    Page<Seat> findByEventIdAndStatus(Long eventId, SeatStatus status, Pageable pageable);

    @Query("SELECT s FROM Seat s WHERE s.event.id = :eventId " +
            "AND (:row IS NULL OR s.rowNumber = :row) " +
            "AND (:status IS NULL OR s.status = :status)")
    Page<Seat> findByEventIdAndFilters(
            @Param("eventId") Long eventId,
            @Param("row") Integer row,
            @Param("status") SeatStatus status,
            Pageable pageable);

    List<Seat> findByEventIdAndStatus(Long eventId, SeatStatus status);

    Optional<Seat> findByEventIdAndRowNumberAndSeatNumber(Long eventId, Integer rowNumber, Integer seatNumber);

    @Query("SELECT s FROM Seat s WHERE s.id = :seatId AND s.status = 'FREE'")
    Optional<Seat> findFreeSeatById(@Param("seatId") Long seatId);
}
