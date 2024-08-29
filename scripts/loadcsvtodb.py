import pandas as pd
from sqlalchemy import create_engine
from dotenv import load_dotenv
import os

# Load environment variables from .env file
load_dotenv()

# Get database credentials from environment variables
user = os.getenv('POSTGRES_USER')
password = os.getenv('POSTGRES_PASSWORD')
host = os.getenv('POSTGRES_HOST')
database = os.getenv('POSTGRES_DB')

# Load CSV into DataFrame
print('CSV parsing started...')
df = pd.read_csv('/Users/anigokul/Downloads/Reviews.csv', low_memory=False)
print('CSV has been read, connecting to database...')

# Create a connection to PostgreSQL
engine = create_engine(f'postgresql://{user}:{password}@{host}/{database}')
print('Database connected, upload started...')

# Upload the DataFrame to PostgreSQL
df.to_sql('amazonreviews', engine, index=False, if_exists='replace')
print('Successfully uploaded to Postgres DB!')
