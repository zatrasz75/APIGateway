DROP TABLE IF EXISTS stop;

CREATE TABLE IF NOT EXISTS stop (
     id SERIAL PRIMARY KEY,
     stop_list TEXT
);


INSERT INTO stop (stop_list) VALUES ('qwerty');

INSERT INTO stop (stop_list) VALUES ('йцукен');

INSERT INTO stop (stop_list) VALUES ('zxvbnm');
