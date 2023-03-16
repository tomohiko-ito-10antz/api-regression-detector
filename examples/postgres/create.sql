DROP TABLE IF EXISTS example_table CASCADE;
CREATE TABLE example_table (
    id serial,
    c0 text,
    c1 integer,
    c2 boolean,
    c3 timestamptz,
    PRIMARY KEY (id)
);
DROP TABLE IF EXISTS child_example_table_1;
CREATE TABLE child_example_table_1 (
    id serial,
    example_table_id integer,
    PRIMARY KEY (id),
    FOREIGN KEY (example_table_id) REFERENCES example_table (id)
);
DROP TABLE IF EXISTS child_example_table_2;
CREATE TABLE child_example_table_2 (
    id serial,
    example_table_id integer,
    PRIMARY KEY (id),
    FOREIGN KEY (example_table_id) REFERENCES example_table (id)
);