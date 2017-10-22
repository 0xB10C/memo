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

-- due to performance reasons, the old feelevel data is discarded after 3h for now
-- the main reason beeing the charting libary and the detailed repensentation
-- 3h is a 'just-to-be-safe' value, likely to increase in the future
CREATE TRIGGER delete_old_fee_data
AFTER INSERT ON State
BEGIN
-- 24h * 60 min * 60 seconds -> 86400 seconds
--  4h * 60 min * 60 seconds -> 14400 seconds
--  3h * 60 min * 60 seconds -> 10800 seconds
    DELETE FROM State WHERE strftime('%s','now') - statetime > 10800;
END;

-- DB Version
CREATE TABLE MemoDBVersion (
    version varchar not NULL,
    PRIMARY KEY (version)
)
-- set DB Version to v1.1.0
-- INSERT INTO MEMODBVERSION VALUES("v1.1.0");
