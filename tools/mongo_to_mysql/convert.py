# -*- coding:utf-8 -*-
import sys
from pymongo import MongoClient
from peewee import *
from models import db, Publishers, Authors, Books


# Convert publishers and authors
def conv_info_models(model, mongo_objects):
    print('Start converting {}'.format(str(model)))

    # Create mysql objects
    with db.transaction():
        for mongo_object in mongo_objects:
            model.create(
                name=mongo_object['name'],
                date_created=mongo_object['date_created'])

    # Make a map of mongo_id to mysql_id
    mongo_id_to_mysql_id = {}
    with db.transaction():
        for mysql_object in model.select():
            mongo_id = None
            for mongo_object in mongo_objects:
                if mongo_object['name'] == mysql_object.name:
                    mongo_id = str(mongo_object['_id'])
                    break
            mongo_id_to_mysql_id[mongo_id] = mysql_object.id

    print('End converting {}'.format(str(model)))
    return mongo_id_to_mysql_id


# Convert books
def conv_books(mongo_books, publisher_ids, author_ids):
    print('Start converting books')

    with db.transaction():
        book_num = mongo_books.count()
        print('Books num:', book_num)
        for i, book in enumerate(mongo_books.find()):
            Books.create(
                kindle=get_book_kindle(book),
                date_publish=book['date_created'],
                image_s_url=book['images']['s']['url'],
                image_s_width=book['images']['s']['size']['width'],
                image_s_height=book['images']['s']['size']['height'],
                image_m_url=book['images']['m']['url'],
                image_m_width=book['images']['m']['size']['width'],
                image_m_height=book['images']['m']['size']['height'],
                image_l_url=book['images']['l']['url'],
                image_l_width=book['images']['l']['size']['width'],
                image_l_height=book['images']['l']['size']['height'],
                asin=book['asin'],
                title=book['title'],
                region=book['region'],
                publisher_id=publisher_ids[str(book['publisher_id'])] if str(book['publisher_id']) in publisher_ids else None,
                author_id=author_ids[str(book['author_id'])] if str(book['author_id']) in author_ids else None,
                date_last_modify=book['date_last_modify'],
                date_created=book['date_created']
            )

            if i % 100 == 0:
                sys.stdout.write('\r{0} {1}%'.format(i + 1, int((i + 1) / book_num * 100)))
                sys.stdout.flush()

    print('')
    print('End converting books')


def get_book_kindle(book):
    value = False
    if 'kindle' in book:
        if book['kindle'] == 1 or book['kindle'] == True:
            value = True
    return value


def main():
    client = MongoClient('localhost', 27017).manganow
    publisher_ids = conv_info_models(Publishers, list(client.publishers.find()))
    author_ids = conv_info_models(Authors, list(client.authors.find()))
    conv_books(client.books, publisher_ids, author_ids)


if __name__ == "__main__":
    main()
