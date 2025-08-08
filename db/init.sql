-- Create tables in public schema of mt_test

-- DROP TABLE IF EXISTS in proper order to avoid FK conflicts
DROP TABLE IF EXISTS public.follows;
DROP TABLE IF EXISTS public.posts;
DROP TABLE IF EXISTS public.users;

CREATE TABLE public.users (
	id int4 NOT NULL GENERATED ALWAYS AS IDENTITY,
	"name" varchar(100) NOT NULL,
	email varchar(255) NOT NULL,
	"password" varchar(255) NOT NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT users_email_key UNIQUE (email),
	CONSTRAINT users_pkey PRIMARY KEY (id)
);

CREATE TABLE public.follows (
	id int4 NOT NULL GENERATED ALWAYS AS IDENTITY,
	follower_id int4 NOT NULL,
	followed_id int4 NOT NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT follows_check CHECK ((follower_id <> followed_id)),
	CONSTRAINT follows_pkey PRIMARY KEY (id),
	CONSTRAINT unique_follow UNIQUE (follower_id, followed_id),
	CONSTRAINT fk_followed FOREIGN KEY (followed_id) REFERENCES public.users(id) ON DELETE CASCADE,
	CONSTRAINT fk_follower FOREIGN KEY (follower_id) REFERENCES public.users(id) ON DELETE CASCADE
);

CREATE TABLE public.posts (
	id int4 NOT NULL GENERATED ALWAYS AS IDENTITY,
	title varchar(200) NOT NULL,
	"content" varchar(5000) NOT NULL,
	user_id int4 NOT NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT posts_pkey PRIMARY KEY (id),
	CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE
);
