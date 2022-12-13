-- CREATE DATABASE
use master
go

drop database if exists yachts
go

create database yachts
go

use yachts
go


-- CREATE TABLES
drop table if exists accounts
go

create table accounts
(
    [id] int primary key identity(1, 1),
    [email] varchar(255) not null,
    [name] varchar(255) not null,
    [password] binary(64) not null,
    [created] datetime not null default current_timestamp,
    [updated] datetime,
)
go

drop table if exists yachts
go

create table yachts
(
    [id] int primary key identity(1, 1),
    [name] varchar(255) not null,
    [type] varchar(255) not null,
    [price] float not null,
    [image] varchar(255) not null,
    [description] varchar(512),
    [created] datetime not null default current_timestamp,
    [updated] datetime,
)
go

drop table if exists reservations
go

create table reservations
(
    [id] int primary key identity(1, 1),
    [from] datetime not null,
    [to] datetime not null,
    [remarks] varchar(255) not null,
    [created_at] datetime not null default current_timestamp,
    [updated_at] datetime,
    [yacht_id] int foreign key references yachts([id]),
    [account_id] int foreign key references accounts([id]),
)
go

use yachts
select * from accounts