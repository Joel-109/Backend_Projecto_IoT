CREATE TABLE Category(
    id_category VARCHAR(50) PRIMARY KEY 
);

CREATE TABLE Status(
    id_status VARCHAR(50) PRIMARY KEY 
);

CREATE TABLE Dish(
    id_dish INTEGER PRIMARY KEY,
    description VARCHAR(100) NOT NULL,
    name VARCHAR(60) NOT NULL,
    image TEXT NOT NULL,
    price FLOAT NOT NULL,
    id_category VARCHAR(50) NOT NULL,
    FOREIGN KEY (id_category) REFERENCES Category(id_category)
);

CREATE TABLE Orders(
    id_order INTEGER PRIMARY KEY ,
    id_desk INTEGER,
    id_status VARCHAR(50),
    FOREIGN KEY (id_status) REFERENCES Status(id_status)
);

CREATE TABLE Invoices(
    id_invoice INTEGER PRIMARY KEY,
    id_order INTEGER,
    id_dish INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (id_order) REFERENCES Orders(id_order),
    FOREIGN KEY (id_dish) REFERENCES Dish(id_dish)
); 
