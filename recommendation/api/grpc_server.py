from concurrent import futures
import time
from pb import rec_svc_pb2_grpc as pb2_grpc
from pb import rec_svc_pb2 as pb2
from google.protobuf import empty_pb2
from service import rec_service as svc
import grpc
import logging

class RecServer(pb2_grpc.RecServiceServicer):

  rec_svc: svc.RecService = None

  def __init__(self, rec_svc):
    self.rec_svc = rec_svc

  def GetRecsForCourse(self, request, context):
    try:
      recs = []
      raw_recs = self.rec_svc.get_recs_for_course(request.course_id, request.recs_nr)
      for rec in raw_recs:
        title, id = rec
        instance = pb2.CourseRec(id=id, title=title)
        recs.append(instance)
      return pb2.RecsResponse(recs = recs)
    except Exception as e:
      context.set_code(grpc.StatusCode.INTERNAL)
      context.set_details(str(e))
    
  def GetRecsForUser(self, request, context):
    try:
      recs = []
      raw_recs = self.rec_svc.get_recs_for_user(request.user_id, request.recs_nr)

      for rec in raw_recs:
        title, id = rec
        instance = pb2.CourseRec(id=id, title=title)
        recs.append(instance)
      return pb2.RecsResponse(recs = recs)
    except Exception as e:
      logging.error(e)
      context.set_code(grpc.StatusCode.INTERNAL)
      context.set_details(str(e))
    
  def AddCourse(self, request, context):
    try:
      request_dict = {
      'author': request.author,
      'categories': " ".join(request.categories),
      'description': request.description,
      'course_id': request.id,
      'title': request.title
    }
      self.rec_svc.add_course(request_dict)
      return empty_pb2.Empty()
    except Exception as e:
      context.set_code(grpc.StatusCode.INTERNAL)
      context.set_details('failed to add course')
    
  def GetServerStatus(self, request, context):
    return pb2.ServerStatus(status="SERVING")

def serve(rec_svc, address):
  logging.info('Starting Rec gRPC Server')
  server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
  pb2_grpc.add_RecServiceServicer_to_server(RecServer(rec_svc), server)
  server.add_insecure_port(address)
  server.start()
  server.wait_for_termination()

  
