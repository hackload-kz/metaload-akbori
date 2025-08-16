-- Create events table
CREATE TABLE events (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    external BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create seats table
CREATE TABLE seats (
    id BIGSERIAL PRIMARY KEY,
    event_id BIGINT NOT NULL REFERENCES events(id),
    row_number INTEGER NOT NULL,
    seat_number INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'FREE' CHECK (status IN ('FREE', 'RESERVED', 'SOLD')),
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(event_id, row_number, seat_number)
);

-- Create bookings table
CREATE TABLE bookings (
    id BIGSERIAL PRIMARY KEY,
    event_id BIGINT NOT NULL REFERENCES events(id),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'PAYMENT_PENDING', 'CONFIRMED', 'CANCELLED')),
    total_amount DECIMAL(10,2),
    currency VARCHAR(3) DEFAULT 'RUB',
    payment_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create booking_seats table (many-to-many relationship)
CREATE TABLE booking_seats (
    id BIGSERIAL PRIMARY KEY,
    booking_id BIGINT NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    seat_id BIGINT NOT NULL REFERENCES seats(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(booking_id, seat_id)
);

-- Create payment_notifications table
CREATE TABLE payment_notifications (
    id BIGSERIAL PRIMARY KEY,
    payment_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    team_slug VARCHAR(255),
    timestamp TIMESTAMP NOT NULL,
    data JSONB,
    processed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX idx_seats_event_id ON seats(event_id);
CREATE INDEX idx_seats_status ON seats(status);
CREATE INDEX idx_bookings_event_id ON bookings(event_id);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_booking_seats_booking_id ON booking_seats(booking_id);
CREATE INDEX idx_booking_seats_seat_id ON booking_seats(seat_id);
CREATE INDEX idx_payment_notifications_payment_id ON payment_notifications(payment_id);
CREATE INDEX idx_payment_notifications_processed ON payment_notifications(processed);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_events_updated_at BEFORE UPDATE ON events
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_seats_updated_at BEFORE UPDATE ON seats
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_bookings_updated_at BEFORE UPDATE ON bookings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
