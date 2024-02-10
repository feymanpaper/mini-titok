use titok_relation;

DROP TABLE IF EXISTS t_film;
CREATE TABLE t_film(
    id INTEGER AUTO_INCREMENT PRIMARY KEY,
    score DECIMAL(2,1) NOT NULL,
    file_name VARCHAR(20) NOT NULL ,
    introduction VARCHAR(30) NOT NULL
);

DROP PROCEDURE IF EXISTS proc_t_film;
DELIMITER $$
CREATE PROCEDURE proc_t_film(total int)
BEGIN DECLARE i INTEGER DEFAULT 1;
    START TRANSACTION;
    WHILE i<=total DO
        INSERT t_film(score, file_name, introduction)
        VALUES (
            ROUND(RAND()*9.9, 1),
            CONCAT('file_name', i),
            CONCAT('introduction', i));
        SET i=i+1;
    END WHILE;
    COMMIT;
END $$
DELIMITER ;

CALL proc_t_film(10000000);

ALTER TABLE t_film ADD INDEX idx_score (score);

# 100w数据浅分页 380ms, 加上score索引后为15-30ms
# 1000w数据2635ms, 加上score索引为15ms
SELECT score, file_name FROM t_film ORDER BY score DESC LIMIT 5, 20;

# 100w数据深分页 587ms, 加上score索引后为500-703ms, 强制索引需要3739ms
# ORDER BY和SELECT字段加联合索引 200ms
# ORDER BY加索引+手动回表 200ms
# 1000w数据深分页4531ms, 加上索引为4544ms
# ORDER BY加索引+手动回表 200ms
SELECT score, file_name FROM t_film ORDER BY score DESC LIMIT 900000, 20;

# 法1, ORDER BY和SELECT字段加联合索引
ALTER TABLE t_film ADD INDEX idx_score_fname (score, file_name);
SELECT score, file_name FROM t_film ORDER BY score DESC LIMIT 900000, 20;

# 法2, ORDER BY加索引+手动回表
SELECT score, file_name FROM t_film a JOIN
    (SELECT id FROM t_film ORDER BY score DESC LIMIT 900000, 20) b
ON a.id = b.id;