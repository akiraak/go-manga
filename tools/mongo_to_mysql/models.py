# -*- coding:utf-8 -*-
from peewee import *


db = MySQLDatabase(
    'manganow',
    host='localhost',
    port=3306,
    user='root',
    passwd='')


class Titles(Model):
    id = BigIntegerField(primary_key=True)
    created_at = DateTimeField()

    class Meta:
        database = db


class Publishers(Model):
    id = BigIntegerField(primary_key=True)
    name = FixedCharField(unique=True, max_length=191)
    created_at = DateTimeField()

    class Meta:
        database = db


class Authors(Model):
    id = BigIntegerField()
    name = FixedCharField(max_length=191)
    created_at = DateTimeField()

    class Meta:
        database = db


class Books(Model):
    id = BigIntegerField()
    asin = FixedCharField(max_length=255)
    date_publish = FixedCharField(max_length=8)
    publish_type = FixedCharField(max_length=255)
    image_s_url = FixedCharField(max_length=255)
    image_s_width = IntegerField()
    image_s_height = IntegerField()
    image_m_url = FixedCharField(max_length=255)
    image_m_width = IntegerField()
    image_m_height = IntegerField()
    image_l_url = FixedCharField(max_length=255)
    image_l_width = IntegerField()
    image_l_height = IntegerField()
    name = TextField()
    region = FixedCharField(max_length=255)
    book_title = ForeignKeyField(Titles, db_column='title_id', related_name='books')
    publisher = ForeignKeyField(Publishers, related_name='books')
    author = ForeignKeyField(Authors, related_name='books')
    updated_at = DateTimeField()
    created_at = DateTimeField()

    class Meta:
        database = db
