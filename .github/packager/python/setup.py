from pathlib import Path

from setuptools import find_packages, setup

setup(
    name="pactus-grpc",
    version="{{ VERSION }}",
    author="Pactus Development Team",
    author_email="info@pactus.org",
    description="gRPC client bindings for the Pactus blockchain",
    long_description=Path("README.md").read_text(encoding="utf-8"),
    long_description_content_type="text/markdown",
    url="https://github.com/pactus-project/pactus",
    packages=find_packages(),
    license="MIT",
    install_requires=[
        "grpcio",
        "protobuf",
    ],
    keywords=["pactus", "blockchain", "grpc"],
    classifiers=[
        "Development Status :: 5 - Production/Stable",
        "Intended Audience :: Developers",
        "Topic :: Software Development :: Build Tools",
        "Operating System :: OS Independent",
    ],
    python_requires=">=3.6",
)
