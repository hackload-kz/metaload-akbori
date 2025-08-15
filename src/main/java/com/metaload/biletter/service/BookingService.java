package com.metaload.biletter.service;

import com.metaload.biletter.dto.CreateBookingRequest;
import com.metaload.biletter.dto.ListBookingsResponseItem;
import com.metaload.biletter.dto.ListBookingsResponseItemSeat;
import com.metaload.biletter.model.Booking;
import com.metaload.biletter.model.BookingSeat;
import com.metaload.biletter.model.Event;
import com.metaload.biletter.model.Seat;
import com.metaload.biletter.repository.BookingRepository;
import com.metaload.biletter.repository.BookingSeatRepository;
import com.metaload.biletter.repository.SeatRepository;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.stream.Collectors;

@Service
@Transactional
public class BookingService {

    private final BookingRepository bookingRepository;
    private final EventService eventService;
    private final SeatRepository seatRepository;
    private final BookingSeatRepository bookingSeatRepository;

    public BookingService(BookingRepository bookingRepository, EventService eventService,
            SeatRepository seatRepository, BookingSeatRepository bookingSeatRepository) {
        this.bookingRepository = bookingRepository;
        this.eventService = eventService;
        this.seatRepository = seatRepository;
        this.bookingSeatRepository = bookingSeatRepository;
    }

    public Booking createBooking(CreateBookingRequest request) {
        Event event = eventService.findById(request.getEventId());

        Booking booking = new Booking();
        booking.setEvent(event);
        booking.setStatus(Booking.BookingStatus.PENDING);

        return bookingRepository.save(booking);
    }

    public List<ListBookingsResponseItem> getAllBookings() {
        List<Booking> bookings = bookingRepository.findAll();

        return bookings.stream()
                .map(this::mapToResponseItem)
                .collect(Collectors.toList());
    }

    public void initiatePayment(Long bookingId) {
        Booking booking = findById(bookingId);
        booking.setStatus(Booking.BookingStatus.PAYMENT_PENDING);
        bookingRepository.save(booking);
    }

    public void cancelBooking(Long bookingId) {
        Booking booking = findById(bookingId);

        // Освобождаем все места
        List<BookingSeat> bookingSeats = bookingSeatRepository.findByBookingId(bookingId);
        for (BookingSeat bookingSeat : bookingSeats) {
            Seat seat = bookingSeat.getSeat();
            seat.setStatus(Seat.SeatStatus.FREE);
            seatRepository.save(seat);
        }

        // Удаляем связи с местами
        bookingSeatRepository.deleteAll(bookingSeats);

        // Отменяем бронирование
        booking.setStatus(Booking.BookingStatus.CANCELLED);
        bookingRepository.save(booking);
    }

    public void selectSeat(Long bookingId, Long seatId) {
        Booking booking = findById(bookingId);
        Seat seat = seatRepository.findById(seatId)
                .orElseThrow(() -> new RuntimeException("Seat not found"));

        // Проверяем, что место свободно
        if (seat.getStatus() != Seat.SeatStatus.FREE) {
            throw new RuntimeException("Seat is not available");
        }

        // Резервируем место
        seat.setStatus(Seat.SeatStatus.RESERVED);
        seatRepository.save(seat);

        // Создаем связь
        BookingSeat bookingSeat = new BookingSeat();
        bookingSeat.setBooking(booking);
        bookingSeat.setSeat(seat);
        bookingSeatRepository.save(bookingSeat);
    }

    public void releaseSeat(Long seatId) {
        Seat seat = seatRepository.findById(seatId)
                .orElseThrow(() -> new RuntimeException("Seat not found"));

        // Находим связь с бронированием
        List<BookingSeat> bookingSeats = bookingSeatRepository.findBySeatId(seatId);

        if (!bookingSeats.isEmpty()) {
            // Удаляем связь
            bookingSeatRepository.deleteAll(bookingSeats);

            // Освобождаем место
            seat.setStatus(Seat.SeatStatus.FREE);
            seatRepository.save(seat);
        }
    }

    public Booking findById(Long id) {
        return bookingRepository.findById(id)
                .orElseThrow(() -> new RuntimeException("Booking not found with id: " + id));
    }

    private ListBookingsResponseItem mapToResponseItem(Booking booking) {
        ListBookingsResponseItem item = new ListBookingsResponseItem();
        item.setId(booking.getId());
        item.setEventId(booking.getEvent().getId());

        // Получаем места для этого бронирования
        List<BookingSeat> bookingSeats = bookingSeatRepository.findByBookingId(booking.getId());
        List<ListBookingsResponseItemSeat> seats = bookingSeats.stream()
                .map(bs -> new ListBookingsResponseItemSeat(bs.getSeat().getId()))
                .collect(Collectors.toList());

        item.setSeats(seats);
        return item;
    }
}
