# train_model.py
import pandas as pd
import numpy as np
import json
import tensorflow as tf
from tensorflow.keras.layers import Input, Dense, Dropout, BatchNormalization, Concatenate
from tensorflow.keras.models import Model
from tensorflow.keras.optimizers import Adam
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import StandardScaler, OneHotEncoder
from sklearn.compose import ColumnTransformer
import joblib
from tensorflow.keras.callbacks import ReduceLROnPlateau


with open("matchmaking_data.json", "r") as f:
    data = json.load(f)["matches"]

df = pd.DataFrame(data)

print(df)

# Create a continuous label (0-1) from match_percentage
df["label"] = df["match_percentage"] / 100.0

# Define the feature columns for each branch
founder_features = ["fund_required", "industry", "funding_stage"]
investor_features = ["total_invested",
                     "preferred_funding_stage", "risk_tolerance"]
print(founder_features)
print(investor_features)

# Preprocess founder features
founder_preprocessor = ColumnTransformer([
    ("num", StandardScaler(), ["fund_required"]),
    ("cat", OneHotEncoder(handle_unknown="ignore"),
     ["industry", "funding_stage"])
])
X_founder = founder_preprocessor.fit_transform(df[founder_features])

print(X_founder)

# Preprocess investor features
investor_preprocessor = ColumnTransformer([
    ("num", StandardScaler(), ["total_invested"]),
    ("cat", OneHotEncoder(handle_unknown="ignore"),
     ["preferred_funding_stage", "risk_tolerance"])
])
X_investor = investor_preprocessor.fit_transform(df[investor_features])
print(X_investor)

y = df["label"].values
joblib.dump(founder_preprocessor, "founder_preprocessor.pkl")
joblib.dump(investor_preprocessor, "investor_preprocessor.pkl")

# Split data for training and testing
Xf_train, Xf_test, Xi_train, Xi_test, y_train, y_test = train_test_split(
    X_founder, X_investor, y, test_size=0.2, random_state=42
)

print(Xf_train, Xf_test, Xi_train, Xi_test, y_train, y_test)

# Define input dimensions
founder_input_dim = X_founder.shape[1]
investor_input_dim = X_investor.shape[1]
print(founder_input_dim, investor_input_dim)

founder_input = Input(shape=(founder_input_dim,), name="founder_input")
f = Dense(128, activation="relu")(founder_input)  # Increased neurons
f = BatchNormalization()(f)
f = Dropout(0.2)(f)  # Reduced dropout rate
f = Dense(64, activation="relu")(f)
f = BatchNormalization()(f)
f = Dropout(0.2)(f)  # Reduced dropout rate
founder_branch = Dense(16, activation="relu")(f)
print(founder_branch)

investor_input = Input(shape=(investor_input_dim,), name="investor_input")
i = Dense(3364, activation="relu")(investor_input)
i = BatchNormalization()(i)
i = Dropout(0.3)(i)
i = Dense(32, activation="relu")(i)
i = BatchNormalization()(i)
i = Dropout(0.2)(i)
investor_branch = Dense(16, activation="relu")(i)
print(investor_branch)

combined = Concatenate()([founder_branch, investor_branch])
x = Dense(124, activation="relu")(combined)
x = BatchNormalization()(x)
x = Dropout(0.3)(x)
x = Dense(32, activation="relu")(x)
x = BatchNormalization()(x)
x = Dropout(0.2)(x)
output = Dense(1, activation="sigmoid", name="match_probability")(x)
print(output)

# Define the model
model = Model(inputs=[founder_input, investor_input], outputs=output)
model.compile(optimizer=Adam(learning_rate=0.001), loss="mse", metrics=["mae"])
reduce_lr = ReduceLROnPlateau(
    monitor='val_loss', factor=0.2, patience=5, min_lr=0.00001)  # Added ReduceLROnPlateau

# Train the model
history = model.fit(
    x=[Xf_train, Xi_train],
    y=y_train,
    validation_data=([Xf_test, Xi_test], y_test),
    epochs=50,
    batch_size=32,
    callbacks=[reduce_lr]
)
# Save the model in Keras v3 format
model.save("matchmaking_model.keras", save_format="keras_v3")
print("✅ Deep learning model trained and saved in Keras v3 format.")

