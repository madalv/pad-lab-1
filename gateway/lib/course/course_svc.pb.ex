defmodule Proto.Course.ServerStatus do
  @moduledoc false

  use Protobuf, syntax: :proto3, protoc_gen_elixir_version: "0.12.0"

  field(:status, 1, type: :string)
end

defmodule Proto.Course.PaginationQuery do
  @moduledoc false

  use Protobuf, syntax: :proto3, protoc_gen_elixir_version: "0.12.0"

  field(:page, 1, type: :uint64)
  field(:limit, 2, type: :uint64)
end

defmodule Proto.Course.Category do
  @moduledoc false

  use Protobuf, syntax: :proto3, protoc_gen_elixir_version: "0.12.0"

  field(:id, 1, type: :string)
  field(:title, 2, type: :string)
end

defmodule Proto.Course.Course do
  @moduledoc false

  use Protobuf, syntax: :proto3, protoc_gen_elixir_version: "0.12.0"

  field(:id, 1, type: :string)
  field(:author_id, 2, type: :string, json_name: "authorId")
  field(:title, 3, type: :string)
  field(:description, 4, type: :string)
  field(:categories, 5, repeated: true, type: Proto.Course.Category)
  field(:created_at, 6, type: Google.Protobuf.Timestamp, json_name: "createdAt")
  field(:updated_at, 7, type: Google.Protobuf.Timestamp, json_name: "updatedAt")
end

defmodule Proto.Course.Chapter do
  @moduledoc false

  use Protobuf, syntax: :proto3, protoc_gen_elixir_version: "0.12.0"

  field(:id, 1, type: :string)
  field(:course_id, 2, type: :string, json_name: "courseId")
  field(:title, 3, type: :string)
  field(:body, 4, type: :string)
  field(:created_at, 5, type: Google.Protobuf.Timestamp, json_name: "createdAt")
  field(:updated_at, 6, type: Google.Protobuf.Timestamp, json_name: "updatedAt")
end

defmodule Proto.Course.ChapterTitle do
  @moduledoc false

  use Protobuf, syntax: :proto3, protoc_gen_elixir_version: "0.12.0"

  field(:id, 1, type: :string)
  field(:title, 2, type: :string)
end

defmodule Proto.Course.CourseWithChapters do
  @moduledoc false

  use Protobuf, syntax: :proto3, protoc_gen_elixir_version: "0.12.0"

  field(:course, 1, type: Proto.Course.Course)
  field(:chapters, 2, repeated: true, type: Proto.Course.ChapterTitle)
end

defmodule Proto.Course.Courses do
  @moduledoc false

  use Protobuf, syntax: :proto3, protoc_gen_elixir_version: "0.12.0"

  field(:courses, 1, repeated: true, type: Proto.Course.Course)
end

defmodule Proto.Course.CourseId do
  @moduledoc false

  use Protobuf, syntax: :proto3, protoc_gen_elixir_version: "0.12.0"

  field(:id, 1, type: :string)
end

defmodule Proto.Course.ChapterId do
  @moduledoc false

  use Protobuf, syntax: :proto3, protoc_gen_elixir_version: "0.12.0"

  field(:id, 1, type: :string)
end

defmodule Proto.Course.UserId do
  @moduledoc false

  use Protobuf, syntax: :proto3, protoc_gen_elixir_version: "0.12.0"

  field(:id, 1, type: :string)
end

defmodule Proto.Course.CreateCourseRequest do
  @moduledoc false

  use Protobuf, syntax: :proto3, protoc_gen_elixir_version: "0.12.0"

  field(:author_id, 1, type: :string, json_name: "authorId")
  field(:title, 2, type: :string)
  field(:description, 3, type: :string)
  field(:category_ids, 4, repeated: true, type: :string, json_name: "categoryIds")
end

defmodule Proto.Course.CreateChapterRequest do
  @moduledoc false

  use Protobuf, syntax: :proto3, protoc_gen_elixir_version: "0.12.0"

  field(:course_id, 1, type: :string, json_name: "courseId")
  field(:title, 2, type: :string)
  field(:body, 3, type: :string)
end

defmodule Proto.Course.CourseIds do
  @moduledoc false

  use Protobuf, syntax: :proto3, protoc_gen_elixir_version: "0.12.0"

  field(:ids, 1, repeated: true, type: :string)
end

defmodule Proto.Course.CourseService.Service do
  @moduledoc false

  use GRPC.Service, name: "proto.course.CourseService", protoc_gen_elixir_version: "0.12.0"

  rpc(:GetCourses, Proto.Course.PaginationQuery, Proto.Course.Courses)

  rpc(:GetCourse, Proto.Course.CourseId, Proto.Course.CourseWithChapters)

  rpc(:CreateCourse, Proto.Course.CreateCourseRequest, Proto.Course.CourseId)

  rpc(:GetChapter, Proto.Course.ChapterId, Proto.Course.Chapter)

  rpc(:CreateChapter, Proto.Course.CreateChapterRequest, Proto.Course.ChapterId)

  rpc(:GetCourseIdsForUser, Proto.Course.UserId, Proto.Course.CourseIds)

  rpc(:GetServerStatus, Google.Protobuf.Empty, Proto.Course.ServerStatus)
end

defmodule Proto.Course.CourseService.Stub do
  @moduledoc false

  use GRPC.Stub, service: Proto.Course.CourseService.Service
end
