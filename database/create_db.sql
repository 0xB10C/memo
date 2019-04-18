-- MySQL dump 10.15  Distrib 10.0.38-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: XX.XX.XX.XX    Database: XXXXXXXXXX
-- ------------------------------------------------------
-- Server version	5.7.25

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `current_mempool`
--

DROP TABLE IF EXISTS `current_mempool`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `current_mempool` (
  `id` tinyint(4) NOT NULL DEFAULT '1',
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `mempoolSize` int(11) NOT NULL,
  `byCount` json NOT NULL,
  `positionsInGreedyBlocks` json NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `current_mempool`
--

LOCK TABLES `current_mempool` WRITE;
/*!40000 ALTER TABLE `current_mempool` DISABLE KEYS */;
INSERT INTO `current_mempool` VALUES (1,'2019-04-18 09:30:06',14076161,'{\"1\": 2286, \"2\": 1524, \"3\": 2569, \"4\": 715, \"5\": 574, \"6\": 1309, \"7\": 1185, \"8\": 544, \"9\": 176, \"10\": 450, \"11\": 137, \"12\": 73, \"13\": 89, \"14\": 40, \"15\": 31, \"16\": 21, \"17\": 22, \"18\": 36, \"19\": 60, \"20\": 75, \"21\": 80, \"22\": 61, \"23\": 35, \"24\": 58, \"25\": 48, \"26\": 79, \"27\": 23, \"28\": 14, \"29\": 24, \"30\": 117, \"31\": 43, \"32\": 18, \"33\": 93, \"34\": 22, \"35\": 21, \"36\": 27, \"37\": 7, \"38\": 12, \"39\": 34, \"40\": 60, \"41\": 11, \"42\": 34, \"43\": 19, \"44\": 25, \"45\": 11, \"46\": 66, \"47\": 18, \"48\": 41, \"49\": 110, \"50\": 324, \"51\": 819, \"52\": 177, \"53\": 90, \"54\": 127, \"55\": 115, \"56\": 1194, \"57\": 123, \"58\": 98, \"59\": 52, \"60\": 82, \"61\": 23, \"62\": 24, \"63\": 24, \"64\": 6, \"65\": 4, \"66\": 5, \"68\": 2, \"69\": 3, \"70\": 2, \"71\": 5, \"72\": 3, \"73\": 1, \"74\": 3, \"75\": 1, \"76\": 5, \"77\": 10, \"80\": 1, \"81\": 1, \"82\": 1, \"83\": 1, \"84\": 5, \"86\": 1, \"88\": 1, \"90\": 2, \"99\": 4, \"100\": 5, \"101\": 2, \"102\": 2, \"103\": 4, \"108\": 1, \"110\": 2, \"112\": 3, \"116\": 1, \"120\": 1, \"123\": 2, \"124\": 1, \"128\": 1, \"129\": 2, \"130\": 1, \"143\": 1, \"147\": 1, \"148\": 2, \"150\": 1, \"163\": 1, \"189\": 1, \"200\": 1, \"214\": 1, \"224\": 1, \"248\": 1, \"277\": 1, \"303\": 1}','[13821, 11917, 11417, 9021, 6620, 4775, 3126, 2273, 1662, 1062, 1028, 991, 687]');
/*!40000 ALTER TABLE `current_mempool` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2019-04-18 13:30:44
