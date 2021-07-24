create table users (
	id serial primary key,
	username varchar(30) unique not null,
	passwordHash text not null,
    isAdmin boolean not null, 
	created timestamp not null,
	updated timestamp not null
);

create table rooms (
	id serial primary key,
	uniqueRoomID text unique not null,
	roomName varchar(50) not null,
	passwordHash text not null,
	ownerID int references users(id) on delete set null,
	created timestamp not null,
	updated timestamp not null
);

create table messages (
	id serial primary key,
	roomID int not null references rooms(id) on delete set null,
	userID int not null references users(id) on delete set null,
	messageText text not null,
	created timestamp not null
);

create table members(
	id serial primary key,
	roomID int not null references rooms(id) on delete set null,
	userID int not null references users(id) on delete set null,
	joined timestamp not null
);

