import pytest
from fastapi.testclient import TestClient
from app import app

client = TestClient(app)

def test_predict_match():
    # Test data
    test_data = {
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
    
    # Make request to the endpoint
    response = client.post("/predict/", json=test_data)
    
    # Assert response
    assert response.status_code == 200
    assert "match_probability" in response.json()
    assert isinstance(response.json()["match_probability"], float)
    assert 0 <= response.json()["match_probability"] <= 100