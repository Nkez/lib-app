CREATE TABLE books
(
    id                SERIAL        NOT NULL UNIQUE PRIMARY KEY,
    book_title        VARCHAR(255)  NOT NULL,
    book_title_native VARCHAR(255),
    book_price        DECIMAL(6, 2) NOT NULL,
    inventory_count   INT           NOT NULL,
    books_in_lib      INT           NOT NULL,
    one_day_price     DECIMAL(6, 2) NOT NULL,
    year_of_published DATE          not null,
    registration_date DATE          not null default CURRENT_DATE,
    number_of_pages   INT           NOT NULL,
    book_state        BOOLEAN       NOT NULL,
    hide_book         BOOLEAN,
    rating            DECIMAL(6, 2)
);

CREATE TABLE book_copies
(
    id           SERIAL NOT NULL UNIQUE PRIMARY KEY,
    defect       VARCHAR(255),
    defect_photo VARCHAR(255)
);

CREATE TABLE books_book_copies
(
    books_id       int references books (id) on delete cascade,
    book_copies_id int references book_copies (id) on delete cascade,
    PRIMARY KEY (books_id, book_copies_id)
);



CREATE TABLE books_photo
(
    id         SERIAL       NOT NULL UNIQUE PRIMARY KEY,
    foto_name varchar(255) ,
    book_photo VARCHAR(255) NOT NULL
);

CREATE TABLE books_books_photo
(
    books_id       int references books (id) on delete cascade,
    books_photo_id int references books_photo (id) on delete cascade,
    PRIMARY KEY (books_id, books_photo_id)
);

CREATE TABLE authors
(
    id               SERIAL       NOT NULL UNIQUE PRIMARY KEY,
    author_firstname varchar(50)  NOT NULL,
    author_lastname  varchar(50)  NOT NULL,
    foto_name varchar(255) ,
    author_photo     VARCHAR(255) NOT NULL
);

CREATE TABLE books_authors
(
    books_id   int references books (id) on delete cascade,
    authors_id int references authors (id) on delete cascade,
    PRIMARY KEY (books_id, authors_id)
);


CREATE TABLE genres
(
    id    SERIAL      NOT NULL UNIQUE PRIMARY KEY,
    genre VARCHAR(40) NOT NULL
);


CREATE TABLE book_genre
(
    book_id  int references books (id) on delete cascade,
    genre_id int references genres (id) on delete cascade,
    PRIMARY KEY (book_id, genre_id)
);


CREATE TABLE users
(
    id              SERIAL PRIMARY KEY NOT NULL,
    last_name       VARCHAR(50)        NOT NULL,
    first_name      VARCHAR(50)        NOT NULL,
    middle_name     VARCHAR(50),
    passport_number VARCHAR(30) UNIQUE,
    birthday        DATE               NOT NULL,
    email_address   VARCHAR(40) UNIQUE NOT NULL,
    address         VARCHAR(150)
);


CREATE TABLE books_users
(
    order_date     DATE          NOT NULL,
    date_to_return DATE          NOT NULL,
    price          DECIMAL(6, 2) NOT NULL,
    rating         DECIMAL(6, 2),
    is_return      BOOLEAN       NOT NULL,
    books_id       int ,
    users_id       int ,
    FOREIGN KEY (books_id) references books (id) on delete cascade,
    FOREIGN KEY (users_id) references users (id) on delete cascade
);

-- order_date     DATE          NOT NULL,
--     date_to_return DATE          NOT NULL,
--     price          DECIMAL(6, 2) NOT NULL,
--     rating         DECIMAL(6, 2),
--     is_return      BOOLEAN       NOT NULL,
--     books_id       int references books (id) on delete cascade,
--      users_id       int references users (id) on delete cascade,
--  PRIMARY KEY (books_id, users_id)
