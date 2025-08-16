package com.metaload.biletter.service;

import com.metaload.biletter.dto.CreateBookingRequest;
import com.metaload.biletter.dto.ListBookingsResponseItem;
import com.metaload.biletter.dto.ListBookingsResponseItemSeat;
import com.metaload.biletter.dto.payment.PaymentInitRequest;
import com.metaload.biletter.model.Booking;
import com.metaload.biletter.model.BookingSeat;
import com.metaload.biletter.model.Event;
import com.metaload.biletter.model.Seat;
import com.metaload.biletter.repository.BookingRepository;
import com.metaload.biletter.repository.BookingSeatRepository;
import com.metaload.biletter.repository.SeatRepository;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.util.List;
import java.util.stream.Collectors;

@Service
@Transactional
public class BookingService {

    private final BookingRepository bookingRepository;
    private final SeatRepository seatRepository;
    private final BookingSeatRepository bookingSeatRepository;
    private final PaymentGatewayService paymentGatewayService;
    private final EventService eventService;

    public BookingService(BookingRepository bookingRepository,
                          SeatRepository seatRepository, BookingSeatRepository bookingSeatRepository,
                          PaymentGatewayService paymentGatewayService, EventService eventService) {
        this.bookingRepository = bookingRepository;
        this.seatRepository = seatRepository;
        this.bookingSeatRepository = bookingSeatRepository;
        this.paymentGatewayService = paymentGatewayService;
        this.eventService = eventService;
    }

    public Booking createBooking(CreateBookingRequest request) {
        Booking booking = new Booking();
        Event event = eventService.findById(request.getEventId());
        booking.setEvent(event);
        booking.setStatus(Booking.BookingStatus.PENDING);

        // Генерируем уникальный orderId
        String orderId = generateOrderId();
        booking.setOrderId(orderId);

        // Устанавливаем userId (в реальном приложении берем из контекста безопасности)
        booking.setUserId(1); // TODO: Заменить на получение из контекста пользователя

        return bookingRepository.save(booking);
    }

    /**
     * Генерирует уникальный orderId для бронирования
     */
    private String generateOrderId() {
        String timestamp = LocalDateTime.now().format(DateTimeFormatter.ofPattern("yyyyMMddHHmmss"));
        String random = String.valueOf((int) (Math.random() * 1000));
        return "BK" + timestamp + random;
    }

    public List<ListBookingsResponseItem> getAllBookings() {
        // В реальном приложении фильтруем по текущему пользователю
        Integer currentUserId = 1; // TODO: Заменить на получение из контекста безопасности
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

        // Обновляем статус на PAYMENT_PENDING
        booking.setStatus(Booking.BookingStatus.PAYMENT_PENDING);
        bookingRepository.save(booking);

        try {
            // Создаем запрос на создание платежа
            PaymentInitRequest paymentRequest = paymentGatewayService.createPaymentRequest(
                    booking.getOrderId(),
                    booking.getTotalAmount() != null ? booking.getTotalAmount().longValue() * 100 : 0L, // Конвертируем
                                                                                                        // в копейки
                    booking.getCurrency() != null ? booking.getCurrency() : "RUB",
                    "Оплата бронирования #" + booking.getOrderId(),
                    "user@example.com" // В реальном приложении берем из контекста пользователя
            );

            // Создаем платеж в платежном шлюзе
            return paymentGatewayService.createPayment(paymentRequest)
                    .map(response -> response.getPaymentURL())
                    .block(); // В реальном приложении лучше использовать async подход

        } catch (Exception e) {
            // В случае ошибки возвращаем бронирование в исходное состояние
            booking.setStatus(Booking.BookingStatus.PENDING);
            bookingRepository.save(booking);
            throw new RuntimeException("Failed to initiate payment: " + e.getMessage(), e);
        }
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

    public Booking findByOrderId(String orderId) {
        return bookingRepository.findByOrderId(orderId)
                .orElseThrow(() -> new RuntimeException("Booking not found with orderId: " + orderId));
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
