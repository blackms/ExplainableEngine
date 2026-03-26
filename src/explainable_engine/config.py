from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    app_name: str = "Explainable Engine"
    version: str = "0.1.0"
    debug: bool = False
    store_backend: str = "memory"  # "memory" | "sqlite"
    sqlite_path: str = "explanations.db"
    log_level: str = "INFO"

    model_config = {"env_prefix": "EE_"}


settings = Settings()
