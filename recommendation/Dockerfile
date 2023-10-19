# Stage 1: Build the application
FROM python:3.9-slim

# Set the working directory in the builder stage
WORKDIR /app

# Copy the requirements file and install the dependencies
COPY requirements.txt ./
RUN pip install -r requirements.txt

# Copy the rest of the application code
COPY . .

# RUN pip install grpcio grpcio-tools SQLAlchemy pandas scikit-learn psycopg2-binary

# Define the command to run your application
CMD ["python3", "__init__.py"]
