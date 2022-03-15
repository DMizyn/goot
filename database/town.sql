CREATE TABLE town
(
    id bigint auto_increment not null,
    name varchar(50) not null unique,
    temple_posx mediumint not null,
    temple_posy mediumint not null,
    temple_posz tinyint not null,
    primary key(id),
    CONSTRAINT UC_TOWN_TEMPLE UNIQUE (temple_posx, temple_posy, temple_posz)
);

INSERT INTO town (name, temple_posx, temple_posy, temple_posz) VALUES ('goots', 10, 10, 10);