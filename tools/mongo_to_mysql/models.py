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
    date_created = DateTimeField()

    class Meta:
        database = db


class Authors(Model):
    id = BigIntegerField()
    name = FixedCharField(max_length=191)
    date_created = DateTimeField()

    class Meta:
        database = db


class Books(Model):
    id = BigIntegerField()
    kindle = BooleanField()
    date_publish = DateTimeField()
    image_s_url = FixedCharField(max_length=255)
    image_s_width = IntegerField()
    image_s_height = IntegerField()
    image_m_url = FixedCharField(max_length=255)
    image_m_width = IntegerField()
    image_m_height = IntegerField()
    image_l_url = FixedCharField(max_length=255)
    image_l_width = IntegerField()
    image_l_height = IntegerField()
    asin = FixedCharField(max_length=255)
    title = TextField()
    region = FixedCharField(max_length=255)
    publisher = ForeignKeyField(Publishers, related_name='books')
    author = ForeignKeyField(Authors, related_name='books')
    date_last_modify = DateTimeField()
    date_created = DateTimeField()

    class Meta:
        database = db
