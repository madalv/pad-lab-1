import logging
import grpc
from sqlalchemy import create_engine
from repository import course_repo as cr
from service import rec_service as svc
from api import grpc_server
from pb import course_svc_pb2_grpc as pb2_grpc
from pb import course_svc_pb2 as pb2


if __name__ == '__main__':
  logging.basicConfig(level=logging.DEBUG, format='%(asctime)s - %(levelname)s - %(message)s')

  # TODO read this out of .env
  alchemyEngine = create_engine('postgresql+psycopg2://admin:password@rec_db:5432/rec_db', pool_recycle=3600)
  channel = grpc.insecure_channel('course_svc:50052') 
  
  # grpc stub for the Course Service
  stub = pb2_grpc.CourseServiceStub(channel)
  repo = cr.CourseRepo(stub, alchemyEngine)

  
  rec_svc = svc.RecService(repo)

  grpc_server.serve(rec_svc)