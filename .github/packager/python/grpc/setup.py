from pathlib import Path

from setuptools import find_packages, setup

setup(
    name="pactus-grpc",
    version="{{ VERSION }}",
    author="Pactus Development Team",
    author_email="info@pactus.org",
    url="https://pactus.org",
    description="Python client for interacting with the Pactus blockchain via gRPC",
    long_description=Path("README.md").read_text(encoding="utf-8"),
    long_description_content_type="text/markdown",
    license="MIT",
    packages=find_packages(),
    keywords=["pactus", "blockchain", "grpc"],
    install_requires=[
        "grpcio",
        "protobuf",
    ],
    classifiers=[
        "Development Status :: 5 - Production/Stable",
        "Intended Audience :: Developers",
        "Topic :: Software Development :: Build Tools",
        "Operating System :: OS Independent",
    ],
    python_requires=">=3.6",
)
