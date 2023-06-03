CREATE TABLE `urls`
(
    `id`        int(11)      NOT NULL AUTO_INCREMENT,
    `code`      varchar(255) NOT NULL COMMENT 'code 短码',
    `data`      varchar(512) NOT NULL COMMENT '数据',
    `exprie`    datetime     DEFAULT NULL COMMENT '过期时间',
    `status`    int(2)       NOT NULL,
    `count`     int(11)      DEFAULT 0,
    `create_by` varchar(255) DEFAULT NULL,
    `create_at` datetime     DEFAULT NULL,
    `update_at` datetime     DEFAULT NULL ON UPDATE current_timestamp(),
    `delete_at` datetime     DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uqi_code` (`code`) USING BTREE COMMENT '唯一索引'
) ENGINE = InnoDB
  AUTO_INCREMENT = 48
  DEFAULT CHARSET = utf8mb4

CREATE TABLE `data`
(
    `id`        int(11)      NOT NULL AUTO_INCREMENT,
    `code`      varchar(255) NOT NULL COMMENT 'code 短码',
    `data`      varchar(512) NOT NULL COMMENT '数据',
    `exprie`    datetime     DEFAULT NULL COMMENT '过期时间',
    `status`    int(2)       NOT NULL,
    `count`     int(11)      DEFAULT 0,
    `create_by` varchar(255) DEFAULT NULL,
    `create_at` datetime     DEFAULT NULL,
    `update_at` datetime     DEFAULT NULL ON UPDATE current_timestamp(),
    `delete_at` datetime     DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uqi_code` (`code`) USING BTREE COMMENT '唯一索引'
) ENGINE = InnoDB
  AUTO_INCREMENT = 48
  DEFAULT CHARSET = utf8mb4

CREATE TABLE `permanent_urls`
(
    `id`        int(11)      NOT NULL AUTO_INCREMENT,
    `code`      varchar(255) NOT NULL COMMENT 'code 短码',
    `data`      varchar(512) NOT NULL COMMENT '数据',
    `exprie`    datetime     DEFAULT NULL COMMENT '过期时间',
    `status`    int(2)       NOT NULL,
    `count`     int(11)      DEFAULT 0,
    `create_by` varchar(255) DEFAULT NULL,
    `create_at` datetime     DEFAULT NULL,
    `update_at` datetime     DEFAULT NULL ON UPDATE current_timestamp(),
    `delete_at` datetime     DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uqi_code` (`code`) USING BTREE COMMENT '唯一索引'
) ENGINE = InnoDB
  AUTO_INCREMENT = 48
  DEFAULT CHARSET = utf8mb4

CREATE TABLE `permanent_data`
(
    `id`        int(11)      NOT NULL AUTO_INCREMENT,
    `code`      varchar(255) NOT NULL COMMENT 'code 短码',
    `data`      varchar(512) NOT NULL COMMENT '数据',
    `exprie`    datetime     DEFAULT NULL COMMENT '过期时间',
    `status`    int(2)       NOT NULL,
    `count`     int(11)      DEFAULT 0,
    `create_by` varchar(255) DEFAULT NULL,
    `create_at` datetime     DEFAULT NULL,
    `update_at` datetime     DEFAULT NULL ON UPDATE current_timestamp(),
    `delete_at` datetime     DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uqi_code` (`code`) USING BTREE COMMENT '唯一索引'
) ENGINE = InnoDB
  AUTO_INCREMENT = 48
  DEFAULT CHARSET = utf8mb4

