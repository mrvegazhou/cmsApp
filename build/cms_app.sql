/*
 Navicat Premium Data Transfer

 Source Server         : postgres-localhost
 Source Server Type    : PostgreSQL
 Source Server Version : 130002 (130002)
 Source Host           : localhost:5432
 Source Catalog        : app_cms_db
 Source Schema         : cms_app

 Target Server Type    : PostgreSQL
 Target Server Version : 130002 (130002)
 File Encoding         : 65001

 Date: 01/01/2024 01:05:53
*/


-- ----------------------------
-- Sequence structure for app_article_collect_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "cms_app"."app_article_collect_id_seq";
CREATE SEQUENCE "cms_app"."app_article_collect_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;
ALTER SEQUENCE "cms_app"."app_article_collect_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for app_article_favorites_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "cms_app"."app_article_favorites_id_seq";
CREATE SEQUENCE "cms_app"."app_article_favorites_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;
ALTER SEQUENCE "cms_app"."app_article_favorites_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for app_article_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "cms_app"."app_article_id_seq";
CREATE SEQUENCE "cms_app"."app_article_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;
ALTER SEQUENCE "cms_app"."app_article_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for app_article_like_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "cms_app"."app_article_like_id_seq";
CREATE SEQUENCE "cms_app"."app_article_like_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;
ALTER SEQUENCE "cms_app"."app_article_like_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for app_favorites_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "cms_app"."app_favorites_id_seq";
CREATE SEQUENCE "cms_app"."app_favorites_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;
ALTER SEQUENCE "cms_app"."app_favorites_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for app_file_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "cms_app"."app_file_id_seq";
CREATE SEQUENCE "cms_app"."app_file_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;
ALTER SEQUENCE "cms_app"."app_file_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for app_imgs_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "cms_app"."app_imgs_id_seq";
CREATE SEQUENCE "cms_app"."app_imgs_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;
ALTER SEQUENCE "cms_app"."app_imgs_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for app_user_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "cms_app"."app_user_id_seq";
CREATE SEQUENCE "cms_app"."app_user_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;
ALTER SEQUENCE "cms_app"."app_user_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for site_config_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "cms_app"."site_config_id_seq";
CREATE SEQUENCE "cms_app"."site_config_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;
ALTER SEQUENCE "cms_app"."site_config_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for site_info_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "cms_app"."site_info_id_seq";
CREATE SEQUENCE "cms_app"."site_info_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;
ALTER SEQUENCE "cms_app"."site_info_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Table structure for app_article
-- ----------------------------
DROP TABLE IF EXISTS "cms_app"."app_article";
CREATE TABLE "cms_app"."app_article" (
  "id" int8 NOT NULL DEFAULT nextval('"cms_app".app_article_id_seq'::regclass),
  "title" varchar(255) COLLATE "pg_catalog"."default",
  "description" varchar(255) COLLATE "pg_catalog"."default",
  "author_id" int8 NOT NULL,
  "content" text COLLATE "pg_catalog"."default",
  "view_count" int8 DEFAULT 0,
  "comment_count" int8 DEFAULT 0,
  "collection_count" int8 DEFAULT 0,
  "like_count" int8 DEFAULT 0,
  "create_time" timestamptz(6) NOT NULL,
  "update_time" timestamptz(6),
  "delete_time" timestamptz(6),
  "cover_url" varchar(255) COLLATE "pg_catalog"."default",
  "share_count" int8,
  "state" int2
)
;
ALTER TABLE "cms_app"."app_article" OWNER TO "postgres";
COMMENT ON COLUMN "cms_app"."app_article"."cover_url" IS '封面';

-- ----------------------------
-- Table structure for app_article_favorites
-- ----------------------------
DROP TABLE IF EXISTS "cms_app"."app_article_favorites";
CREATE TABLE "cms_app"."app_article_favorites" (
  "id" int8 NOT NULL DEFAULT nextval('"cms_app".app_article_favorites_id_seq'::regclass),
  "article_id" int8 NOT NULL,
  "favorites_id" int8 NOT NULL,
  "create_time" timestamptz(6) NOT NULL
)
;
ALTER TABLE "cms_app"."app_article_favorites" OWNER TO "postgres";

