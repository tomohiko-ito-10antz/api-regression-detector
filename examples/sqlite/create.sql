DROP TABLE IF EXISTS example_table;
CREATE TABLE example_table (
    id integer PRIMARY KEY AUTOINCREMENT,
    c0 text,
    c1 integer,
    c2 boolean,
    c3 datetime
);
DROP TABLE IF EXISTS child_example_table_1;
CREATE TABLE child_example_table_1 (
    id integer PRIMARY KEY AUTOINCREMENT,
    example_table_id integer,
    FOREIGN KEY (example_table_id) REFERENCES example_table (id)
);
DROP TABLE IF EXISTS child_example_table_2;
CREATE TABLE child_example_table_2 (
    id integer,
    example_table_id integer,
    PRIMARY KEY(id),
    FOREIGN KEY (example_table_id) REFERENCES example_table (id)
);
PRAGMA foreign_keys = 1;
