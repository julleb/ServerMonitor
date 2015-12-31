/*change to random db before we drop it*/
\c julle; 
drop database servermonitor;
create database servermonitor;

\c servermonitor;

create table server (
ip varchar(16) PRIMARY KEY
);


create table information(
info_id SERIAL PRIMARY KEY,
cpu_temp int,
cpu_load int,
memory_usage int,
memory_total int,
date Timestamp
);


create table has(
ip varchar(16) REFERENCES server (ip),
info_id int REFERENCES information (info_id),
PRIMARY KEY (ip,info_id)
);



