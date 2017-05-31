# -*- coding:utf-8 -*-
import sys
import datetime
from pymongo import MongoClient
from peewee import *
from models import db, Publishers, Authors, Books, Titles


# Convert publishers and authors
def conv_info_models(model, mongo_objects):
    print('Start converting {}'.format(str(model)))

    # Create mysql objects
    with db.transaction():
        for mongo_object in mongo_objects:
            model.create(
                name=mongo_object['name'].strip(),
                created_at=mongo_object['date_created'])

    # Make a map of mongo_id to mysql_id
    mongo_id_to_mysql_id = {str(mongo_object['_id']): i + 1 for i, mongo_object in enumerate(mongo_objects)}

    print('End converting {}'.format(str(model)))
    return mongo_id_to_mysql_id


# Convert books
def conv_books(mongo_books, publisher_ids, author_ids):
    print('Start converting books')

    with db.transaction():
        for i, book in enumerate(mongo_books):
            Books.create(
                asin=book['asin'],
                publish_type=book['binding'] if 'binding' in book else None,
                date_publish='%d%02d%02d' % (book['date_publish'].year, book['date_publish'].month, book['date_publish'].day),
                image_s_url=book['images']['s']['url'],
                image_s_width=book['images']['s']['size']['width'],
                image_s_height=book['images']['s']['size']['height'],
                image_m_url=book['images']['m']['url'],
                image_m_width=book['images']['m']['size']['width'],
                image_m_height=book['images']['m']['size']['height'],
                image_l_url=book['images']['l']['url'],
                image_l_width=book['images']['l']['size']['width'],
                image_l_height=book['images']['l']['size']['height'],
                name=book['title'],
                region=book['region'],
                publisher_id=publisher_ids[str(book['publisher_id'])] if str(book['publisher_id']) in publisher_ids else None,
                author_id=author_ids[str(book['author_id'])] if str(book['author_id']) in author_ids else None,
                updated_at=book['date_last_modify'],
                created_at=book['date_created']
            )

            if i % 100 == 0:
                book_num = len(mongo_books)
                sys.stdout.write('\r{0} {1}%'.format(i + 1, int((i + 1) / book_num * 100)))
                sys.stdout.flush()

    # Make a map of mongo_id to mysql_id
    mongo_id_to_mysql_id = {str(mongo_object['_id']): i + 1 for i, mongo_object in enumerate(mongo_books)}

    print('')
    print('End converting books')

    return mongo_id_to_mysql_id


def update_book_title(asins, title_id):
    Books.update(book_title=title_id).where(Books.asin << asins).execute()


def create_titles(mongo_books, book_ids):
    print('Start creating titles')

    with db.transaction():
        title_id = 0
        for i, book in enumerate(mongo_books):
            if not 'tree_type' in book or book['tree_type'] == 'main':
                Titles.create(
                    created_at=datetime.datetime.now()
                )
                title_id += 1
                asins = [book['asin']]
                if 'sub_asins' in book and book['sub_asins'] and len(book['sub_asins']):
                    asins.extend(book['sub_asins'])
                update_book_title(asins, title_id)

            if i % 100 == 0:
                book_num = len(mongo_books)
                sys.stdout.write('\r{0} {1}%'.format(i + 1, int((i + 1) / book_num * 100)))
                sys.stdout.flush()

    print('')
    print('End creating titles')


def main():
    client = MongoClient('localhost', 27017).manganow
    publisher_ids = conv_info_models(Publishers, list(client.publishers.find()))
    author_ids = conv_info_models(Authors, list(client.authors.find()))
    mongo_books = list(client.books.find())
    book_ids = conv_books(mongo_books, publisher_ids, author_ids)
    create_titles(mongo_books, book_ids)


if __name__ == "__main__":
    main()
