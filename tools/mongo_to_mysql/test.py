# -*- coding:utf-8 -*-
import sys
from pymongo import MongoClient
from bson.objectid import ObjectId
from peewee import *
from models import db, Books


def main():
    client = MongoClient('localhost', 27017).manganow
    mongo_publishers = list(client.publishers.find())
    mongo_authors = list(client.authors.find())

    tests = [
        {   'title': 'ふたりごはん (MARBLE COMICS)',
            'publisher_name': 'ソフトライン 東京漫画社',
            'author_name': 'テラシマ'},
        {   'title': '鋼の錬金術師 完全版(10) (ガンガンコミックスデラックス)',
            'publisher_name': 'スクウェア・エニックス',
            'author_name': '荒川 弘'},
        {   'title': 'ノーゲーム・ノーライフ5 ゲーマー兄妹は強くてニューゲームがお嫌いなようです (MF文庫J)',
            'publisher_name': 'メディアファクトリー',
            'author_name': '榎宮祐'},
        {   'title': 'はたらく細胞（２） (シリウスコミックス)',
            'publisher_name': '講談社',
            'author_name': '清水茜'},
    ]

    for test in tests:
        book = Books.get(Books.title == test['title'])
        assert(book.publisher.name == test['publisher_name'])
        assert(book.author.name == test['author_name'])
        print('OK', test['title'], test['publisher_name'], test['author_name'])

    print('Success!')


if __name__ == "__main__":
    main()
