package com.metaload.biletter.service;

import com.github.f4b6a3.uuid.UuidCreator;
import com.metaload.biletter.dto.CreateBookingRequest;
import com.metaload.biletter.dto.ListBookingsResponseItem;
import com.metaload.biletter.dto.ListBookingsResponseItemSeat;
import com.metaload.biletter.dto.event.CreateOrderResponse;
import com.metaload.biletter.dto.payment.PaymentInitRequest;
import com.metaload.biletter.dto.payment.PaymentInitResponse;
import com.metaload.biletter.model.*;
import com.metaload.biletter.model.domainevents.BookingCreatedEvent;
import com.metaload.biletter.model.domainevents.SeatAddedToBookingEvent;
import com.metaload.biletter.model.domainevents.SeatRemovedFromBookingEvent;
import com.metaload.biletter.repository.BookingRepository;
import com.metaload.biletter.repository.BookingSeatRepository;
import com.metaload.biletter.repository.SeatRepository;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.ApplicationEventPublisher;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.List;
import java.util.stream.Collectors;

@Service
public class BookingService {
    private static final Logger log = LoggerFactory.getLogger(BookingService.class);

    private final BookingRepository bookingRepository;
    private final SeatRepository seatRepository;
    private final BookingSeatRepository bookingSeatRepository;
    private final PaymentGatewayService paymentGatewayService;
    private final UserService userService;
    private final EventService eventService;
    private final EventProviderService eventProviderService;
    private final ApplicationEventPublisher applicationEventPublisher;

    @Autowired
    private BookingService selfProxy;

    public BookingService(BookingRepository bookingRepository,
                          SeatRepository seatRepository,
                          BookingSeatRepository bookingSeatRepository,
                          PaymentGatewayService paymentGatewayService,
                          UserService userService,
                          EventService eventService,
                          EventProviderService eventProviderService,
                          ApplicationEventPublisher applicationEventPublisher) {
        this.bookingRepository = bookingRepository;
        this.seatRepository = seatRepository;
        this.bookingSeatRepository = bookingSeatRepository;
        this.paymentGatewayService = paymentGatewayService;
        this.userService = userService;
        this.eventService = eventService;
        this.eventProviderService = eventProviderService;
        this.applicationEventPublisher = applicationEventPublisher;
    }

    @Transactional
    public Booking createBooking(CreateBookingRequest request) {
        Booking booking = new Booking();
        Event event = eventService.findById(request.getEventId());
        booking.setEvent(event);
        booking.setStatus(Booking.BookingStatus.PENDING);

        // Устанавливаем userId
        User currentUser = userService.getCurrentUser();
        booking.setUserId(currentUser.getUserId());

        Booking savedBooking = bookingRepository.save(booking);

        // Сейчас Заказ создается асинхронно, можно сделать синхронным если нужно
        // createOrderForBooking(savedBooking.getId());

        // Публикуем доменное событие о создании брони
        BookingCreatedEvent bookingEvent = new BookingCreatedEvent(
                savedBooking.getId(),
                savedBooking.getOrderId(),
                savedBooking.getEvent().getId(),
                savedBooking.getUserId(),
                savedBooking.getStatus().name(),
                savedBooking.getCreatedAt()
        );
        applicationEventPublisher.publishEvent(bookingEvent);

        return savedBooking;
    }

    @Transactional(readOnly = true)
    public List<ListBookingsResponseItem> getAllBookings() {
        // В реальном приложении фильтруем по текущему пользователю
        Integer currentUserId = userService.getCurrentUser().getUserId();
        List<Booking> bookings = bookingRepository.findByUserId(currentUserId);

        return bookings.stream()
                .map(this::mapToResponseItem)
                .collect(Collectors.toList());
    }

