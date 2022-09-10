
CREATE TABLE roles (
    role_id INTEGER NOT NULL PRIMARY KEY,
    userrole VARCHAR(100) NOT NULL
);

CREATE TABLE users (
    userid INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT, 
    userphone VARCHAR(100) NOT NULL,
    username VARCHAR(100) NOT NULL,
    usersurname VARCHAR(100) NOT NULL,
    userpatrynomic VARCHAR(100),
    useremail VARCHAR(100) NOT NULL,
    userhash VARCHAR(100) NOT NULL,
    userdescription VARCHAR(1000) NOT NULL,
    userphoto VARCHAR(150) NOT NULL,
    userrole INTEGER NOT NULL,
    FOREIGN KEY(userrole) REFERENCES roles(role_id)
    );

CREATE TABLE news (
    id INTEGER NOT NULL AUTO_INCREMENT, 
    headline_date TIMESTAMP NOT NULL,
    title VARCHAR(200) NOT NULL, 
    headline_text VARCHAR(100), 
    photo VARCHAR(150) NOT NULL, 
    author_id INTEGER NOT NULL,
    PRIMARY KEY(id),
    FOREIGN KEY(author_id) REFERENCES users(userid)
);

CREATE TABLE tags (
    tag_id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT DEFAULT 0,
    tag VARCHAR(100) NOT NULL
 );

CREATE TABLE tags_news (
    headline_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    FOREIGN KEY(headline_id) REFERENCES news(id),
    FOREIGN KEY(tag_id) REFERENCES tags(tag_id)
);

INSERT INTO roles (role_id, userrole) VALUES 
    (1, 'admin'), 
    (2, 'user');

INSERT INTO users (userphone, username, usersurname, userpatrynomic, useremail, userhash, userdescription, userphoto, userrole) VALUES 
    ('89512590156', 'Pavel', 'Desyukevich', 'Yurievich', 'cool.pumba01@yandex.ru', 'qwerty', 'Very stupid moran', 'morans/moran.jpg', 1);

INSERT INTO news (headline_date, title, headline_text, photo, author_id) VALUES 
    (DATE(NOW()), 'First title', 'Author has fucked your mummy', 'their/photo.jpg', 1), 
    (DATE(NOW()), 'Second title', 'Author has fucked your mummy again', 'their/new/photo.jpg', 1);

INSERT INTO tags (tag) VALUES
    ('fuck you'),
    ('fuck you again');

INSERT INTO tags_news (headline_id, tag_id) VALUES
    (3, 1),
    (3, 2),
    (4, 1),
    (4, 2);

