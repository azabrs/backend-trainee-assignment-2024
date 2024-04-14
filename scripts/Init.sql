CREATE TABLE banners_data (
    id INT PRIMARY KEY,
    content VARCHAR(100) NOT NULL

);

CREATE TABLE banners (
    id SERIAL PRIMARY KEY,
    feature_id INT NOT NULL,
    data_id INT NOT NULL UNIQUE,
    is_active BOOLEAN NOT NULL,
    FOREIGN KEY (data_id) REFERENCES banners_data(id)
);



CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE banner_tags (
    banner_id INT NOT NULL,
    tag_id INT NOT NULL,
    PRIMARY KEY (banner_id, tag_id),
    FOREIGN KEY (banner_id) REFERENCES banners(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON UPDATE CASCADE
);

CREATE TABLE users(
    login TEXT PRIMARY KEY,
    password_hash TEXT NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE
);