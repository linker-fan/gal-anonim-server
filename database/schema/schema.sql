create table users (
	id serial primary key,
	username varchar(30) unique not null,
	passwordHash text not null,
	created timestamp not null,
	updated timestamp not null
);

create table rooms (
	id text primary key not null,
	ownerID integer not null references users(id) on delete set null
);

create table room_members (
	id serial primary key,
	roomID integer not null references rooms(id) on delete cascade,
	userID integer references users(id) on delete set null,
);
