

/*$ sudo su postgres -c psql*/
/*jullebulle pw*/
/* \l - show dbs */
/* \dt schemes */

/*change to random db before we drop it*/

/* \i /home/jb/go/src/github.com/julleb/ServerMonitor/database.sql */

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
memory_usage int,
memory_total int,
total_memory int
);


create table has(
ip varchar(16) REFERENCES server (ip),
info_id int REFERENCES information (info_id),
PRIMARY KEY (ip,info_id)
);


/*select * from server NATURAL JOIN information ON id =  */


/*serial primary key --> increment the primary key in each insert*/
