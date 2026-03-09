--
-- PostgreSQL database dump
--

\restrict 8A7sykcbWn6skfC0SwHPsCQpGj3jW3cn8bfOLLXQLC0nMQMU1Sh6XMHbyqUJuFo

-- Dumped from database version 17.6
-- Dumped by pg_dump version 18.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA public;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: blacklisted_tokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.blacklisted_tokens (
    id bigint NOT NULL,
    token_hash character varying(255) NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: blacklisted_tokens_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.blacklisted_tokens_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: blacklisted_tokens_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.blacklisted_tokens_id_seq OWNED BY public.blacklisted_tokens.id;


--
-- Name: blog_views; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.blog_views (
    id bigint NOT NULL,
    blog_id bigint NOT NULL,
    visitor_hash character varying(64) NOT NULL,
    viewed_date date DEFAULT CURRENT_DATE NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: blog_views_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.blog_views_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: blog_views_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.blog_views_id_seq OWNED BY public.blog_views.id;


--
-- Name: blogs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.blogs (
    id bigint NOT NULL,
    title character varying(200) NOT NULL,
    slug character varying(200) NOT NULL,
    content text,
    author_id bigint NOT NULL,
    status character varying(20) DEFAULT 'draft'::character varying NOT NULL,
    tags character varying(500),
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone,
    cover_image_url text,
    cover_image_key text,
    views_total bigint DEFAULT 0 NOT NULL
);


--
-- Name: blogs_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.blogs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: blogs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.blogs_id_seq OWNED BY public.blogs.id;


--
-- Name: coupon_usages; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.coupon_usages (
    id bigint NOT NULL,
    coupon_id bigint NOT NULL,
    order_id bigint NOT NULL,
    user_id bigint NOT NULL,
    discount_applied_cents bigint DEFAULT 0 NOT NULL,
    commission_usd_cents bigint DEFAULT 0 NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: coupon_usages_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.coupon_usages_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: coupon_usages_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.coupon_usages_id_seq OWNED BY public.coupon_usages.id;


--
-- Name: coupons; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.coupons (
    id bigint NOT NULL,
    creator_id bigint NOT NULL,
    code character varying(32) NOT NULL,
    type character varying(20) NOT NULL,
    discount_usd_cents bigint DEFAULT 0 NOT NULL,
    max_discount_cents bigint DEFAULT 10000 NOT NULL,
    commission_pct numeric(5,2) DEFAULT 0 NOT NULL,
    status character varying(20) DEFAULT 'pending'::character varying NOT NULL,
    usage_limit integer,
    used_count integer DEFAULT 0 NOT NULL,
    expires_at timestamp with time zone,
    approved_by bigint,
    approved_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT coupons_status_check CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'approved'::character varying, 'rejected'::character varying])::text[]))),
    CONSTRAINT coupons_type_check CHECK (((type)::text = ANY ((ARRAY['discount'::character varying, 'affiliate'::character varying, 'both'::character varying])::text[])))
);


--
-- Name: coupons_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.coupons_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: coupons_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.coupons_id_seq OWNED BY public.coupons.id;


--
-- Name: memberships; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.memberships (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    plan character varying(20) NOT NULL,
    status character varying(20) DEFAULT 'active'::character varying NOT NULL,
    starts_at timestamp with time zone NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT memberships_plan_check CHECK (((plan)::text = ANY ((ARRAY['free'::character varying, 'basic'::character varying, 'pro'::character varying, 'enterprise'::character varying])::text[]))),
    CONSTRAINT memberships_status_check CHECK (((status)::text = ANY ((ARRAY['active'::character varying, 'cancelled'::character varying, 'expired'::character varying])::text[])))
);


--
-- Name: memberships_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.memberships_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: memberships_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.memberships_id_seq OWNED BY public.memberships.id;


