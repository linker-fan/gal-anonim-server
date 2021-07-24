create table users (
	id serial primary key,
	username varchar(30) unique not null,
	passwordHash text not null,
    isAdmin boolean not null, 
	created timestamp not null,
	updated timestamp not null
);

create table room_members (
	id serial primary key,
	userID integer references users(id) on delete set null
); 

create table rooms (
	id text primary key not null,
	roomName varchar(50) not null,
	roomPasswordHash text not null, 
	ownerID integer not null references users(id) on delete set null,
	roomMembersID integer not null references room_members(id) on delete set null
);