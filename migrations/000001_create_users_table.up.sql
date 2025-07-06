CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid () NOT NULL,
    login text NOT NULL,
    password_hash bytea NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT users_login_key UNIQUE (login)
);
