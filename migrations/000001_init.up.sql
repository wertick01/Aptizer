
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
    participants INTEGER,
    PRIMARY KEY(id),
    FOREIGN KEY(author_id) REFERENCES users(userid)
);

CREATE TABLE participants (
    userid INTEGER NOT NULL,
    headline_id INTEGER NOT NULL,
    FOREIGN KEY(userid) REFERENCES users(userid),
    FOREIGN KEY(headline_id) REFERENCES news(id)
);

CREATE TABLE tags (
    tag_id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
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
    ('89512590156', 'Pavel', 'Desyukevich', 'Yurievich', 'cool.pumba01@yandex.ru', 'qwerty', 'Very stupid moran', 'morans/moran.jpg', 1),
    ('89614780777', 'Vyacheslav', 'Bogatov', 'Alexandrovich', 'adolph@yandex.ru', 'qwerty', 'Very stupid nazi', 'nazis/nazi.jpg', 2);

INSERT INTO news (headline_date, title, headline_text, photo, author_id, participants) VALUES 
    (NOW(), 'First title', 'Author has fucked your mummy', 'their/photo.jpg', 1, 1), 
    (NOW(), 'Second title', 'Author has fucked your mummy again', 'their/new/photo.jpg', 1, 2),
    (NOW(), 'Third title', 'Author has fucked your mummy thirdly', 'their/third/photo.jpg', 2, 2),
    (NOW(), 'Fourth title', 'Author has fucked your mummy fourthly', 'their/fourth/photo.jpg', 1, 3);

INSERT INTO tags (tag) VALUES
    ('fuck you'),
    ('fuck you again'),
    ('third fuck');

INSERT INTO tags_news (headline_id, tag_id) VALUES
    (1, 1),
    (1, 2),
    (2, 1),
    (2, 2),
    (3, 1),
    (3, 2),
    (4, 3);


INSERT INTO participants (userid, headline_id) VALUES
    (1, 1),
    (1, 2),
    (1, 3),
    (2, 1),
    (2, 2),
    (2, 4);