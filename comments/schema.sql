DROP TABLE IF EXISTS comments;

CREATE TABLE comments (
                          id SERIAL PRIMARY KEY,
                          news_id INT,
                          content TEXT NOT NULL DEFAULT 'empty',
                          pubtime BIGINT NOT NULL DEFAULT extract (epoch from now())


);

INSERT INTO comments(news_id,content)  VALUES (1,'комментарий');

INSERT INTO comments(news_id,content)  VALUES (2,'ups  проверка');
