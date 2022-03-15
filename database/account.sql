CREATE TABLE account
(
    id bigint auto_increment not null ,
    password varchar(50) not null,
    type tinyint,
    primary key(id)
);

INSERT INTO account (password, type) VALUES ('demo', 1);