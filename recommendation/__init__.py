import logging
import os
import grpc
from sqlalchemy import create_engine
from repository import course_repo as cr
from service import rec_service as svc
from api import grpc_server
from pb import course_svc_pb2_grpc as pb2_grpc
from pb import course_svc_pb2 as pb2
from dotenv import load_dotenv


if __name__ == '__main__':
  logging.basicConfig(level=logging.DEBUG, format='%(asctime)s - %(levelname)s - %(message)s')

  app_mode = os.environ.get("APP_MODE", "prod")
  if app_mode == "prod":
    load_dotenv(dotenv_path=".prod.env")
  elif app_mode == "dev":
    load_dotenv(dotenv_path=".local.env")
  else:
    logging.fatal(f"Invalid mode: {app_mode}. Supported modes are 'dev' and 'prod'.")

  alchemyEngine = create_engine(os.getenv("POSTGRES_URL"), pool_recycle=3600)
  channel = grpc.insecure_channel(os.getenv("COURSE_SVC_ADDRESS"))
  
  # grpc stub for the Course Service
  stub = pb2_grpc.CourseServiceStub(channel)
  repo = cr.CourseRepo(stub, alchemyEngine)

  
  rec_svc = svc.RecService(repo)
  address = os.getenv("GRPC_ADDRESS")

  grpc_server.serve(rec_svc, address)