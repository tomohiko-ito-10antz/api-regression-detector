CREATE TABLE example_table (
    id int64 not null,
    c0 string(MAX),
    c1 int64,
    c2 bool,
    c3 timestamp
) PRIMARY KEY (id);
CREATE TABLE child_example_table_1 (
    id int64 not null,
    example_table_id int64,
    FOREIGN KEY (example_table_id) REFERENCES example_table (id)
) PRIMARY KEY (id);
CREATE TABLE child_example_table_2 (
    id int64 not null,
    example_table_id int64
) PRIMARY KEY (id),
    INTERLEAVE IN PARENT example_table;