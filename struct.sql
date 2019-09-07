create database hma
	with owner postgres;

create table if not exists band_links
(
	id serial not null
		constraint band_links_pkey
			primary key,
	created_at timestamp(6) with time zone,
	updated_at timestamp(6) with time zone,
	deleted_at timestamp(6) with time zone,
	name text,
	url text
);

alter table band_links owner to hma;

create index if not exists idx_band_links_deleted_at
	on band_links (deleted_at);

create index if not exists url_idx
	on band_links (url);

create table if not exists bands_genres
(
	band_id integer not null,
	genre_id integer not null
);

alter table bands_genres owner to hma;

create table if not exists countries
(
	id serial not null
		constraint countries_pk
			primary key,
	name varchar not null,
	created_at timestamp(6) with time zone,
	updated_at timestamp(6) with time zone,
	deleted_at timestamp(6) with time zone
);

alter table countries owner to hma;

create unique index if not exists countries_id_uindex
	on countries (id);

create unique index if not exists countries_name_uindex
	on countries (name);

create index if not exists idx_countries_deleted_at
	on countries (deleted_at);

create table if not exists artists
(
	id serial not null
		constraint artists_pk
			primary key,
	real_name varchar,
	country_id integer
		constraint artists_country_fk
			references countries
			on update set null on delete set null,
	born_year integer,
	gender varchar,
	image integer
);

alter table artists owner to hma;

create unique index if not exists artists_id_uindex
	on artists (id);

create table if not exists genres
(
	id serial not null
		constraint genres_pk
			primary key,
	name varchar not null,
	created_at timestamp(6) with time zone,
	updated_at timestamp(6) with time zone,
	deleted_at timestamp(6) with time zone
);

alter table genres owner to hma;

create unique index if not exists genres_id_uindex
	on genres (id);

create unique index if not exists genres_name_uindex
	on genres (name);

create index if not exists idx_genres_deleted_at
	on genres (deleted_at);

create table if not exists labels
(
	id serial not null
		constraint labels_pk
			primary key,
	name varchar not null,
	created_at timestamp(6) with time zone,
	updated_at timestamp(6) with time zone,
	deleted_at timestamp(6) with time zone
);

alter table labels owner to hma;

create index if not exists idx_labels_deleted_at
	on labels (deleted_at);

create unique index if not exists labels_id_uindex
	on labels (id);

create unique index if not exists labels_name_uindex
	on labels (name);

create table if not exists bands
(
	id serial not null
		constraint bands_pk
			primary key,
	name varchar not null,
	status varchar not null,
	country_id integer
		constraint bands_country_fk
			references countries
			on update set null on delete set null,
	formed_in integer,
	years_active varchar,
	label_id integer
		constraint bands_label_fk
			references labels
			on update set null on delete set null,
	description varchar,
	image_logo varchar,
	created_at timestamp(6) with time zone,
	updated_at timestamp(6) with time zone,
	deleted_at timestamp(6) with time zone,
	platform_id varchar not null,
	image_band varchar
);

comment on column bands.formed_in is 'year';

comment on column bands.years_active is '2000-2010';

alter table bands owner to hma;

create unique index if not exists bands_id_uindex
	on bands (id);

create index if not exists idx_bands_deleted_at
	on bands (deleted_at, deleted_at, deleted_at);

create unique index if not exists bands_platform_id_uindex
	on bands (platform_id);

create table if not exists albums
(
	id serial not null
		constraint albums_pk
			primary key,
	band_id integer
		constraint albums_band_fk
			references bands
			on update cascade on delete cascade,
	type varchar,
	name varchar not null,
	year integer,
	release_date timestamp(6),
	image varchar,
	total_time varchar,
	label_id integer
		constraint albums_label_fk
			references labels
			on update set null on delete set null,
	platform_id integer not null,
	created_at timestamp(6) with time zone,
	updated_at timestamp(6) with time zone,
	deleted_at timestamp(6) with time zone
);

alter table albums owner to hma;

create unique index if not exists albums_id_uindex
	on albums (id);

create unique index if not exists albums_platform_id_uindex
	on albums (platform_id);

create index if not exists idx_albums_deleted_at
	on albums (deleted_at);

