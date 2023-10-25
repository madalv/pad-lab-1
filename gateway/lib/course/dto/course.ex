defmodule Course.Dto.CreateCourse do
  defstruct [:author_id, :title, :description, :category_ids]

  defimpl Nestru.Decoder do
    def from_map_hint(_value, _ctx, map) do
      if Nestru.has_key?(map, :author_id) &&
           Nestru.has_key?(map, :title) do
        {:ok, %{}}
      else
        {:error, "missing keys: title or author_id"}
      end
    end
  end
end
