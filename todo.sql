CREATE TABLE `tbl_todo`
(
    `todo_id`      varchar(50) NOT NULL,
    `todo`         varchar(200) DEFAULT NULL,
    `is_completed` varchar(5),
    `created_at`   varchar(30),
    PRIMARY KEY (`todo_id`)
)