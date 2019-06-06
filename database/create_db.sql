SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;


CREATE TABLE `current_mempool` (
  `id` tinyint(4) NOT NULL DEFAULT '1',
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `mempoolSize` int(11) NOT NULL,
  `byCount` json NOT NULL,
  `positionsInGreedyBlocks` json NOT NULL
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

CREATE TABLE `historicalMempool` (
  `timeframe` int(11) NOT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `countInBuckets` json NOT NULL,
  `feeInBuckets` json NOT NULL,
  `sizeInBuckets` json NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `pastBlocks` (
  `height` int(11) NOT NULL,
  `txCount` int(11) NOT NULL,
  `size` int(11) NOT NULL,
  `weight` int(11) NOT NULL,
  `timestamp` timestamp NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

CREATE TABLE `timeInMempool` (
  `id` int(11) NOT NULL,
  `timeAxis` json NOT NULL,
  `feerateAxis` json NOT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `transactionsStats` (
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `segwitCount` int(11) NOT NULL,
  `rbfCount` int(11) NOT NULL,
  `txCount` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


ALTER TABLE `current_mempool`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `id` (`id`);

ALTER TABLE `historicalMempool`
  ADD PRIMARY KEY (`timeframe`,`timestamp`),
  ADD UNIQUE KEY `timeframe` (`timeframe`,`timestamp`);

ALTER TABLE `pastBlocks`
  ADD PRIMARY KEY (`height`),
  ADD UNIQUE KEY `height` (`height`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;