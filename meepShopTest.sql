--
-- PostgreSQL database dump
--

-- Dumped from database version 14.4
-- Dumped by pg_dump version 14.4

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: users; Type: TABLE; Schema: public; Owner: test
--

CREATE TABLE public.users (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    first_name text,
    last_name text,
    display_name text,
    account_number text,
    deposit numeric
);


ALTER TABLE public.users OWNER TO test;

--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: test
--

COPY public.users (id, created_at, updated_at, deleted_at, first_name, last_name, display_name, account_number, deposit) FROM stdin;
ea900cd3-6825-48f0-98a8-37b50bc92903	2024-02-21 23:21:55.397227+08	2024-02-21 23:27:12.929578+08	\N	yanxian	li	yanxianli	999405903669680	100
e0e04293-691b-48d3-9838-6534d1f214c6	2024-02-21 23:22:11.224615+08	2024-02-21 23:27:12.925651+08	\N	billy	li	billyli	870621047041519	50
\.


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: test
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_users_account_number; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX idx_users_account_number ON public.users USING btree (account_number);


--
-- Name: idx_users_id; Type: INDEX; Schema: public; Owner: test
--

CREATE UNIQUE INDEX idx_users_id ON public.users USING btree (id);


--
-- PostgreSQL database dump complete
--

