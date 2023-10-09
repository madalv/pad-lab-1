import logging
import pandas as pd
from sklearn.feature_extraction.text import CountVectorizer
from sklearn.metrics.pairwise import cosine_similarity

class RecSystem:
  dataframe: pd.DataFrame = None
  similarities = None
  vectorizer = CountVectorizer()

  def __init__(self, courses):
    logging.info("Init Recommendation System")
    self.dataframe = pd.DataFrame(courses, columns=['course_id', 'title', 'description', 'author', 'categories'])
    self.similarities = self.process_similarities(self.dataframe)
    self.dataframe = self.clean_dataframe(self.dataframe)
  
  def append_course(self, course: dict):
    self.dataframe.loc[len(self.dataframe)] = course
    self.similarities = self.process_similarities(self.dataframe)

  def clean_text(self, author):
    result = str(author).lower()
    return result.replace(' ', '')

  def clean_dataframe(self, df):
    df['author'] = df['author'].apply(self.clean_text)
    df['title'] = df['title'].str.lower()
    df['description'] = df['description'].str.lower()
    df['categories'] = df['categories'].str.lower()
    df = df.drop_duplicates(subset='course_id', keep='first')
    return df

  def process_similarities(self, df):
    # df = self.clean_dataframe(df)

    df['data'] = df[df.columns[1:]].apply(
        lambda x: ' '.join(x.dropna().astype(str)),
        axis=1
    )

    vectorized = self.vectorizer.fit_transform(df['data'])
    cos_vec = cosine_similarity(vectorized)
    similarities = pd.DataFrame(cos_vec, columns=df['course_id'], index=df[['title', 'course_id']]).reset_index()

    return similarities
  
  def get_recs(self, course_id, nr):
    """
    Returns list of tuples of the form [(title, id), (title, id), ...]
    """
    recs = pd.DataFrame(self.similarities.nlargest(nr + 1, course_id)['index'])
    recs = recs[recs['index'].apply(lambda x: x[1]) != course_id]
    recs['index'] = recs['index'].apply(lambda x: (x[0].title(), x[1]))
    return recs.values.flatten().tolist()
