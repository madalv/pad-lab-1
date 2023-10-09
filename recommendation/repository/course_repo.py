import logging
from sqlalchemy import MetaData, Table, Column, Integer, String

class CourseRepo:
  course_grpc_client = None
  engine = None
  metadata = None
  course_data = None
  
  def __init__(self, client, db_conn):
    logging.info("Init Course Repository")
    self.course_grpc_client = client
    self.engine = db_conn
    
    self.metadata = MetaData()
    self.course_data = Table(
      'course_data',
      self.metadata,
      Column('course_id', String, primary_key=True),
      Column('title', String),
      Column('description', String),
      Column('author', String),
      Column('categories', String)
      )
    self.metadata.create_all(self.engine)

  def fetch_courses(self):
    with self.engine.connect() as connection:
      select_query = self.course_data.select()
      results = connection.execute(select_query).fetchall()
    
    return results
  
  def add_course(self, course: dict):
    with self.engine.connect() as connection:
        insert_query = self.course_data.insert().values(course)
        connection.execute(insert_query)
        connection.commit()

  def fetch_user_courses(self, user_id):
    # TODO implement
    return
