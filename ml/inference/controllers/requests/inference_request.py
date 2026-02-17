from pydantic import BaseModel, Field

class PredictRequest(BaseModel):
    map_name: str
    mode: str
    rank: int
    team_a: list[int] = Field(default_factory=list)
    team_b: list[int] = Field(default_factory=list)

class RecommendRequest(BaseModel):
    map_name: str
    mode: str
    rank: int
    ally_brawlers: list[int] = Field(default_factory=list)
    enemy_brawlers: list[int] = Field(default_factory=list)
    candidate_brawlers: list[int] | None = None
    banned_brawlers: list[int] = Field(default_factory=list)
    top_k: int | None = None