    @Transactional
    public String initiatePayment(Long bookingId) {
        // Используем SELECT FOR UPDATE для предотвращения конкурентного доступа
        Booking booking = bookingRepository.findByIdForUpdate(bookingId)
                .orElseThrow(() -> new RuntimeException("Booking not found with id: " + bookingId));

        if (booking.getStatus() != Booking.BookingStatus.PENDING) {
            throw new RuntimeException("Cannot initiate payment for booking with status: " + booking.getStatus());
        }

        // Вычисляем total_amount как сумму цен всех забронированных мест
        List<BookingSeat> bookingSeats = bookingSeatRepository.findByBookingId(bookingId);
        if (bookingSeats.isEmpty()) {
            throw new RuntimeException("Cannot initiate payment for booking without selected seats");
        }

        // Проверяем, что все места действительно зарезервированы этим бронированием
        for (BookingSeat bookingSeat : bookingSeats) {
            if (bookingSeat.getSeat().getStatus() != Seat.SeatStatus.RESERVED) {
                throw new RuntimeException("Seat " + bookingSeat.getSeat().getId() + " is not in RESERVED status");
            }
        }

        BigDecimal calculatedTotalAmount = bookingSeats.stream()
                .map(bookingSeat -> bookingSeat.getSeat().getPrice())
                .reduce(BigDecimal.ZERO, BigDecimal::add);

        // Обновляем статус на PAYMENT_PENDING и устанавливаем рассчитанную сумму
        booking.setStatus(Booking.BookingStatus.PAYMENT_PENDING);
        booking.setTotalAmount(calculatedTotalAmount);
        Booking savedBooking = bookingRepository.saveAndFlush(booking);
        long totalAmountInTiyn = calculatedTotalAmount.multiply(new BigDecimal("100")).longValue();

        String email = userService.getCurrentUser().getEmail();
        try {
            // Создаем запрос на создание платежа
            PaymentInitRequest paymentRequest = paymentGatewayService.createPaymentRequest(
                    savedBooking.getOrderId(),
                    totalAmountInTiyn,
                    "KZT",
                    "Оплата бронирования #" + savedBooking.getOrderId(),
                    email
            );

            // Создаем платеж в платежном шлюзе
            return paymentGatewayService.createPayment(paymentRequest)
                    .map(PaymentInitResponse::getPaymentURL)
                    .block();

        } catch (Exception e) {
            // В случае ошибки возвращаем бронирование в исходное состояние
            savedBooking.setStatus(Booking.BookingStatus.PENDING);
            savedBooking.setTotalAmount(null);
            bookingRepository.saveAndFlush(savedBooking);
            throw new RuntimeException("Failed to initiate payment: " + e.getMessage(), e);
        }
    }

    @Transactional
    public void cancelBooking(Long bookingId) {
        Booking booking = findById(bookingId);

        // Освобождаем все места с пессимистичными блокировками
        List<BookingSeat> bookingSeats = bookingSeatRepository.findByBookingId(bookingId);
        for (BookingSeat bookingSeat : bookingSeats) {
            releaseSeat(bookingSeat.getSeat().getId());
        }

        // Отменяем бронирование
        booking.setStatus(Booking.BookingStatus.CANCELLED);
        bookingRepository.save(booking);
    }

    @Transactional
    public void selectSeat(Long bookingId, Long seatId) {
        Booking booking = findById(bookingId);

        // Используем пессимистичную блокировку для места
        Seat seat = seatRepository.findByIdForUpdate(seatId)
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

        // Публикуем локальное Spring событие, которое будет обработано после коммита транзакции
        SeatAddedToBookingEvent seatEvent = new SeatAddedToBookingEvent(
                booking.getId(),
                seat.getId(),
                booking.getEvent().getId(),
                booking.getOrderId(),
                booking.getUserId(),
                seat.getRowNumber(),
                seat.getSeatNumber(),
                LocalDateTime.now()
        );
        applicationEventPublisher.publishEvent(seatEvent);
    }

    @Transactional
    public void releaseSeat(Long seatId) {
        // Используем пессимистичную блокировку для места
        Seat seat = seatRepository.findByIdForUpdate(seatId)
                .orElseThrow(() -> new RuntimeException("Seat not found"));

        // Находим связь с бронированием
        List<BookingSeat> bookingSeats = bookingSeatRepository.findBySeatId(seatId);

        if (!bookingSeats.isEmpty()) {
            // Сохраняем информацию о брони до удаления связи
            BookingSeat bookingSeat = bookingSeats.get(0);
            Booking booking = bookingSeat.getBooking();

            // Удаляем связь
            bookingSeatRepository.deleteAll(bookingSeats);

            // Освобождаем место
            seat.setStatus(Seat.SeatStatus.FREE);
            seatRepository.save(seat);

            // Публикуем локальное Spring событие о удалении места из брони
            SeatRemovedFromBookingEvent seatEvent = new SeatRemovedFromBookingEvent(
                    booking.getId(),
                    seat.getId(),
                    booking.getEvent().getId(),
                    booking.getOrderId(),
                    booking.getUserId(),
                    seat.getRowNumber(),
                    seat.getSeatNumber(),
                    LocalDateTime.now()
            );
            applicationEventPublisher.publishEvent(seatEvent);
        }
    }

    @Transactional(readOnly = true)
    public Booking findById(Long id) {
        return bookingRepository.findById(id)
                .orElseThrow(() -> new RuntimeException("Booking not found with id: " + id));
    }

    @Transactional(readOnly = true)
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

    @Transactional
    public void createOrderForBooking(Long bookingId) {
        Booking booking = findById(bookingId);

        if (EventService.MAIN_EVENT.equals(booking.getEvent().getId())) {
            CreateOrderResponse response = eventProviderService.createOrder().block();
            booking.setOrderId(response.getOrderId());
        } else {
            String orderId = UuidCreator.getTimeOrderedEpoch().toString();
            booking.setOrderId(orderId);
        }

        bookingRepository.save(booking);
    }
}
