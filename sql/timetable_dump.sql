--
-- PostgreSQL database dump
--

-- Dumped from database version 10.21 (Ubuntu 10.21-0ubuntu0.18.04.1)
-- Dumped by pg_dump version 10.21 (Ubuntu 10.21-0ubuntu0.18.04.1)

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
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: book; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.book (
    bid integer NOT NULL,
    bno smallint,
    btype smallint,
    bname character varying(20) NOT NULL
);


ALTER TABLE public.book OWNER TO postgres;

--
-- Name: day; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.day (
    dayid integer NOT NULL,
    dayname character varying(20) NOT NULL
);


ALTER TABLE public.day OWNER TO postgres;

--
-- Name: day_uniform; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.day_uniform (
    dayid integer,
    uid integer
);


ALTER TABLE public.day_uniform OWNER TO postgres;

--
-- Name: period; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.period (
    pid integer NOT NULL,
    pstarttime time without time zone NOT NULL,
    pendtime time without time zone NOT NULL,
    ptype smallint,
    pdescr character varying(20)
);


ALTER TABLE public.period OWNER TO postgres;

--
-- Name: sub_book; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sub_book (
    subid integer,
    bid integer
);


ALTER TABLE public.sub_book OWNER TO postgres;

--
-- Name: subject; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.subject (
    subid integer NOT NULL,
    subname character varying(20) NOT NULL
);


ALTER TABLE public.subject OWNER TO postgres;

--
-- Name: tt; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tt (
    ttid integer NOT NULL,
    dayid integer,
    pid integer,
    subid integer
);


ALTER TABLE public.tt OWNER TO postgres;

--
-- Name: uniform; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.uniform (
    uid integer NOT NULL,
    utype smallint,
    uname character varying(10) NOT NULL
);


ALTER TABLE public.uniform OWNER TO postgres;

--
-- Data for Name: book; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.book (bid, bno, btype, bname) FROM stdin;
1	0	0	Notes
2	1	0	Rough book
3	2	0	Dication
4	3	0	English
5	4	0	Math
6	5	0	EVS
7	6	0	Hindi
8	7	0	Marathi
9	8	1	English textbook
10	9	1	Math textbook
11	10	1	Hindi textbook
12	11	1	Marathi textbook
\.


--
-- Data for Name: day; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.day (dayid, dayname) FROM stdin;
1	Monday
2	Tuesday
3	Wednesday
4	Thursday
5	Friday
\.


--
-- Data for Name: day_uniform; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.day_uniform (dayid, uid) FROM stdin;
\.


--
-- Data for Name: period; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.period (pid, pstarttime, pendtime, ptype, pdescr) FROM stdin;
\.


--
-- Data for Name: sub_book; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.sub_book (subid, bid) FROM stdin;
\.


--
-- Data for Name: subject; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.subject (subid, subname) FROM stdin;
1	ENGLISH
2	EVS
3	MATH
4	HINDI
5	MARATHI
6	DRAWING
7	COMPUTER
8	VALUE EDUCATION
9	PE
10	GK
11	STORY TELLING
12	MUSIC
13	SUPW
14	ZERO HOUR
15	LIBRARY
\.


--
-- Data for Name: tt; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.tt (ttid, dayid, pid, subid) FROM stdin;
\.


--
-- Data for Name: uniform; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.uniform (uid, utype, uname) FROM stdin;
1	1	Regular
2	2	PE uniform
\.


--
-- Name: book book_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.book
    ADD CONSTRAINT book_pkey PRIMARY KEY (bid);


--
-- Name: day day_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.day
    ADD CONSTRAINT day_pkey PRIMARY KEY (dayid);


--
-- Name: period period_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.period
    ADD CONSTRAINT period_pkey PRIMARY KEY (pid);


--
-- Name: subject subject_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.subject
    ADD CONSTRAINT subject_pkey PRIMARY KEY (subid);


--
-- Name: tt tt_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tt
    ADD CONSTRAINT tt_pkey PRIMARY KEY (ttid);


--
-- Name: uniform uniform_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.uniform
    ADD CONSTRAINT uniform_pkey PRIMARY KEY (uid);


--
-- Name: sub_book fk_bid; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sub_book
    ADD CONSTRAINT fk_bid FOREIGN KEY (bid) REFERENCES public.book(bid);


--
-- Name: day_uniform fk_dayid; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.day_uniform
    ADD CONSTRAINT fk_dayid FOREIGN KEY (dayid) REFERENCES public.day(dayid);


--
-- Name: tt fk_did; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tt
    ADD CONSTRAINT fk_did FOREIGN KEY (dayid) REFERENCES public.day(dayid);


--
-- Name: tt fk_pid; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tt
    ADD CONSTRAINT fk_pid FOREIGN KEY (pid) REFERENCES public.period(pid);


--
-- Name: tt fk_sid; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tt
    ADD CONSTRAINT fk_sid FOREIGN KEY (subid) REFERENCES public.subject(subid);


--
-- Name: sub_book fk_subbid; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sub_book
    ADD CONSTRAINT fk_subbid FOREIGN KEY (subid) REFERENCES public.subject(subid);


--
-- Name: day_uniform fk_uid; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.day_uniform
    ADD CONSTRAINT fk_uid FOREIGN KEY (uid) REFERENCES public.uniform(uid);


--
-- PostgreSQL database dump complete
--

