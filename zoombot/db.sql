CREATE TABLE `oauth_state` (
  `state` char(24) NOT NULL,
  `identifier` varchar(128) NOT NULL,
  `conv_id` char(64) NOT NULL,
  `msg_id` char(64) NOT NULL,
  `is_complete` boolean NOT NULL DEFAULT 0,
  PRIMARY KEY (`state`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `oauth` (
  `identifier` varchar(128) NOT NULL,
  `ctime` datetime NOT NULL,
  `mtime` datetime NOT NULL,
  `access_token` varchar(1024) NOT NULL,
  `token_type` varchar(64) NOT NULL,
  `refresh_token` varchar(1024) NOT NULL,
  `expiry` datetime NOT NULL,
  PRIMARY KEY (`identifier`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `user` (
  `user_id` varchar(128) NOT NULL,
  `account_id` varchar(128) NOT NULL,
  `identifier` varchar(128) NOT NULL,
  PRIMARY KEY (`user_id`, `account_id`, `identifier`),
  INDEX (`user_id`, `account_id`),
  FOREIGN KEY (`identifier`)
    REFERENCES oauth (`identifier`)
    ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