--
-- Name: order_items; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.order_items (
    id bigint NOT NULL,
    order_id bigint NOT NULL,
    product_id bigint NOT NULL,
    type character varying(20) NOT NULL,
    quantity integer NOT NULL,
    unit_price bigint NOT NULL,
    CONSTRAINT order_items_quantity_check CHECK ((quantity > 0)),
    CONSTRAINT order_items_type_check CHECK (((type)::text = ANY ((ARRAY['physical'::character varying, 'digital'::character varying])::text[]))),
    CONSTRAINT order_items_unit_price_check CHECK ((unit_price > 0))
);


--
-- Name: order_items_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.order_items_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: order_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.order_items_id_seq OWNED BY public.order_items.id;


--
-- Name: orders; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.orders (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    status character varying(20) DEFAULT 'pending'::character varying NOT NULL,
    total_amount bigint NOT NULL,
    currency character varying(10) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    service_id bigint,
    amount_usd bigint DEFAULT 0 NOT NULL,
    payment_provider character varying(20),
    coupon_id bigint,
    affiliate_id bigint,
    CONSTRAINT orders_payment_provider_check CHECK (((payment_provider)::text = ANY ((ARRAY['stripe'::character varying, 'razorpay'::character varying])::text[]))),
    CONSTRAINT orders_status_check CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'confirmed'::character varying, 'dispatched'::character varying, 'delivered'::character varying, 'completed'::character varying, 'cancelled'::character varying])::text[]))),
    CONSTRAINT orders_total_amount_check CHECK ((total_amount >= 0))
);


--
-- Name: orders_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.orders_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: orders_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.orders_id_seq OWNED BY public.orders.id;


--
-- Name: otps; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.otps (
    id bigint NOT NULL,
    user_id bigint,
    email character varying(100) NOT NULL,
    code character varying(6) NOT NULL,
    type character varying(20) NOT NULL,
    used boolean DEFAULT false NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT otps_type_check CHECK (((type)::text = ANY ((ARRAY['register'::character varying, 'reset'::character varying])::text[])))
);


--
-- Name: otps_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.otps_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: otps_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.otps_id_seq OWNED BY public.otps.id;


--
-- Name: payments; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.payments (
    id character varying(100) NOT NULL,
    user_id bigint NOT NULL,
    order_id character varying(100) NOT NULL,
    amount bigint NOT NULL,
    currency character varying(10) NOT NULL,
    status character varying(20) DEFAULT 'pending'::character varying NOT NULL,
    provider character varying(50) NOT NULL,
    provider_txn_id character varying(255),
    customer_email character varying(100) NOT NULL,
    description text,
    refunded_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    provider_payment_id text,
    refund_id text,
    CONSTRAINT payments_amount_check CHECK ((amount > 0)),
    CONSTRAINT payments_provider_check CHECK (((provider)::text = ANY ((ARRAY['stripe'::character varying, 'razorpay'::character varying])::text[]))),
    CONSTRAINT payments_status_check CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'completed'::character varying, 'failed'::character varying, 'refunded'::character varying])::text[])))
);


--
-- Name: payouts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.payouts (
    id bigint NOT NULL,
    recipient_id bigint NOT NULL,
    type character varying(30) NOT NULL,
    amount_cents bigint NOT NULL,
    status character varying(20) DEFAULT 'pending'::character varying NOT NULL,
    order_id bigint,
    reference text,
    paid_by bigint,
    paid_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT payouts_amount_cents_check CHECK ((amount_cents > 0)),
    CONSTRAINT payouts_status_check CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'paid'::character varying])::text[]))),
    CONSTRAINT payouts_type_check CHECK (((type)::text = ANY ((ARRAY['partner_sale'::character varying, 'affiliate_commission'::character varying])::text[])))
);


--
-- Name: payouts_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.payouts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: payouts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.payouts_id_seq OWNED BY public.payouts.id;


--
-- Name: products; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.products (
    id bigint NOT NULL,
    name character varying(200) NOT NULL,
    description text,
    type character varying(20) NOT NULL,
    price bigint NOT NULL,
    currency character varying(10) NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT products_price_check CHECK ((price > 0)),
    CONSTRAINT products_type_check CHECK (((type)::text = ANY ((ARRAY['physical'::character varying, 'digital'::character varying])::text[])))
);


