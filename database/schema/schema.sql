create table users (
	id serial primary key,
	username varchar(30) unique not null,
	passwordHash text not null,
	created timestamp not null,
	updated timestamp not null
);
