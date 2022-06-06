--
-- PostgreSQL database dump
--

-- Dumped from database version 12.10
-- Dumped by pg_dump version 12.10

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: balances; Type: TABLE; Schema: public; Owner: marsel
--

CREATE TABLE public.balances (
    id integer NOT NULL,
    current numeric(32,2) DEFAULT 0,
    withdrawn numeric(32,2) DEFAULT 0,
    user_id integer NOT NULL,
    CONSTRAINT balances_current_check CHECK ((current >= (0)::numeric)),
    CONSTRAINT balances_withdrawn_check CHECK ((withdrawn >= (0)::numeric))
);


ALTER TABLE public.balances OWNER TO marsel;

--
-- Name: balances_id_seq; Type: SEQUENCE; Schema: public; Owner: marsel
--

CREATE SEQUENCE public.balances_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.balances_id_seq OWNER TO marsel;

--
-- Name: balances_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: marsel
--

ALTER SEQUENCE public.balances_id_seq OWNED BY public.balances.id;


--
-- Name: balances_user_id_seq; Type: SEQUENCE; Schema: public; Owner: marsel
--

CREATE SEQUENCE public.balances_user_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.balances_user_id_seq OWNER TO marsel;

--
-- Name: balances_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: marsel
--

ALTER SEQUENCE public.balances_user_id_seq OWNED BY public.balances.user_id;


--
-- Name: orders; Type: TABLE; Schema: public; Owner: marsel
--

CREATE TABLE public.orders (
    id integer NOT NULL,
    number character varying,
    status character varying NOT NULL,
    accrual numeric(32,2),
    uploaded_at timestamp with time zone NOT NULL,
    user_id integer NOT NULL
);


ALTER TABLE public.orders OWNER TO marsel;

--
-- Name: orders_id_seq; Type: SEQUENCE; Schema: public; Owner: marsel
--

CREATE SEQUENCE public.orders_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.orders_id_seq OWNER TO marsel;

--
-- Name: orders_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: marsel
--

ALTER SEQUENCE public.orders_id_seq OWNED BY public.orders.id;


--
-- Name: orders_user_id_seq; Type: SEQUENCE; Schema: public; Owner: marsel
--

CREATE SEQUENCE public.orders_user_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.orders_user_id_seq OWNER TO marsel;

--
-- Name: orders_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: marsel
--

ALTER SEQUENCE public.orders_user_id_seq OWNED BY public.orders.user_id;


--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: marsel
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO marsel;

--
-- Name: users; Type: TABLE; Schema: public; Owner: marsel
--

CREATE TABLE public.users (
    id integer NOT NULL,
    login character varying,
    passwordhash character varying
);


ALTER TABLE public.users OWNER TO marsel;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: marsel
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO marsel;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: marsel
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: withdrawals; Type: TABLE; Schema: public; Owner: marsel
--

CREATE TABLE public.withdrawals (
    id integer NOT NULL,
    "order" character varying,
    sum numeric(32,2) DEFAULT 0,
    processed_at timestamp with time zone NOT NULL,
    user_id integer NOT NULL
);


ALTER TABLE public.withdrawals OWNER TO marsel;

--
-- Name: withdrawals_id_seq; Type: SEQUENCE; Schema: public; Owner: marsel
--

CREATE SEQUENCE public.withdrawals_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.withdrawals_id_seq OWNER TO marsel;

--
-- Name: withdrawals_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: marsel
--

ALTER SEQUENCE public.withdrawals_id_seq OWNED BY public.withdrawals.id;


--
-- Name: withdrawals_user_id_seq; Type: SEQUENCE; Schema: public; Owner: marsel
--

CREATE SEQUENCE public.withdrawals_user_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.withdrawals_user_id_seq OWNER TO marsel;

--
-- Name: withdrawals_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: marsel
--

ALTER SEQUENCE public.withdrawals_user_id_seq OWNED BY public.withdrawals.user_id;


--
-- Name: balances id; Type: DEFAULT; Schema: public; Owner: marsel
--

ALTER TABLE ONLY public.balances ALTER COLUMN id SET DEFAULT nextval('public.balances_id_seq'::regclass);


--
-- Name: balances user_id; Type: DEFAULT; Schema: public; Owner: marsel
--

ALTER TABLE ONLY public.balances ALTER COLUMN user_id SET DEFAULT nextval('public.balances_user_id_seq'::regclass);


--
-- Name: orders id; Type: DEFAULT; Schema: public; Owner: marsel
--

ALTER TABLE ONLY public.orders ALTER COLUMN id SET DEFAULT nextval('public.orders_id_seq'::regclass);


--
-- Name: orders user_id; Type: DEFAULT; Schema: public; Owner: marsel
--

ALTER TABLE ONLY public.orders ALTER COLUMN user_id SET DEFAULT nextval('public.orders_user_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: marsel
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: withdrawals id; Type: DEFAULT; Schema: public; Owner: marsel
--

ALTER TABLE ONLY public.withdrawals ALTER COLUMN id SET DEFAULT nextval('public.withdrawals_id_seq'::regclass);


--
-- Name: withdrawals user_id; Type: DEFAULT; Schema: public; Owner: marsel
--

ALTER TABLE ONLY public.withdrawals ALTER COLUMN user_id SET DEFAULT nextval('public.withdrawals_user_id_seq'::regclass);


--
-- Name: balances balances_pkey; Type: CONSTRAINT; Schema: public; Owner: marsel
--

ALTER TABLE ONLY public.balances
    ADD CONSTRAINT balances_pkey PRIMARY KEY (id);


--
-- Name: balances balances_user_id_key; Type: CONSTRAINT; Schema: public; Owner: marsel
--

ALTER TABLE ONLY public.balances
    ADD CONSTRAINT balances_user_id_key UNIQUE (user_id);


--
-- Name: orders orders_number_key; Type: CONSTRAINT; Schema: public; Owner: marsel
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_number_key UNIQUE (number);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: marsel
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- Name: users users_login_key; Type: CONSTRAINT; Schema: public; Owner: marsel
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_login_key UNIQUE (login);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: marsel
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: withdrawals withdrawals_pkey; Type: CONSTRAINT; Schema: public; Owner: marsel
--

ALTER TABLE ONLY public.withdrawals
    ADD CONSTRAINT withdrawals_pkey PRIMARY KEY (id);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: marsel
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- Name: orders fk_user; Type: FK CONSTRAINT; Schema: public; Owner: marsel
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: balances fk_user; Type: FK CONSTRAINT; Schema: public; Owner: marsel
--

ALTER TABLE ONLY public.balances
    ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

