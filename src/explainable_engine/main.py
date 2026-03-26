from fastapi import FastAPI

from explainable_engine.config import settings


def create_app() -> FastAPI:
    application = FastAPI(
        title=settings.app_name,
        version=settings.version,
        description=(
            "Transform any numerical/decisional output into an explicit,"
            " queryable, verifiable causal chain"
        ),
    )
    return application


app = create_app()
