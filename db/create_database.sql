-- each measurement represents a state.
CREATE TABLE State (
    state_id INTEGER PRIMARY KEY,
    statetime timestamp NOT NULL
);

-- a feelevel has a key and a value
-- the key is the fee/size-ratio and the value is the amount (tally) of tx in the mempool, seen at a given state
CREATE TABLE Feelevel (
    spb int NOT NULL,
    state_id int NOT NULL,
    tally int, -- amount of tx
    value int, -- fees per sat p byte
    size int,  -- size per sat p byte
    PRIMARY KEY (spb, state_id),
    FOREIGN KEY (state_id) REFERENCES State(state_id) ON DELETE CASCADE
);


-- bitcoin core groups transactions buckets to avoid tracking each transaction feerate independently
-- the smallest bucket has a feerate of 1 sat/byte
-- the bucketsize is spaced exponentially and increases 5% every bucket
-- a bucket cointains a amount (tally) of transaction is it
CREATE TABLE Bucketlevel (
    bucket int NOT NULL,
    state_id int not NULL,
    tally int,
    PRIMARY KEY (bucket, state_id),
    FOREIGN KEY (state_id) REFERENCES STATE(state_id) ON DELETE CASCADE
);


CREATE VIEW v_4hData AS
    SELECT * FROM Feelevel INNER JOIN State ON Feelevel.state_id = State.state_id
    WHERE
        strftime('%s','now') - State.statetime <= 4*60*60; -- 4 * 60 * 60 seconds = 4h

CREATE VIEW v_24hData AS
    SELECT * FROM Feelevel INNER JOIN State ON Feelevel.state_id = State.state_id
    WHERE
        State.state_id % 6 = 0 -- every 6th to have one value every 12 min and 120 total in 24h
        AND
        strftime('%s','now') - State.statetime <= 24*60*60; -- 24 * 60 * 60 seconds = 24h

CREATE VIEW v_7dData AS
    SELECT * FROM Feelevel INNER JOIN State ON Feelevel.state_id = State.state_id
    WHERE
        State.state_id % (6*7) = 0 -- every 42th to have one value every 84 min and 120 total in 7d / 168h
        AND
        strftime('%s','now') - State.statetime <= 7*24*60*60; --  7 * 24 * 60 * 60 seconds = 7d

CREATE VIEW v_usedStateIDs AS
    SELECT state_id FROM v_4hData
    UNION
    SELECT state_id FROM v_24hData
    UNION
    SELECT state_id FROM v_7dData;

CREATE VIEW v_deleteFromBucketlevel AS
    SELECT * FROM Bucketlevel INNER JOIN State ON Bucketlevel.state_id = State.state_id WHERE strftime('%s','now') - State.statetime > 4*60*60; --  4h * 60 min * 60 seconds -> 14400 seconds

CREATE TRIGGER delete_old_data
AFTER INSERT ON State
BEGIN
    DELETE FROM State WHERE state_id NOT IN v_usedStateIDs AND strftime('%s','now') - State.statetime > 4*60*60; -- 4 * 60 * 60 seconds = 4h
END;
