package com.metaload.biletter.repository;

import com.metaload.biletter.model.Event;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.time.LocalDate;
import java.util.List;

@Repository
public interface EventRepository extends JpaRepository<Event, Long> {

    @Query("SELECT e FROM Event e WHERE " +
            "(:query IS NULL OR LOWER(e.title) LIKE LOWER(CONCAT('%', :query, '%'))) AND " +
            "(:date IS NULL OR DATE(e.createdAt) = :date)")
    List<Event> findEventsByQueryAndDate(@Param("query") String query, @Param("date") LocalDate date);

    List<Event> findByExternal(Boolean external);
}
