defmodule Chapter.Dto.CreateChapter do
  defstruct [:course_id, :title, :body]

  defimpl Nestru.Decoder do
    def from_map_hint(_value, _ctx, map) do
      if Nestru.has_key?(map, :course_id) &&
           Nestru.has_key?(map, :title) &&
           Nestru.has_key?(map, :body) do
        {:ok, %{}}
      else
        {:error, "missing keys: course_id, title, or body"}
      end
    end
  end
end
