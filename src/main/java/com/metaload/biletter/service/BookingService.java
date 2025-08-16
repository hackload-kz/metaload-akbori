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
    private final SeatRepository seatRepository;
    private final BookingSeatRepository bookingSeatRepository;
    private final UserService userService;

    public BookingService(BookingRepository bookingRepository,
            SeatRepository seatRepository, BookingSeatRepository bookingSeatRepository,
            UserService userService) {
        this.bookingRepository = bookingRepository;
        this.seatRepository = seatRepository;
        this.bookingSeatRepository = bookingSeatRepository;
        this.userService = userService;
    }

    public Booking createBooking(CreateBookingRequest request) {
        Booking booking = new Booking();
        booking.setEventId(request.getEventId());
        booking.setUserId(userService.getCurrentUser().getUserId());
        booking.setStatus(Booking.BookingStatus.PENDING);

        return bookingRepository.save(booking);
    }

    public List<ListBookingsResponseItem> getAllBookings() {
        Integer currentUserId = userService.getCurrentUser().getUserId();
        List<Booking> bookings = bookingRepository.findByUserId(currentUserId);

        return bookings.stream()
                .map(this::mapToResponseItem)
                .collect(Collectors.toList());
    }

    public String initiatePayment(Long bookingId) {
        Booking booking = findById(bookingId);
        if (booking.getStatus() != Booking.BookingStatus.PENDING) {
            throw new RuntimeException("Cannot initiate payment for booking with status: " + booking.getStatus());
        }
        booking.setStatus(Booking.BookingStatus.PAYMENT_PENDING);
        bookingRepository.save(booking);

        // Возвращаем URL для оплаты (в реальном приложении это будет URL платежного
        // шлюза)
        return "https://payment-gateway.example.com/pay/" + bookingId;
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
        item.setEventId(booking.getEventId());

        // Получаем места для этого бронирования
        List<BookingSeat> bookingSeats = bookingSeatRepository.findByBookingId(booking.getId());
        List<ListBookingsResponseItemSeat> seats = bookingSeats.stream()
                .map(bs -> new ListBookingsResponseItemSeat(bs.getSeat().getId()))
                .collect(Collectors.toList());

        item.setSeats(seats);
        return item;
    }
}
