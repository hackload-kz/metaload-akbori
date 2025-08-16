package com.metaload.biletter.repository;

import com.metaload.biletter.model.Event;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;

import java.time.LocalDateTime;

@Repository
public interface EventRepository extends JpaRepository<Event, Long> {
    @Query("select e from Event e where (e.title ilike concat('%',:query,'%') or e.description ilike concat('%',:query,'%')) and e.datetimeStart >= :date order by e.datetimeStart")
    Page<Event> find(String query, LocalDateTime date, Pageable pageable);

    @Query("select e from Event e where (:date is null or e.datetimeStart >= :date) order by e.datetimeStart")
    Page<Event> find(LocalDateTime date, Pageable pageable);

    @Query("select e from Event e where (e.title ilike concat('%',:query,'%') or e.description ilike concat('%',:query,'%')) order by e.datetimeStart")
    Page<Event> find(String query, Pageable pageable);

}
