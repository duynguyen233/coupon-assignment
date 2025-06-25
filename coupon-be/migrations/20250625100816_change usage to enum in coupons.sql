-- Modify "coupons" table
ALTER TABLE `coupons` MODIFY COLUMN `usage` enum('manual','auto') NOT NULL;
