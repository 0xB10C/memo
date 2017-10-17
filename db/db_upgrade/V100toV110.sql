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

-- DB Version
CREATE TABLE MemoDBVersion (
    version int not NULL,
    PRIMARY KEY (version)
)

-- set DB Version to v1.1.0
INSERT INTO MemoDBVersion VALUES("v1.1.0");