-- ----------------------------
-- Table structure for app_article_like
-- ----------------------------
DROP TABLE IF EXISTS "cms_app"."app_article_like";
CREATE TABLE "cms_app"."app_article_like" (
  "id" int8 NOT NULL DEFAULT nextval('"cms_app".app_article_like_id_seq'::regclass),
  "article_id" int8 NOT NULL,
  "user_id" int8 NOT NULL,
  "create_time" timestamptz(6) NOT NULL
)
;
ALTER TABLE "cms_app"."app_article_like" OWNER TO "postgres";

-- ----------------------------
-- Table structure for app_favorites
-- ----------------------------
DROP TABLE IF EXISTS "cms_app"."app_favorites";
CREATE TABLE "cms_app"."app_favorites" (
  "id" int8 NOT NULL DEFAULT nextval('"cms_app".app_favorites_id_seq'::regclass),
  "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "create_time" timestamptz(6) NOT NULL,
  "user_id" int8 NOT NULL
)
;
ALTER TABLE "cms_app"."app_favorites" OWNER TO "postgres";

-- ----------------------------
-- Table structure for app_file
-- ----------------------------
DROP TABLE IF EXISTS "cms_app"."app_file";
CREATE TABLE "cms_app"."app_file" (
  "id" int8 NOT NULL DEFAULT nextval('"cms_app".app_file_id_seq'::regclass),
  "url" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "base_dir" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "type" int2 NOT NULL,
  "media_type" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "size" float8 NOT NULL,
  "create_time" timestamptz(6) NOT NULL,
  "update_time" timestamptz(6)
)
;
ALTER TABLE "cms_app"."app_file" OWNER TO "postgres";

-- ----------------------------
-- Table structure for app_imgs
-- ----------------------------
DROP TABLE IF EXISTS "cms_app"."app_imgs";
CREATE TABLE "cms_app"."app_imgs" (
  "id" int8 NOT NULL DEFAULT nextval('"cms_app".app_imgs_id_seq'::regclass),
  "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "path" varchar(255) COLLATE "pg_catalog"."default",
  "type" int2 NOT NULL,
  "create_time" timestamptz(6) NOT NULL,
  "update_time" timestamptz(6),
  "tags" varchar(255) COLLATE "pg_catalog"."default",
  "resource_id" int8 DEFAULT 0,
  "width" int4,
  "height" int4,
  "user_id" int8
)
;
ALTER TABLE "cms_app"."app_imgs" OWNER TO "postgres";

-- ----------------------------
-- Table structure for app_imgs_temp
-- ----------------------------
DROP TABLE IF EXISTS "cms_app"."app_imgs_temp";
CREATE TABLE "cms_app"."app_imgs_temp" (
  "id" int8 NOT NULL,
  "url" varchar(255) COLLATE "pg_catalog"."default",
  "base_dir" varchar(255) COLLATE "pg_catalog"."default",
  "create_time" timestamptz(6),
  "type" int2,
  "resource_id" int8
)
;
ALTER TABLE "cms_app"."app_imgs_temp" OWNER TO "postgres";

