# ************************************************************
# Sequel Pro SQL dump
# Version 4541
#
# http://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: localhost (MySQL 5.5.42)
# Database: kudo_admin
# Generation Time: 2017-11-19 09:29:05 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table admin
# ------------------------------------------------------------

DROP TABLE IF EXISTS `admin`;

CREATE TABLE `admin` (
  `id_admin` int(11) NOT NULL AUTO_INCREMENT,
  `nama_admin` varchar(255) NOT NULL,
  `email_admin` varchar(255) NOT NULL,
  `peran` varchar(255) NOT NULL,
  `kata_sandi_admin` varchar(255) NOT NULL,
  `tanggal_diperbaharui` datetime NOT NULL,
  PRIMARY KEY (`id_admin`),
  UNIQUE KEY `email_admin` (`email_admin`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

LOCK TABLES `admin` WRITE;
/*!40000 ALTER TABLE `admin` DISABLE KEYS */;

INSERT INTO `admin` (`id_admin`, `nama_admin`, `email_admin`, `peran`, `kata_sandi_admin`, `tanggal_diperbaharui`)
VALUES
	(1,'Kudo Admin','admin@kudo.co.id','Administrator','$2a$14$HDuB/fqVA0P1rpago3HM4.tueRNXbceOQ6h5pVvRXfBekXtQ/PfM2','2017-11-19 13:50:00'),
	(2,'CS Kudo','cs@kudo.co.id','Administrator','$2a$14$CjBIFNPrvTwDnixBHM8YxOrFST1yxsR0ncRTdZ5dHu1Xm1rbedQZS','2017-11-19 13:34:28');

/*!40000 ALTER TABLE `admin` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table akses
# ------------------------------------------------------------

DROP TABLE IF EXISTS `akses`;

CREATE TABLE `akses` (
  `id_akses` int(11) NOT NULL AUTO_INCREMENT,
  `nama_akses` varchar(255) NOT NULL,
  `id_grup_akses` int(11) NOT NULL,
  PRIMARY KEY (`id_akses`),
  KEY `id_grup_akses` (`id_grup_akses`),
  CONSTRAINT `fk_akses_id_grup_akses` FOREIGN KEY (`id_grup_akses`) REFERENCES `grup_akses` (`id_grup_akses`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

LOCK TABLES `akses` WRITE;
/*!40000 ALTER TABLE `akses` DISABLE KEYS */;

INSERT INTO `akses` (`id_akses`, `nama_akses`, `id_grup_akses`)
VALUES
	(1,'Belanja voucher pulsa',1),
	(2,'Belanja bahan sembako',1),
	(3,'Membeli banyak voucher pulsa untuk tujuan disimpan (di-stok)',2),
	(4,'Membeli banyak bahan sembako untuk tujuan disimpan (di-stok)',2),
	(5,'Menjual kelebihan stok voucher pulsa ke pikak Kudo',3),
	(6,'Menjual kelebihan stok bahan sembako ke pikak Kudo',3);

/*!40000 ALTER TABLE `akses` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table grup
# ------------------------------------------------------------

DROP TABLE IF EXISTS `grup`;

CREATE TABLE `grup` (
  `id_grup` int(11) NOT NULL AUTO_INCREMENT,
  `id_grup_pengguna` int(11) NOT NULL,
  `id_grup_akses` int(11) NOT NULL,
  PRIMARY KEY (`id_grup`),
  KEY `id_grup_pengguna` (`id_grup_pengguna`),
  KEY `id_grup_akses` (`id_grup_akses`),
  CONSTRAINT `fk_grup_id_grup_akses` FOREIGN KEY (`id_grup_akses`) REFERENCES `grup_akses` (`id_grup_akses`),
  CONSTRAINT `fk_grup_id_grup_pengguna` FOREIGN KEY (`id_grup_pengguna`) REFERENCES `grup_pengguna` (`id_grup_pengguna`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

LOCK TABLES `grup` WRITE;
/*!40000 ALTER TABLE `grup` DISABLE KEYS */;

INSERT INTO `grup` (`id_grup`, `id_grup_pengguna`, `id_grup_akses`)
VALUES
	(1,1,1),
	(2,2,1),
	(3,2,2),
	(4,3,1),
	(5,3,2),
	(6,3,3);

/*!40000 ALTER TABLE `grup` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table grup_akses
# ------------------------------------------------------------

DROP TABLE IF EXISTS `grup_akses`;

CREATE TABLE `grup_akses` (
  `id_grup_akses` int(11) NOT NULL AUTO_INCREMENT,
  `nama_grup_akses` varchar(255) NOT NULL,
  PRIMARY KEY (`id_grup_akses`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

LOCK TABLES `grup_akses` WRITE;
/*!40000 ALTER TABLE `grup_akses` DISABLE KEYS */;

INSERT INTO `grup_akses` (`id_grup_akses`, `nama_grup_akses`)
VALUES
	(1,'Belanja'),
	(2,'Menyetok barang'),
	(3,'Menjual kelebihan/sisa barang');

/*!40000 ALTER TABLE `grup_akses` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table grup_pengguna
# ------------------------------------------------------------

DROP TABLE IF EXISTS `grup_pengguna`;

CREATE TABLE `grup_pengguna` (
  `id_grup_pengguna` int(11) NOT NULL AUTO_INCREMENT,
  `nama_grup_pengguna` varchar(255) NOT NULL,
  `keterangan` varchar(255) NOT NULL,
  PRIMARY KEY (`id_grup_pengguna`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

LOCK TABLES `grup_pengguna` WRITE;
/*!40000 ALTER TABLE `grup_pengguna` DISABLE KEYS */;

INSERT INTO `grup_pengguna` (`id_grup_pengguna`, `nama_grup_pengguna`, `keterangan`)
VALUES
	(1,'Agen','Tingkat yang paling rendah (agen biasa)'),
	(2,'Master','Tingkat yang lebih tinggi dari agen'),
	(3,'Super Master','Tingkat yang paling tinggi');

/*!40000 ALTER TABLE `grup_pengguna` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table pengguna
# ------------------------------------------------------------

DROP TABLE IF EXISTS `pengguna`;

CREATE TABLE `pengguna` (
  `id_pengguna` int(11) NOT NULL AUTO_INCREMENT,
  `nama_pengguna` varchar(255) NOT NULL,
  `email_pengguna` varchar(255) NOT NULL,
  `deskripsi` text NOT NULL,
  `tanggal_dibuat` datetime NOT NULL,
  `tanggal_diperbaharui` datetime NOT NULL,
  `id_admin` int(11) NOT NULL,
  `id_grup_pengguna` int(11) NOT NULL,
  PRIMARY KEY (`id_pengguna`),
  UNIQUE KEY `email_pengguna` (`email_pengguna`),
  KEY `id_admin` (`id_admin`),
  KEY `id_grup_pengguna` (`id_grup_pengguna`),
  CONSTRAINT `fk_id_admin` FOREIGN KEY (`id_admin`) REFERENCES `admin` (`id_admin`),
  CONSTRAINT `fk_pengguna_id_grup_pengguna` FOREIGN KEY (`id_grup_pengguna`) REFERENCES `grup_pengguna` (`id_grup_pengguna`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

LOCK TABLES `pengguna` WRITE;
/*!40000 ALTER TABLE `pengguna` DISABLE KEYS */;

INSERT INTO `pengguna` (`id_pengguna`, `nama_pengguna`, `email_pengguna`, `deskripsi`, `tanggal_dibuat`, `tanggal_diperbaharui`, `id_admin`, `id_grup_pengguna`)
VALUES
	(1,'Wilson Brain','iconbrain96@gmail.com','Pengguna Kudo yang aktif. Wilson Brain adalah salah satu pengguna yang aktif. Wilson Brain telah bergabung dengan Kudo sejak 2017.','2017-11-15 05:56:14','2017-11-19 14:06:08',1,1),
	(2,'Kevin','kevin@gmail','Pengguna Kudo yang pasif. Kevin cukup pasif dan tidak up to date dengan sistem Kudo.','2017-11-15 05:56:14','2017-11-18 22:38:09',1,3),
	(3,'Andy','andy@gmail.com','Pengguna Kudo yang aktif. Andy adalah salah satu pengguna yang aktif. Andy telah bergabung dengan Kudo sejak 2017.','2017-11-15 05:56:14','2017-11-15 05:56:14',1,1),
	(4,'Dicky','dicky@gmail.com','Pengguna Kudo yang aktif. Dicky adalah salah satu pengguna yang aktif. Dicky telah bergabung dengan Kudo sejak 2017.','2017-11-15 05:56:14','2017-11-19 13:59:43',1,3),
	(5,'Roderick','roderick@gmail.com','Pengguna Kudo yang aktif. Roderick adalah salah satu pengguna yang aktif. Roderick telah bergabung dengan Kudo sejak 2017.','2017-11-15 05:56:14','2017-11-15 05:56:14',1,3),
	(6,'Santo','santo@gmail.com','Pengguna Kudo yang pasif. Santo cukup pasif dan tidak up to date dengan sistem Kudo.','2017-11-15 05:56:14','2017-11-15 05:56:14',1,3),
	(7,'Alvin','alvin@gmail.com','Pengguna Kudo yang pasif. Alvin cukup pasif dan tidak up to date dengan sistem Kudo.','2017-11-15 05:56:14','2017-11-15 05:56:14',1,2),
	(8,'Ricky','ricky@gmail.com','Pengguna Kudo yang pasif. Ricky cukup pasif dan tidak up to date dengan sistem Kudo.','2017-11-15 05:56:14','2017-11-15 05:56:14',1,2),
	(9,'Bodhi','bodhi@gmail.com','Pengguna Kudo yang pasif. Bodhi cukup pasif dan tidak up to date dengan sistem Kudo.','2017-11-15 05:56:14','2017-11-15 05:56:14',1,2),
	(10,'Hendrix','hendrix@gmail.com','Pengguna Kudo yang aktif. Hendrix adalah salah satu pengguna yang aktif. Hendrix telah bergabung dengan Kudo sejak 2017.','2017-11-15 05:56:14','2017-11-15 05:56:14',1,1),
	(11,'Richard','richard@gmail.com','Pengguna Kudo yang aktif. Richard adalah salah satu pengguna yang aktif. Richard telah bergabung dengan Kudo sejak 2017.','2017-11-15 05:56:14','2017-11-15 05:56:14',1,1);

/*!40000 ALTER TABLE `pengguna` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
