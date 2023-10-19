# Use an official Elixir runtime as a parent image
FROM elixir:latest

# Set the environment to build and run the app
ENV MIX_ENV=prod

# Set the working directory in the container
WORKDIR /app

# Copy the mix.exs and mix.lock files to install dependencies
COPY mix.exs mix.lock ./

# Install hex and rebar (Elixir package manager and build tool)
RUN mix local.hex --force && mix local.rebar --force

# Install and compile the application dependencies
RUN mix do deps.get, deps.compile

# Copy the rest of the application code
COPY . .

# Build the release
RUN mix release

# Start your application
CMD ["_build/prod/rel/gateway/bin/gateway", "start"]
