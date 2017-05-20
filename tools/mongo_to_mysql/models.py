# -*- coding:utf-8 -*-
from peewee import *


db = MySQLDatabase(
    'manganow',
    host='localhost',
    port=3306,
    user='root',
    passwd='')


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
    tree_type = FixedCharField(unique=True, max_length=191)
    sub_asins = FixedCharField(max_length=255)
    kindle = BooleanField()
    date_publish = DateField()
    image_s_url = FixedCharField(max_length=255)
    image_s_width = IntegerField()
    image_s_height = IntegerField()
    image_m_url = FixedCharField(max_length=255)
    image_m_width = IntegerField()
    image_m_height = IntegerField()
    image_l_url = FixedCharField(max_length=255)
    image_l_width = IntegerField()
    image_l_height = IntegerField()
    title = TextField()
    region = FixedCharField(max_length=255)
    publisher = ForeignKeyField(Publishers, related_name='books')
    author = ForeignKeyField(Authors, related_name='books')
    updated_at = DateTimeField()
    created_at = DateTimeField()

    class Meta:
        database = db
