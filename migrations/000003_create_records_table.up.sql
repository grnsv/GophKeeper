CREATE TABLE public.records (
	id uuid NOT NULL,
	user_id uuid NOT NULL,
	type public.record_type NOT NULL,
	data bytea NOT NULL,
	metadata jsonb DEFAULT '{}'::jsonb NOT NULL,
	version int8 DEFAULT 0 NOT NULL,
	CONSTRAINT records_pk PRIMARY KEY (id, user_id),
	CONSTRAINT records_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE
);

CREATE INDEX records_user_id_idx ON public.records USING btree (user_id);