--
-- Name: products_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.products_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: products_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.products_id_seq OWNED BY public.products.id;


--
-- Name: refresh_tokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.refresh_tokens (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    token_hash character varying(255) NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    revoked boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: refresh_tokens_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.refresh_tokens_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: refresh_tokens_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.refresh_tokens_id_seq OWNED BY public.refresh_tokens.id;


--
-- Name: saved_blogs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.saved_blogs (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    blog_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: saved_blogs_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.saved_blogs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: saved_blogs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.saved_blogs_id_seq OWNED BY public.saved_blogs.id;


--
-- Name: service_views; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.service_views (
    id bigint NOT NULL,
    service_id bigint NOT NULL,
    visitor_hash character varying(64) NOT NULL,
    viewed_date date DEFAULT CURRENT_DATE NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: service_views_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.service_views_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: service_views_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.service_views_id_seq OWNED BY public.service_views.id;


--
-- Name: services; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.services (
    id bigint NOT NULL,
    owner_id bigint NOT NULL,
    title character varying(200) NOT NULL,
    description text,
    price_usd_cents bigint NOT NULL,
    status character varying(20) DEFAULT 'pending'::character varying NOT NULL,
    cover_image_url text,
    cover_image_key text,
    views_total bigint DEFAULT 0 NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone,
    CONSTRAINT services_price_usd_cents_check CHECK ((price_usd_cents > 0)),
    CONSTRAINT services_status_check CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'approved'::character varying, 'rejected'::character varying, 'live'::character varying, 'paused'::character varying])::text[])))
);


--
-- Name: services_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.services_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: services_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.services_id_seq OWNED BY public.services.id;


--
-- Name: sessions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sessions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id bigint NOT NULL,
    token_hash character varying(255) NOT NULL,
    user_agent character varying(512) DEFAULT ''::character varying NOT NULL,
    ip_address character varying(45) DEFAULT ''::character varying NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    revoked boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_active timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: ticket_messages; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ticket_messages (
    id bigint NOT NULL,
    ticket_id bigint NOT NULL,
    sender_id bigint NOT NULL,
    message text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: ticket_messages_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.ticket_messages_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: ticket_messages_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.ticket_messages_id_seq OWNED BY public.ticket_messages.id;


--
-- Name: tickets; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tickets (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    subject character varying(200) NOT NULL,
    description text NOT NULL,
    status character varying(20) DEFAULT 'open'::character varying NOT NULL,
    priority character varying(20) DEFAULT 'medium'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    assigned_to bigint,
    closed_at timestamp with time zone,
    CONSTRAINT tickets_priority_check CHECK (((priority)::text = ANY ((ARRAY['low'::character varying, 'medium'::character varying, 'high'::character varying, 'urgent'::character varying])::text[]))),
    CONSTRAINT tickets_status_check CHECK (((status)::text = ANY ((ARRAY['open'::character varying, 'in_progress'::character varying, 'resolved'::character varying, 'closed'::character varying])::text[])))
);


--
-- Name: tickets_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tickets_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tickets_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tickets_id_seq OWNED BY public.tickets.id;


--
-- Name: upgrade_requests; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.upgrade_requests (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    from_role character varying(20) NOT NULL,
    to_role character varying(20) DEFAULT 'partner'::character varying NOT NULL,
    status character varying(20) DEFAULT 'pending'::character varying NOT NULL,
    reviewed_by bigint,
    reviewed_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT upgrade_requests_status_check CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'approved'::character varying, 'rejected'::character varying])::text[])))
);


--
-- Name: upgrade_requests_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.upgrade_requests_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: upgrade_requests_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.upgrade_requests_id_seq OWNED BY public.upgrade_requests.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    name character varying(100) NOT NULL,
    email character varying(100) NOT NULL,
    password_hash character varying(255),
    role character varying(20) DEFAULT 'user'::character varying NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone,
    is_verified boolean DEFAULT false NOT NULL,
    upgrade_requested_at timestamp with time zone,
    CONSTRAINT users_role_check CHECK (((role)::text = ANY ((ARRAY['user'::character varying, 'writer'::character varying, 'partner'::character varying, 'worker'::character varying, 'admin'::character varying, 'superadmin'::character varying])::text[])))
);


