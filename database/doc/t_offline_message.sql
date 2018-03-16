/*
set names utf8;
create TABLE IF NOT EXISTS t_offline_message
(
    uid                 int unsigned not NULL comment 'uid',
    id                  big int unsigned not NULL comment 'id',
    send_time           Datetime not NULL comment 'send_time',
    fuid                int unsigned not NULL comment 'fuid',
    funame              VARCHAR(16) not NULL comment 'funame',
    ffigure             int unsigned not NULL comment 'ffigure',
    et_data             blob not NULL comment 'msg data',
    PRIMARY KEY (uid,id)
)engine = InnoDb default charset utf8;
*/