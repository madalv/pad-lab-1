defmodule Proto.Rec.ServerStatus do
  @moduledoc false

  use Protobuf, syntax: :proto3, protoc_gen_elixir_version: "0.12.0"

  field :status, 1, type: :string
end

defmodule Proto.Rec.CourseRecsRequest do
  @moduledoc false

  use Protobuf, syntax: :proto3, protoc_gen_elixir_version: "0.12.0"

  field :course_id, 1, type: :string, json_name: "courseId"
  field :recs_nr, 2, type: :uint64, json_name: "recsNr"
end

defmodule Proto.Rec.UserRecsRequest do
  @moduledoc false

  use Protobuf, syntax: :proto3, protoc_gen_elixir_version: "0.12.0"

  field :user_id, 1, type: :string, json_name: "userId"
  field :recs_nr, 2, type: :uint64, json_name: "recsNr"
end

defmodule Proto.Rec.Course do
  @moduledoc false

  use Protobuf, syntax: :proto3, protoc_gen_elixir_version: "0.12.0"

  field :id, 1, type: :string
  field :title, 2, type: :string
  field :description, 3, type: :string
  field :categories, 4, repeated: true, type: :string
  field :author, 5, type: :string
end

defmodule Proto.Rec.CourseRec do
  @moduledoc false

  use Protobuf, syntax: :proto3, protoc_gen_elixir_version: "0.12.0"

  field :id, 1, type: :string
  field :title, 2, type: :string
end

defmodule Proto.Rec.RecsResponse do
  @moduledoc false

  use Protobuf, syntax: :proto3, protoc_gen_elixir_version: "0.12.0"

  field :recs, 1, repeated: true, type: Proto.Rec.CourseRec
end

defmodule Proto.Rec.RecService.Service do
  @moduledoc false

  use GRPC.Service, name: "proto.rec.RecService", protoc_gen_elixir_version: "0.12.0"

  rpc :GetRecsForCourse, Proto.Rec.CourseRecsRequest, Proto.Rec.RecsResponse

  rpc :GetRecsForUser, Proto.Rec.UserRecsRequest, Proto.Rec.RecsResponse

  rpc :AddCourse, Proto.Rec.Course, Google.Protobuf.Empty

  rpc :GetServerStatus, Google.Protobuf.Empty, Proto.Rec.ServerStatus
end

defmodule Proto.Rec.RecService.Stub do
  @moduledoc false

  use GRPC.Stub, service: Proto.Rec.RecService.Service
end