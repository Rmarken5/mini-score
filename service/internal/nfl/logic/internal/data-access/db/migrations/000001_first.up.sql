CREATE EXTENSION "uuid-ossp";

CREATE TABLE TEAM
(
    ID           UUID                              DEFAULT uuid_generate_v4() PRIMARY KEY,
    NAME         TEXT                     NOT NULL,
    ABBREVIATION TEXT NOT NULL,
    CREATED_AT   timestamp with time zone NOT NULL DEFAULT NOW(),
    UPDATED_AT   timestamp with time zone NOT NULL DEFAULT NOW(),
    DELETED_AT   timestamp with time zone
);

CREATE TABLE GAME
(
    ID           TEXT PRIMARY KEY          NOT NULL,
    GAME_TIME    timestamp with time zone,
    QUARTER      TEXT,
    GAME_CLOCK   TEXT,
    AWAY_TEAM    UUID                     NOT NULL,
    HOME_TEAM    UUID                     NOT NULL,
    CREATED_AT   timestamp with time zone NOT NULL DEFAULT NOW(),
    UPDATED_AT   timestamp with time zone NOT NULL DEFAULT NOW(),
    DELETED_AT   timestamp with time zone,

    FOREIGN KEY (AWAY_TEAM) REFERENCES TEAM (ID),
    FOREIGN KEY (HOME_TEAM) REFERENCES TEAM (ID)
);

CREATE TABLE GAME_QUARTER_SCORE
(
    ID         UUID                              DEFAULT uuid_generate_v4() PRIMARY KEY,
    GAME_ID    text                     NOT NULL,
    TEAM_ID    UUID                     NOT NULL,
    QUARTER    TEXT                     NOT NULL,
    SCORE      INT                      NOT NULL DEFAULT 0,
    CREATED_AT timestamp with time zone NOT NULL DEFAULT NOW(),
    UPDATED_AT timestamp with time zone NOT NULL DEFAULT NOW(),
    DELETED_AT timestamp with time zone,

    FOREIGN KEY (GAME_ID) REFERENCES GAME (ID),
    FOREIGN KEY (TEAM_ID) REFERENCES TEAM (ID),

    CONSTRAINT GAME_TEAM_QUARTER UNIQUE (GAME_ID, TEAM_ID, QUARTER)
);

