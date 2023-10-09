from google.protobuf import empty_pb2 as _empty_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class CourseRecsRequest(_message.Message):
    __slots__ = ["course_id", "recs_nr"]
    COURSE_ID_FIELD_NUMBER: _ClassVar[int]
    RECS_NR_FIELD_NUMBER: _ClassVar[int]
    course_id: str
    recs_nr: int
    def __init__(self, course_id: _Optional[str] = ..., recs_nr: _Optional[int] = ...) -> None: ...

class UserRecsRequest(_message.Message):
    __slots__ = ["user_id", "recs_nr"]
    USER_ID_FIELD_NUMBER: _ClassVar[int]
    RECS_NR_FIELD_NUMBER: _ClassVar[int]
    user_id: str
    recs_nr: int
    def __init__(self, user_id: _Optional[str] = ..., recs_nr: _Optional[int] = ...) -> None: ...

class Course(_message.Message):
    __slots__ = ["id", "title", "description", "categories", "author"]
    ID_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    CATEGORIES_FIELD_NUMBER: _ClassVar[int]
    AUTHOR_FIELD_NUMBER: _ClassVar[int]
    id: str
    title: str
    description: str
    categories: _containers.RepeatedScalarFieldContainer[str]
    author: str
    def __init__(self, id: _Optional[str] = ..., title: _Optional[str] = ..., description: _Optional[str] = ..., categories: _Optional[_Iterable[str]] = ..., author: _Optional[str] = ...) -> None: ...

class CourseRec(_message.Message):
    __slots__ = ["id", "title"]
    ID_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    id: str
    title: str
    def __init__(self, id: _Optional[str] = ..., title: _Optional[str] = ...) -> None: ...

class RecsResponse(_message.Message):
    __slots__ = ["recs"]
    RECS_FIELD_NUMBER: _ClassVar[int]
    recs: _containers.RepeatedCompositeFieldContainer[CourseRec]
    def __init__(self, recs: _Optional[_Iterable[_Union[CourseRec, _Mapping]]] = ...) -> None: ...
