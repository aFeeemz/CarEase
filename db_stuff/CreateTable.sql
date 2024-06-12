-- Create the Users table with isadmin column
CREATE TABLE Users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    deposit_amount DECIMAL(10, 2) DEFAULT 0,
    isadmin BOOLEAN DEFAULT FALSE
);

-- Create the Category table
CREATE TABLE Category (
    id SERIAL PRIMARY KEY,
    category_name VARCHAR(255) NOT NULL
);

-- Create the Cars table
CREATE TABLE Cars (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    availability BOOLEAN NOT NULL DEFAULT TRUE,
    rental_costs DECIMAL(10, 2) NOT NULL,
    vin_number VARCHAR(255) NOT NULL,
    color VARCHAR(255) NOT NULL,
    transmission VARCHAR(255) NOT NULL,
    category_id INTEGER NOT NULL,
    FOREIGN KEY (category_id) REFERENCES Category(id)
);

-- Create the RentalHistory table
CREATE TABLE RentalHistory (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    car_id INTEGER NOT NULL,
    rental_date DATE NOT NULL,
    return_date DATE NOT NULL,
    actual_return_date DATE,
    total_cost DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES Users(id),
    FOREIGN KEY (car_id) REFERENCES Cars(id)
);