-- ----------------------------
-- Table structure for app_user
-- ----------------------------
DROP TABLE IF EXISTS "cms_app"."app_user";
CREATE TABLE "cms_app"."app_user" (
  "id" int8 NOT NULL DEFAULT nextval('"cms_app".app_user_id_seq'::regclass),
  "nickname" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "email" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "password" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "salt" varchar COLLATE "pg_catalog"."default" NOT NULL,
  "update_time" timestamptz(6),
  "delete_time" timestamptz(6),
  "refresh_token" varchar COLLATE "pg_catalog"."default",
  "expir_time" timestamptz(6),
  "avatar_url" varchar COLLATE "pg_catalog"."default",
  "about" varchar COLLATE "pg_catalog"."default",
  "create_time" timestamptz(6),
  "avatar" varchar(255) COLLATE "pg_catalog"."default",
  "phone" varchar(255) COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "cms_app"."app_user" OWNER TO "postgres";
COMMENT ON COLUMN "cms_app"."app_user"."avatar_url" IS '头像';
COMMENT ON COLUMN "cms_app"."app_user"."about" IS '简介';

-- ----------------------------
-- Table structure for site_config
-- ----------------------------
DROP TABLE IF EXISTS "cms_app"."site_config";
CREATE TABLE "cms_app"."site_config" (
  "id" int4 NOT NULL DEFAULT nextval('"cms_app".site_config_id_seq'::regclass),
  "delete_time" timestamptz(6),
  "update_time" timestamptz(6),
  "create_time" timestamptz(6) NOT NULL,
  "version" int2 NOT NULL,
  "email" varchar COLLATE "pg_catalog"."default",
  "jwt" varchar COLLATE "pg_catalog"."default",
  "phone" varchar COLLATE "pg_catalog"."default",
  "qq" varchar COLLATE "pg_catalog"."default",
  "wechat" varchar COLLATE "pg_catalog"."default",
  "weibo" varchar COLLATE "pg_catalog"."default",
  "type" int2
)
;
ALTER TABLE "cms_app"."site_config" OWNER TO "postgres";

-- ----------------------------
-- Table structure for site_info
-- ----------------------------
DROP TABLE IF EXISTS "cms_app"."site_info";
CREATE TABLE "cms_app"."site_info" (
  "id" int4 NOT NULL DEFAULT nextval('"cms_app".site_info_id_seq'::regclass),
  "content" text COLLATE "pg_catalog"."default" NOT NULL,
  "type" int2 NOT NULL,
  "create_time" timestamptz(6) NOT NULL,
  "update_time" timestamptz(6),
  "title" varchar(255) COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "cms_app"."site_info" OWNER TO "postgres";

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"cms_app"."app_article_collect_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"cms_app"."app_article_favorites_id_seq"', 1, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"cms_app"."app_article_id_seq"', 152, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"cms_app"."app_article_like_id_seq"', 19, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"cms_app"."app_favorites_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"cms_app"."app_file_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"cms_app"."app_imgs_id_seq"', 140, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"cms_app"."app_user_id_seq"', 11, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"cms_app"."site_config_id_seq"', 1, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"cms_app"."site_info_id_seq"', 1, true);

-- ----------------------------
-- Primary Key structure for table app_article
-- ----------------------------
ALTER TABLE "cms_app"."app_article" ADD CONSTRAINT "app_article_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table app_article_favorites
-- ----------------------------
ALTER TABLE "cms_app"."app_article_favorites" ADD CONSTRAINT "app_article_favorites_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table app_article_like
-- ----------------------------
ALTER TABLE "cms_app"."app_article_like" ADD CONSTRAINT "app_article_like_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table app_favorites
-- ----------------------------
ALTER TABLE "cms_app"."app_favorites" ADD CONSTRAINT "app_favorites_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table app_file
-- ----------------------------
ALTER TABLE "cms_app"."app_file" ADD CONSTRAINT "app_file_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table app_imgs
-- ----------------------------
ALTER TABLE "cms_app"."app_imgs" ADD CONSTRAINT "app_imgs_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table app_imgs_temp
-- ----------------------------
ALTER TABLE "cms_app"."app_imgs_temp" ADD CONSTRAINT "app_imgs_temp_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table app_user
-- ----------------------------
ALTER TABLE "cms_app"."app_user" ADD CONSTRAINT "app_user_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table site_config
-- ----------------------------
ALTER TABLE "cms_app"."site_config" ADD CONSTRAINT "sys_config_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table site_info
-- ----------------------------
ALTER TABLE "cms_app"."site_info" ADD CONSTRAINT "app_introduce_pkey" PRIMARY KEY ("id");

CREATE INDEX idx_name ON "cms_app"."app_imgs" (name);
CREATE INDEX idx_email ON "cms_app"."app_user" (email);
CREATE INDEX idx_name_temp ON "cms_app"."app_imgs_temp" (name);
CREATE INDEX idx_article_id ON "cms_app"."app_article_history" (article_id);

