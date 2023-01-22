DROP TABLE IF EXISTS example_table;
CREATE TABLE example_table (
    id int64 not null,
    c0 string(MAX),
    c1 int64,
    c2 bool,
    c3 timestamp
) primary key(id);