create table if not exists albums_members
(
	id serial not null
		constraint albums_members_pk
			primary key,
	album_id integer
		constraint albums_members_albums_fk
			references albums
			on update set null on delete set null,
	artist_id integer
		constraint albums_members_artists_fk
			references artists
			on update set null on delete set null,
	position varchar
);

alter table albums_members owner to hma;

create unique index if not exists albums_members_id_uindex
	on albums_members (id);

create table if not exists bands_members
(
	id serial not null
		constraint bands_members_pk
			primary key,
	band_id integer
		constraint bands_members_bands_fk
			references bands
			on update restrict on delete restrict,
	artist_id integer
		constraint bands_members_artists_fk
			references artists
			on update restrict on delete restrict,
	position varchar,
	type varchar
);

comment on column bands_members.type is '(current / past)';

alter table bands_members owner to hma;

create unique index if not exists bands_members_id_uindex
	on bands_members (id);

create table if not exists lyrical_themes
(
	id serial not null
		constraint lyrichal_themes_pk
			primary key,
	name varchar not null,
	created_at timestamp(6) with time zone,
	updated_at timestamp(6) with time zone,
	deleted_at timestamp(6) with time zone
);

alter table lyrical_themes owner to hma;

create index if not exists idx_lyrical_themes_deleted_at
	on lyrical_themes (deleted_at);

create unique index if not exists lyrichal_themes_id_uindex
	on lyrical_themes (id);

create unique index if not exists lyrichal_themes_name_uindex
	on lyrical_themes (name);

create table if not exists bands_lyrical_themes
(
	band_id integer not null
		constraint bands_lyrical_themes_band_fk
			references bands
			on update cascade on delete cascade,
	lyrical_theme_id integer not null
		constraint bands_lyrical_themes_lyrical_theme_id_fk
			references lyrical_themes
			on update cascade on delete cascade
);

alter table bands_lyrical_themes owner to hma;

create table if not exists reviews
(
	id serial not null
		constraint reviews_pk
			primary key,
	album_id integer
		constraint reviews_albums_fk
			references albums
			on update set null on delete set null,
	rating varchar,
	date timestamp(6),
	text varchar
);

alter table reviews owner to hma;

create unique index if not exists reviews_id_uindex
	on reviews (id);

create table if not exists similar_bands
(
	id serial not null
		constraint similar_bands_pk
			primary key,
	band_id integer
		constraint similar_bands_bands_fk
			references bands,
	similar_band_id integer
		constraint similar_bands_bands_similar_fk
			references bands,
	score integer
);

alter table similar_bands owner to hma;

create unique index if not exists similar_bands_id_uindex
	on similar_bands (id);

create table if not exists songs
(
	id serial not null
		constraint songs_pk
			primary key,
	name varchar,
	band_id integer
		constraint songs_bands_id_fk
			references bands,
	album_id integer
		constraint songs_albums_fk
			references albums,
	position integer,
	time varchar,
	lyrics varchar,
	platform_id integer not null,
	created_at timestamp(6) with time zone,
	updated_at timestamp(6) with time zone,
	deleted_at timestamp(6) with time zone
);

alter table songs owner to hma;

create index if not exists idx_songs_deleted_at
	on songs (deleted_at);

create unique index if not exists songs_id_uindex
	on songs (id);

create unique index if not exists songs_platform_id_uindex
	on songs (platform_id);

create table if not exists latest_band_updates
(
	id serial not null
		constraint latest_band_update_pk
			primary key,
	band_id integer not null
		constraint latest_band_update_bands_id_fk
			references bands,
	created_at timestamp,
	updated_at timestamp
);

alter table latest_band_updates owner to hma;

create unique index if not exists latest_band_update_id_uindex
	on latest_band_updates (id);

create table if not exists upcoming_albums
(
	id serial not null
		constraint upcoming_albums_pk
			primary key,
	album_id integer not null
		constraint upcoming_albums_albums_id_fk
			references albums,
	created_at timestamp not null,
	updated_at timestamp
);

alter table upcoming_albums owner to hma;

create unique index if not exists upcoming_albums_id_uindex
	on upcoming_albums (id);

