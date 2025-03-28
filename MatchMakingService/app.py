# app.py
from fastapi import FastAPI
import joblib
import tensorflow as tf
import pandas as pd

app = FastAPI()

# Load the trained model and preprocessors
model = tf.keras.models.load_model("matchmaking_model.keras")
print(model.summary())
founder_preprocessor = joblib.load("founder_preprocessor.pkl")
investor_preprocessor = joblib.load("investor_preprocessor.pkl")


@app.post("/predict/")
async def predict_match(data: dict):
    """
    Expected JSON structure:
    {
      "founder": {
         "fund_required": 500000,
         "industry": "AI/ML",
         "funding_stage": "Series A"
      },
      "investor": {
         "total_invested": 1000000,
         "preferred_funding_stage": "Series A",
         "risk_tolerance": "Moderate"
      }
    }
    """
    # Create DataFrames from the input data
    df_founder = pd.DataFrame([data["founder"]])
    df_investor = pd.DataFrame([data["investor"]])

    # Preprocess the data
    X_founder = founder_preprocessor.transform(df_founder)
    X_investor = investor_preprocessor.transform(df_investor)

    # Predict match probability
    prediction = model.predict([X_founder, X_investor])
    match_probability = float(prediction[0][0]) * 100  # Convert to percentage

    # trancate the match_probability to 2 decimal points
    match_probability = round(match_probability, 2)

# write test codes

    return {"match_probability": match_probability}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=4040)
