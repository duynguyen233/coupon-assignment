-- Create "coupons" table
CREATE TABLE `coupons` (
  `coupon_code` varchar(255) NOT NULL,
  `title` varchar(255) NOT NULL,
  `description` text NOT NULL,
  `coupon_type` enum('fixed','percentage') NOT NULL,
  `usage` text NOT NULL,
  `expired_at` datetime NOT NULL,
  `coupon_value` decimal(10,2) NOT NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  PRIMARY KEY (`coupon_code`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