--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: webhook_events; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.webhook_events (
    id character varying(100) NOT NULL,
    provider character varying(50) NOT NULL,
    event_type character varying(100) NOT NULL,
    payload bytea,
    signature character varying(255),
    processed boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: blacklisted_tokens id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blacklisted_tokens ALTER COLUMN id SET DEFAULT nextval('public.blacklisted_tokens_id_seq'::regclass);


--
-- Name: blog_views id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blog_views ALTER COLUMN id SET DEFAULT nextval('public.blog_views_id_seq'::regclass);


--
-- Name: blogs id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blogs ALTER COLUMN id SET DEFAULT nextval('public.blogs_id_seq'::regclass);


--
-- Name: coupon_usages id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.coupon_usages ALTER COLUMN id SET DEFAULT nextval('public.coupon_usages_id_seq'::regclass);


--
-- Name: coupons id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.coupons ALTER COLUMN id SET DEFAULT nextval('public.coupons_id_seq'::regclass);


--
-- Name: memberships id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.memberships ALTER COLUMN id SET DEFAULT nextval('public.memberships_id_seq'::regclass);


--
-- Name: order_items id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.order_items ALTER COLUMN id SET DEFAULT nextval('public.order_items_id_seq'::regclass);


--
-- Name: orders id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.orders ALTER COLUMN id SET DEFAULT nextval('public.orders_id_seq'::regclass);


--
-- Name: otps id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.otps ALTER COLUMN id SET DEFAULT nextval('public.otps_id_seq'::regclass);


--
-- Name: payouts id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payouts ALTER COLUMN id SET DEFAULT nextval('public.payouts_id_seq'::regclass);


--
-- Name: products id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.products ALTER COLUMN id SET DEFAULT nextval('public.products_id_seq'::regclass);


--
-- Name: refresh_tokens id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refresh_tokens ALTER COLUMN id SET DEFAULT nextval('public.refresh_tokens_id_seq'::regclass);


--
-- Name: saved_blogs id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.saved_blogs ALTER COLUMN id SET DEFAULT nextval('public.saved_blogs_id_seq'::regclass);


--
-- Name: service_views id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.service_views ALTER COLUMN id SET DEFAULT nextval('public.service_views_id_seq'::regclass);


--
-- Name: services id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.services ALTER COLUMN id SET DEFAULT nextval('public.services_id_seq'::regclass);


--
-- Name: ticket_messages id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ticket_messages ALTER COLUMN id SET DEFAULT nextval('public.ticket_messages_id_seq'::regclass);


--
-- Name: tickets id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tickets ALTER COLUMN id SET DEFAULT nextval('public.tickets_id_seq'::regclass);


--
-- Name: upgrade_requests id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.upgrade_requests ALTER COLUMN id SET DEFAULT nextval('public.upgrade_requests_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: blacklisted_tokens blacklisted_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blacklisted_tokens
    ADD CONSTRAINT blacklisted_tokens_pkey PRIMARY KEY (id);


--
-- Name: blog_views blog_views_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blog_views
    ADD CONSTRAINT blog_views_pkey PRIMARY KEY (id);


--
-- Name: blogs blogs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blogs
    ADD CONSTRAINT blogs_pkey PRIMARY KEY (id);


--
-- Name: coupon_usages coupon_usages_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.coupon_usages
    ADD CONSTRAINT coupon_usages_pkey PRIMARY KEY (id);


--
-- Name: coupons coupons_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.coupons
    ADD CONSTRAINT coupons_pkey PRIMARY KEY (id);


--
-- Name: memberships memberships_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.memberships
    ADD CONSTRAINT memberships_pkey PRIMARY KEY (id);


--
-- Name: order_items order_items_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_pkey PRIMARY KEY (id);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- Name: otps otps_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.otps
    ADD CONSTRAINT otps_pkey PRIMARY KEY (id);


--
-- Name: payments payments_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_pkey PRIMARY KEY (id);


--
-- Name: payouts payouts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payouts
    ADD CONSTRAINT payouts_pkey PRIMARY KEY (id);


--
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);


--
-- Name: refresh_tokens refresh_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_pkey PRIMARY KEY (id);


--
-- Name: saved_blogs saved_blogs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.saved_blogs
    ADD CONSTRAINT saved_blogs_pkey PRIMARY KEY (id);


--
-- Name: saved_blogs saved_blogs_user_id_blog_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.saved_blogs
    ADD CONSTRAINT saved_blogs_user_id_blog_id_key UNIQUE (user_id, blog_id);


--
-- Name: service_views service_views_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.service_views
    ADD CONSTRAINT service_views_pkey PRIMARY KEY (id);


--
-- Name: services services_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.services
    ADD CONSTRAINT services_pkey PRIMARY KEY (id);


--
-- Name: sessions sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_pkey PRIMARY KEY (id);


--
-- Name: ticket_messages ticket_messages_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ticket_messages
    ADD CONSTRAINT ticket_messages_pkey PRIMARY KEY (id);


--
-- Name: tickets tickets_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tickets
    ADD CONSTRAINT tickets_pkey PRIMARY KEY (id);


--
-- Name: upgrade_requests upgrade_requests_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.upgrade_requests
    ADD CONSTRAINT upgrade_requests_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: webhook_events webhook_events_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.webhook_events
    ADD CONSTRAINT webhook_events_pkey PRIMARY KEY (id);


--
-- Name: idx_blacklisted_tokens_hash; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_blacklisted_tokens_hash ON public.blacklisted_tokens USING btree (token_hash);


--
-- Name: idx_blog_views_blog_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_blog_views_blog_id ON public.blog_views USING btree (blog_id);


--
-- Name: idx_blog_views_date; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_blog_views_date ON public.blog_views USING btree (viewed_date);


--
-- Name: idx_blog_views_unique; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_blog_views_unique ON public.blog_views USING btree (blog_id, visitor_hash, viewed_date);


--
-- Name: idx_blogs_author_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_blogs_author_id ON public.blogs USING btree (author_id);


--
-- Name: idx_blogs_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_blogs_deleted_at ON public.blogs USING btree (deleted_at);


--
-- Name: idx_blogs_slug; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_blogs_slug ON public.blogs USING btree (slug);


--
-- Name: idx_blogs_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_blogs_status ON public.blogs USING btree (status);


--
-- Name: idx_coupon_usages_coupon_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_coupon_usages_coupon_id ON public.coupon_usages USING btree (coupon_id);


--
-- Name: idx_coupon_usages_order_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_coupon_usages_order_id ON public.coupon_usages USING btree (order_id);


--
-- Name: idx_coupon_usages_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_coupon_usages_user_id ON public.coupon_usages USING btree (user_id);


--
-- Name: idx_coupons_code; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_coupons_code ON public.coupons USING btree (code);


--
-- Name: idx_coupons_creator; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_coupons_creator ON public.coupons USING btree (creator_id);


--
-- Name: idx_coupons_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_coupons_status ON public.coupons USING btree (status);


--
-- Name: idx_memberships_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_memberships_status ON public.memberships USING btree (status);


--
-- Name: idx_memberships_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_memberships_user_id ON public.memberships USING btree (user_id);


--
-- Name: idx_order_items_order_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_order_items_order_id ON public.order_items USING btree (order_id);


--
-- Name: idx_order_items_product_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_order_items_product_id ON public.order_items USING btree (product_id);


--
-- Name: idx_orders_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_orders_status ON public.orders USING btree (status);


--
-- Name: idx_orders_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_orders_user_id ON public.orders USING btree (user_id);


--
-- Name: idx_otps_email; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_otps_email ON public.otps USING btree (email);


--
-- Name: idx_otps_expires_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_otps_expires_at ON public.otps USING btree (expires_at);


--
-- Name: idx_otps_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_otps_type ON public.otps USING btree (type);


--
-- Name: idx_payments_order_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_payments_order_id ON public.payments USING btree (order_id);


--
-- Name: idx_payments_provider; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_payments_provider ON public.payments USING btree (provider);


--
-- Name: idx_payments_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_payments_status ON public.payments USING btree (status);


--
-- Name: idx_payments_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_payments_user_id ON public.payments USING btree (user_id);


--
-- Name: idx_payouts_recipient; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_payouts_recipient ON public.payouts USING btree (recipient_id);


--
-- Name: idx_payouts_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_payouts_status ON public.payouts USING btree (status);


--
-- Name: idx_payouts_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_payouts_type ON public.payouts USING btree (type);


--
-- Name: idx_products_is_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_products_is_active ON public.products USING btree (is_active);


--
-- Name: idx_products_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_products_type ON public.products USING btree (type);


--
-- Name: idx_refresh_tokens_token_hash; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_refresh_tokens_token_hash ON public.refresh_tokens USING btree (token_hash);


--
-- Name: idx_refresh_tokens_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_refresh_tokens_user_id ON public.refresh_tokens USING btree (user_id);


--
-- Name: idx_saved_blogs_blog_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_saved_blogs_blog_id ON public.saved_blogs USING btree (blog_id);


--
-- Name: idx_saved_blogs_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_saved_blogs_user_id ON public.saved_blogs USING btree (user_id);


--
-- Name: idx_service_views_date; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_service_views_date ON public.service_views USING btree (viewed_date);


--
-- Name: idx_service_views_service_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_service_views_service_id ON public.service_views USING btree (service_id);


--
-- Name: idx_service_views_unique; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_service_views_unique ON public.service_views USING btree (service_id, visitor_hash, viewed_date);


--
-- Name: idx_services_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_services_deleted_at ON public.services USING btree (deleted_at);


--
-- Name: idx_services_owner_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_services_owner_id ON public.services USING btree (owner_id);


--
-- Name: idx_services_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_services_status ON public.services USING btree (status);


--
-- Name: idx_sessions_expires_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_sessions_expires_at ON public.sessions USING btree (expires_at);


--
-- Name: idx_sessions_revoked; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_sessions_revoked ON public.sessions USING btree (revoked);


--
-- Name: idx_sessions_token_hash; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_sessions_token_hash ON public.sessions USING btree (token_hash);


--
-- Name: idx_sessions_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_sessions_user_id ON public.sessions USING btree (user_id);


--
-- Name: idx_ticket_messages_sender_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_ticket_messages_sender_id ON public.ticket_messages USING btree (sender_id);


--
-- Name: idx_ticket_messages_ticket_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_ticket_messages_ticket_id ON public.ticket_messages USING btree (ticket_id);


--
-- Name: idx_tickets_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_tickets_status ON public.tickets USING btree (status);


--
-- Name: idx_tickets_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_tickets_user_id ON public.tickets USING btree (user_id);


--
-- Name: idx_upgrade_requests_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_upgrade_requests_status ON public.upgrade_requests USING btree (status);


--
-- Name: idx_upgrade_requests_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_upgrade_requests_user_id ON public.upgrade_requests USING btree (user_id);


--
-- Name: idx_users_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_users_email ON public.users USING btree (email);


--
-- Name: idx_webhook_events_processed; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_webhook_events_processed ON public.webhook_events USING btree (processed);


--
-- Name: idx_webhook_events_provider; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_webhook_events_provider ON public.webhook_events USING btree (provider);


--
-- Name: blog_views blog_views_blog_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blog_views
    ADD CONSTRAINT blog_views_blog_id_fkey FOREIGN KEY (blog_id) REFERENCES public.blogs(id) ON DELETE CASCADE;


--
-- Name: blogs blogs_author_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.blogs
    ADD CONSTRAINT blogs_author_id_fkey FOREIGN KEY (author_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: coupon_usages coupon_usages_coupon_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.coupon_usages
    ADD CONSTRAINT coupon_usages_coupon_id_fkey FOREIGN KEY (coupon_id) REFERENCES public.coupons(id) ON DELETE CASCADE;


--
-- Name: coupon_usages coupon_usages_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.coupon_usages
    ADD CONSTRAINT coupon_usages_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE CASCADE;


--
-- Name: coupon_usages coupon_usages_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.coupon_usages
    ADD CONSTRAINT coupon_usages_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: coupons coupons_approved_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.coupons
    ADD CONSTRAINT coupons_approved_by_fkey FOREIGN KEY (approved_by) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: coupons coupons_creator_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.coupons
    ADD CONSTRAINT coupons_creator_id_fkey FOREIGN KEY (creator_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: memberships memberships_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.memberships
    ADD CONSTRAINT memberships_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: order_items order_items_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE CASCADE;


--
-- Name: order_items order_items_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id) ON DELETE RESTRICT;


--
-- Name: orders orders_affiliate_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_affiliate_id_fkey FOREIGN KEY (affiliate_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: orders orders_coupon_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_coupon_id_fkey FOREIGN KEY (coupon_id) REFERENCES public.coupons(id) ON DELETE SET NULL;


--
-- Name: orders orders_service_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_service_id_fkey FOREIGN KEY (service_id) REFERENCES public.services(id) ON DELETE SET NULL;


--
-- Name: orders orders_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE RESTRICT;


--
-- Name: otps otps_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.otps
    ADD CONSTRAINT otps_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: payments payments_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE RESTRICT;


--
-- Name: payouts payouts_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payouts
    ADD CONSTRAINT payouts_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE SET NULL;


--
-- Name: payouts payouts_paid_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payouts
    ADD CONSTRAINT payouts_paid_by_fkey FOREIGN KEY (paid_by) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: payouts payouts_recipient_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payouts
    ADD CONSTRAINT payouts_recipient_id_fkey FOREIGN KEY (recipient_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: refresh_tokens refresh_tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: saved_blogs saved_blogs_blog_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.saved_blogs
    ADD CONSTRAINT saved_blogs_blog_id_fkey FOREIGN KEY (blog_id) REFERENCES public.blogs(id) ON DELETE CASCADE;


--
-- Name: saved_blogs saved_blogs_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.saved_blogs
    ADD CONSTRAINT saved_blogs_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: service_views service_views_service_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.service_views
    ADD CONSTRAINT service_views_service_id_fkey FOREIGN KEY (service_id) REFERENCES public.services(id) ON DELETE CASCADE;


--
-- Name: services services_owner_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.services
    ADD CONSTRAINT services_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: sessions sessions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sessions
    ADD CONSTRAINT sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: ticket_messages ticket_messages_sender_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ticket_messages
    ADD CONSTRAINT ticket_messages_sender_id_fkey FOREIGN KEY (sender_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: ticket_messages ticket_messages_ticket_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ticket_messages
    ADD CONSTRAINT ticket_messages_ticket_id_fkey FOREIGN KEY (ticket_id) REFERENCES public.tickets(id) ON DELETE CASCADE;


--
-- Name: tickets tickets_assigned_to_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tickets
    ADD CONSTRAINT tickets_assigned_to_fkey FOREIGN KEY (assigned_to) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: tickets tickets_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tickets
    ADD CONSTRAINT tickets_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: upgrade_requests upgrade_requests_reviewed_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.upgrade_requests
    ADD CONSTRAINT upgrade_requests_reviewed_by_fkey FOREIGN KEY (reviewed_by) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: upgrade_requests upgrade_requests_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.upgrade_requests
    ADD CONSTRAINT upgrade_requests_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict 8A7sykcbWn6skfC0SwHPsCQpGj3jW3cn8bfOLLXQLC0nMQMU1Sh6XMHbyqUJuFo

