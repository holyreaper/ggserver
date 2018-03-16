
/*set names utf8;
create TABLE IF NOT EXISTS t_user
(
    uid                 int unsigned not NULL comment 'uid',
    uname               VARCHAR(16) not NULL comment 'uname',
    create_time         Datetime not NULL comment 'create_time',
    last_login_time     Datetime not NULL comment 'last_login_time',
    last_logout_time    Datetime not NULL comment 'last_logout_time',
    exp                 int unsigned not NULL comment 'exp',
    level               int unsigned not NULL comment 'level',
    figure              int unsigned not NULL comment 'figure',
    et_data             blob not NULL comment 'extra data',
    PRIMARY KEY (uid)
)engine = InnoDb default charset utf8;
*/