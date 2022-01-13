CREATE TABLE books
(
    id                SERIAL        NOT NULL UNIQUE PRIMARY KEY,
    book_title        VARCHAR(150)  NOT NULL,
    book_title_native VARCHAR(150),
    book_price        DECIMAL(6, 2) NOT NULL,
    inventory_count   INT           NOT NULL,
    books_in_lib      INT           NOT NULL,
    one_day_price     DECIMAL(6, 2) NOT NULL,
    year_of_published DATE  not null,
    registration_date DATE     not null default CURRENT_DATE,
    number_of_pages   INT           NOT NULL,
    book_state        BOOLEAN       NOT NULL,
    hide_book         BOOLEAN,
    rating            INT,
    book_foto         VARCHAR(255)  NOT NULL
);

CREATE TABLE authors
(
    id                SERIAL       NOT NULL UNIQUE PRIMARY KEY,
    authors_firstname varchar(150) NOT NULL,
    authors_lastname  varchar(150) NOT NULL,
    authors_foto      varchar(255) NOT NULL
);

CREATE TABLE genres
(
    id    SERIAL       NOT NULL UNIQUE PRIMARY KEY,
    genre VARCHAR(150) NOT NULL

);

CREATE TABLE book_author
(
    book_id   int references books (id) on delete cascade,
    author_id int references authors (id) on delete cascade,
    PRIMARY KEY (book_id, author_id)
);0

CREATE TABLE book_genre
(
    book_id  int references books (id) on delete cascade,
    genre_id int references genres (id) on delete cascade,
    PRIMARY KEY (book_id, genre_id)
);


CREATE TABLE users
(
    id              SERIAL PRIMARY KEY NOT NULL,
    last_name       VARCHAR(40)        NOT NULL,
    first_name      VARCHAR(40)        NOT NULL,
    middle_name     VARCHAR(40),
    passport_number VARCHAR(30) UNIQUE,
    birthday        DATE               NOT NULL,
    email_address   VARCHAR(40) UNIQUE NOT NULL,
    address         VARCHAR(50)
);


CREATE TABLE order_cart
(
    last_name      VARCHAR(40)        NOT NULL,
    first_name     VARCHAR(40)        NOT NULL,
    email_address  VARCHAR(40) UNIQUE NOT NULL,
    book1          VARCHAR(40),
    book2          VARCHAR(40),
    book3          VARCHAR(40),
    book4          VARCHAR(40),
    book5          VARCHAR(40),
    price          DECIMAL(6, 2),
    date_to_return DATE          not null
);

CREATE TABLE return_cart
(
    date_to_return DATE  not null default CURRENT_DATE,
    price DECIMAL(6,2) NOT NULL,
    rating INT,
    defect_foto VARCHAR(255),
    defect VARCHAR(255),
    is_book_defect BOOLEAN
);