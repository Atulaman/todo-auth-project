PGDMP                        |            todo #   16.4 (Ubuntu 16.4-0ubuntu0.24.04.2) #   16.4 (Ubuntu 16.4-0ubuntu0.24.04.2)     s           0    0    ENCODING    ENCODING        SET client_encoding = 'UTF8';
                      false            t           0    0 
   STDSTRINGS 
   STDSTRINGS     (   SET standard_conforming_strings = 'on';
                      false            u           0    0 
   SEARCHPATH 
   SEARCHPATH     8   SELECT pg_catalog.set_config('search_path', '', false);
                      false            v           1262    16389    todo    DATABASE     p   CREATE DATABASE todo WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.UTF-8';
    DROP DATABASE todo;
                postgres    false            �            1259    16398    Tasks    TABLE     X   CREATE TABLE public."Tasks" (
    id integer,
    description character varying(500)
);
    DROP TABLE public."Tasks";
       public         heap    postgres    false            �            1259    16547    auth    TABLE     x   CREATE TABLE public.auth (
    username character varying(100) NOT NULL,
    password character varying(20) NOT NULL
);
    DROP TABLE public.auth;
       public         heap    postgres    false            �            1259    16584    session    TABLE     �   CREATE TABLE public.session (
    username character varying(100) NOT NULL,
    session_id character varying(200) NOT NULL,
    created_at timestamp without time zone NOT NULL
);
    DROP TABLE public.session;
       public         heap    postgres    false            n          0    16398    Tasks 
   TABLE DATA           2   COPY public."Tasks" (id, description) FROM stdin;
    public          postgres    false    215   �       o          0    16547    auth 
   TABLE DATA           2   COPY public.auth (username, password) FROM stdin;
    public          postgres    false    216   (       p          0    16584    session 
   TABLE DATA           C   COPY public.session (username, session_id, created_at) FROM stdin;
    public          postgres    false    217   c       �           2606    16551    auth auth_username_key 
   CONSTRAINT     U   ALTER TABLE ONLY public.auth
    ADD CONSTRAINT auth_username_key UNIQUE (username);
 @   ALTER TABLE ONLY public.auth DROP CONSTRAINT auth_username_key;
       public            postgres    false    216            �           2606    16588    session session_pkey 
   CONSTRAINT     Z   ALTER TABLE ONLY public.session
    ADD CONSTRAINT session_pkey PRIMARY KEY (session_id);
 >   ALTER TABLE ONLY public.session DROP CONSTRAINT session_pkey;
       public            postgres    false    217            �           2606    16589    session session_username_fkey    FK CONSTRAINT     �   ALTER TABLE ONLY public.session
    ADD CONSTRAINT session_username_fkey FOREIGN KEY (username) REFERENCES public.auth(username);
 G   ALTER TABLE ONLY public.session DROP CONSTRAINT session_username_fkey;
       public          postgres    false    217    3291    216            n   2   x�3�,I,�V0�2�0��L .#È��0�2�0̸�!s�=... \��      o   +   x�K,)�I�M������.u042�L,��M�(ID����� W6�      p   �   x�]�AN1��?��N�QJ;w� &n(�q�n��[&j`�������v�hAT��V�&E����zZ��	�}6k@X#J7�=�=��b�gg4��?�L�g���z�:��騊���5I:��R �Y��!�O����������o���]�:�-C�L�Kv�<e�,���Ǥy�>��+���
�/e�Ntmņ��1Ck�!��5yd�o,�X�[����+��?
_�Փ i8���8� I�i�     