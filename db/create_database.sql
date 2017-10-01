CREATE TABLE State (
    state_id INTEGER PRIMARY KEY,
    statetime timestamp NOT NULL
);

CREATE TABLE Feelevel (
    spb int NOT NULL,
    state_id int NOT NULL,
    tally int,
    PRIMARY KEY (spb,state_id),
    FOREIGN KEY (state_id) REFERENCES State(state_id) ON DELETE CASCADE
);

CREATE TRIGGER delete_old_fee_data
    AFTER INSERT ON State
BEGIN
    -- 24h * 60 min * 60 seconds -> 86400 seconds
    --  3h * 60 min * 60 seconds -> 10800 seconds
    DELETE FROM State WHERE strftime('%s','now') - statetime > 10800;
END;
