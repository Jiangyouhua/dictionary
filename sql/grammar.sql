-- 商品分类表
CREATE TABLE catalogue (
  `id` INT(10) UNSIGNED AUTO_INCREMENT,
  `parent_id` INT(10) UNSIGNED COMMENT 'self.id',
  `title` VARCHAR(255) COMMENT '目录标题',
  `info` VARCHAR(255) COMMENT '目录信息',
  `rank` INT(10) UNSIGNED COMMENT '排序，由小到大排序',
  `state` TINYINT(1) DEFAULT 1,
  `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `update_time` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`parent_id`) REFERENCES catalogue (`id`)
) COMMENT '目录';
CREATE TABLE clause (
  `id` INT(10) UNSIGNED AUTO_INCREMENT,
  `catalogue_id` INT(10) UNSIGNED COMMENT 'catalogue.id',
  `parent_id` INT(10) UNSIGNED COMMENT 'self.id',
  `title` VARCHAR(255) COMMENT '条目标题',
  `info` VARCHAR(255) COMMENT '条目信息',
  `image` VARCHAR(255) COMMENT '条目图片',
  `state` TINYINT(1) DEFAULT 1,
  `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `update_time` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`catalogue_id`) REFERENCES catalogue (`id`),
  FOREIGN KEY (`parent_id`) REFERENCES clause (`id`)
) COMMENT '条目';
CREATE TABLE example (
  `id` INT(10) UNSIGNED AUTO_INCREMENT,
  `clause_id` INT(10) UNSIGNED COMMENT 'clause.id',
  `english` VARCHAR(255) COMMENT '示例英文',
  `chinese` VARCHAR(255) COMMENT '示例中文',
  `emphasis` VARCHAR(255) COMMENT '示例重点',
  `state` TINYINT(1) DEFAULT 1,
  `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `update_time` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`clause_id`) REFERENCES clause (`id`)
) COMMENT '示例';
CREATE TABLE question (
  `id` INT(10) UNSIGNED AUTO_INCREMENT,
  `clause_id` INT(10) UNSIGNED COMMENT 'clause.id',
  `title` VARCHAR(255) COMMENT '题目',
  `option` VARCHAR(255) COMMENT '题目选项，多个使用|分割',
  `answer` VARCHAR(255) COMMENT '题目答案，多个使用|分割',
  `state` TINYINT(1) DEFAULT 1,
  `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `update_time` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`clause_id`) REFERENCES clause (`id`)
) COMMENT '题目';