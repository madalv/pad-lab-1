import logging
from sqlalchemy import create_engine
from repository import course_repo as cr
from service import rec_service as svc
from api import grpc_server as grpc


if __name__ == '__main__':
  logging.basicConfig(level=logging.DEBUG, format='%(asctime)s - %(levelname)s - %(message)s')

  # TODO read this out of .env
  alchemyEngine = create_engine('postgresql+psycopg2://admin:pass@localhost:5436/rec_db', pool_recycle=3600)
  repo = cr.CourseRepo(None, alchemyEngine)
  
  rec_svc = svc.RecService(repo)

  grpc.serve(rec_svc)